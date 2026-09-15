package main

import (
	"github.com/mimokpl/miwin-toolkit/miwin/internal/migrate"
)

func init() {
	rootCmd.AddCommand(migrate.CmdMigrate)
}
