package project

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/mimokpl/miwin-toolkit/miwin/internal/buf"
	"github.com/mimokpl/miwin-toolkit/miwin/internal/pkg"
)

// CmdProject represents the project command.
var CmdProject = &cobra.Command{
	Use:     "project [name]",
	Aliases: []string{"proj"},
	Short:   "create a new project scaffold",
	Long:    "Create a project using the repository template. Example: gow new project helloworld",
	Args:    cobra.ExactArgs(1),
	RunE:    Run,
	SilenceUsage: true,
}

var (
	repoURL    string
	branch     string
	timeout    string
	moduleName string
	nomod      bool
)

const (
	GithubRepoURL = "https://github.com/mimokpl/miwin-admin-template.git"
	GiteeRepoURL  = "https://gitee.com/miwin/miwin-admin-template.git"

	// FallbackRepoURL 上游可用的同款模板，作为 miwin 模板仓库尚未发布时的兜底。
	FallbackRepoURL = "https://github.com/tx7do/go-wind-admin-template.git"
)

// templateFallbackRepoURL 兜底模板仓库，可用环境变量 MIWIN_TEMPLATE_FALLBACK_REPO 覆盖。
func templateFallbackRepoURL() string {
	if v := strings.TrimSpace(os.Getenv("MIWIN_TEMPLATE_FALLBACK_REPO")); v != "" {
		return v
	}
	return FallbackRepoURL
}

// repoReachable 用 git ls-remote 探测仓库是否可克隆。
func repoReachable(repoURL string, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return exec.CommandContext(ctx, "git", "ls-remote", "--exit-code", repoURL, "HEAD").Run() == nil
}

// 模板仓库地址可用环境变量覆盖，便于指向自建模板或上游临时仓库：
//   MIWIN_TEMPLATE_REPO        GitHub 地址
//   MIWIN_TEMPLATE_GITEE_REPO  Gitee 地址（GitHub 不可达时的回退）
//   MIWIN_TEMPLATE_MODULE      模板的 go module 路径
func templateRepoURL() string {
	if v := strings.TrimSpace(os.Getenv("MIWIN_TEMPLATE_REPO")); v != "" {
		return v
	}
	return GithubRepoURL
}

func templateGiteeRepoURL() string {
	if v := strings.TrimSpace(os.Getenv("MIWIN_TEMPLATE_GITEE_REPO")); v != "" {
		return v
	}
	return GiteeRepoURL
}

func canReach(addr string, d time.Duration) bool {
	conn, err := net.DialTimeout("tcp", addr, d)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func init() {
	timeout = "60s"

	CmdProject.Flags().StringVarP(&repoURL, "repo-url", "r", templateRepoURL(), "layout repo")
	CmdProject.Flags().StringVarP(&branch, "branch", "b", branch, "repo branch")
	CmdProject.Flags().StringVarP(&timeout, "timeout", "t", timeout, "time out")
	CmdProject.Flags().StringVarP(&moduleName, "module", "m", moduleName, "set go module name, if not set, use project name")
	CmdProject.Flags().BoolVarP(&nomod, "nomod", "", nomod, "retain go mod")
}

func Run(cmd *cobra.Command, args []string) error {
	// Default endpoint (no explicit -r): prefer GitHub, fall back to Gitee if
	// unreachable. Probed here rather than in init() so unrelated commands
	// don't pay a blocking network round-trip on startup.
	if !cmd.Flags().Changed("repo-url") {
		if canReach("github.com:443", 3*time.Second) {
			repoURL = templateRepoURL()
		} else {
			repoURL = templateGiteeRepoURL()
		}
		// 默认模板仓库尚未发布时，自动回退到上游可用模板，保证开箱能用。
		if !repoReachable(repoURL, 10*time.Second) {
			if fb := templateFallbackRepoURL(); fb != "" && fb != repoURL {
				log.Printf("⚠️  模板仓库 %s 不可用，回退到 %s", repoURL, fb)
				repoURL = fb
			}
		}
	}

	t, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("invalid timeout value %q: %w", timeout, err)
	}

	parentCtx := cmd.Context()
	if parentCtx == nil {
		parentCtx = context.Background()
	}
	ctx, cancel := context.WithTimeout(parentCtx, t)
	defer cancel()

	name := ""
	if len(args) == 0 {
		prompt := &survey.Input{
			Message: "What is project name ?",
			Help:    "Created project name.",
		}
		err = survey.AskOne(prompt, &name)
		if err != nil || name == "" {
			return nil
		}
	} else {
		name = args[0]
	}

	projectName, workingDir := processProjectParams(name)

	if isDirExists(workingDir, projectName) {
		fmt.Printf("🚫 %s already exists\n", projectName)
		prompt := &survey.Confirm{
			Message: "📂 Do you want to override the folder ?",
			Help:    "Delete the existing folder and create the project.",
		}
		var override bool
		e := survey.AskOne(prompt, &override)
		if e != nil {
			return nil
		}
		if !override {
			return nil
		}
		_ = os.RemoveAll(filepath.Join(workingDir, projectName))
	}

	fmt.Printf("🚀 Creating project %s, layout repo is %s, please wait a moment.\n\n", projectName, repoURL)

	p := &Project{
		Name:   projectName,
		Module: projectName,
	}
	if moduleName != "" {
		p.Module = moduleName
	}

	done := make(chan error, 1)
	go func() {
		if !nomod {
			done <- p.New(ctx, workingDir, repoURL, branch)
			return
		}
		projectRoot := getGoModProjectRoot(workingDir)
		if goModIsNotExistIn(projectRoot) {
			done <- fmt.Errorf("🚫 go.mod don't exists in %s", projectRoot)
			return
		}

			packagePath, e := filepath.Rel(projectRoot, filepath.Join(workingDir, projectName))
			if e != nil {
				done <- fmt.Errorf("🚫 failed to get relative path: %v", e)
				return
			}
		packagePath = strings.ReplaceAll(packagePath, "\\", "/")

		mod, e := pkg.ModulePath(filepath.Join(projectRoot, "go.mod"))
		if e != nil {
			done <- fmt.Errorf("🚫 failed to parse `go.mod`: %v", e)
			return
		}
		// Get the relative path for adding a project based on Go modules
		p.Path = filepath.Join(strings.TrimPrefix(workingDir, projectRoot+"/"), p.Name)
		done <- p.Add(ctx, workingDir, repoURL, branch, mod, packagePath)
	}()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("project creation timed out")
		}
		return fmt.Errorf("failed to create project: %w", ctx.Err())

	case createErr := <-done:
		if createErr != nil {
			return fmt.Errorf("failed to create project: %w（默认模板仓库 %s 可能不存在；可用 -r/--repo-url 或环境变量 MIWIN_TEMPLATE_REPO 指定模板）", createErr, templateRepoURL())
		}

		// 先按 proto 生成 api 代码，再 go mod tidy：
		// 模板仓库通常不提交 api/gen，先生成才能解析这些 import。
		if err = buf.GenerateFromPath(ctx, filepath.Join(workingDir, projectName, "api")); err != nil {
			return fmt.Errorf("failed to generate api code: %w", err)
		}

		if err = pkg.GoModTidy(ctx, filepath.Join(workingDir, projectName)); err != nil {
			return fmt.Errorf("failed to run `go mod tidy`: %w", err)
		}

		fmt.Printf("✅ Project %s created successfully at %s\n", projectName, filepath.Join(workingDir, projectName))
		return nil
	}
}
