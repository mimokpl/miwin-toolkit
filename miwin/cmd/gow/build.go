package main

import (
	"github.com/mimokpl/miwin-toolkit/miwin/internal/build"
)

func init() {
	rootCmd.AddCommand(build.CmdBuild)
}
