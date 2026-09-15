package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/mimokpl/miwin-toolkit/miwin/internal/extract"
	"github.com/mimokpl/miwin-toolkit/miwin/internal/generate"
	"github.com/mimokpl/miwin-toolkit/miwin/internal/project"
	"github.com/mimokpl/miwin-toolkit/miwin/internal/run"
)

var rootCmd = &cobra.Command{
	Use:   "gow",
	Short: "gow CLI",
	Long:  "gow is the CLI for Miwin framework.",
}

func init() {
	rootCmd.AddCommand(project.CmdProject)
	rootCmd.AddCommand(run.CmdRun)
	rootCmd.AddCommand(generate.CmdGenerate)
	rootCmd.AddCommand(extract.CmdExtract)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
