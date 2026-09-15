package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mimokpl/miwin-toolkit/miwin/internal/pkg"
)

func TestGenerate(t *testing.T) {
	err := Generate(context.Background(), GeneratorOptions{
		GenerateMain:     true,
		GenerateServer:   true,
		GenerateService:  true,
		GenerateData:     true,
		GenerateMakefile: true,
		GenerateConfigs:  true,

		ProjectModule: "github.com/miwin-example",
		ProjectName:   "miwin-example",
		ServiceName:   "user",

		Servers:   []string{"rest", "grpc"},
		DbClients: []string{"ent", "redis"},

		OutputPath: t.TempDir(),
	})
	assert.Nil(t, err)
}

func TestExtractProjectName(t *testing.T) {
	projectModule := "github.com/miwin-example"
	projectName := pkg.ExtractProjectName(projectModule)
	assert.Equal(t, "miwin-example", projectName)

	projectModule = "miwin-example"
	projectName = pkg.ExtractProjectName(projectModule)
	assert.Equal(t, "miwin-example", projectName)
}
