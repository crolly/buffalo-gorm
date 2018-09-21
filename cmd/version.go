package cmd

import (
	"fmt"

	"github.com/crolly/buffalo-gorm/gorm"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "current version of gorm",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("gorm", gorm.Version)
		return nil
	},
}

func init() {
	gormCmd.AddCommand(versionCmd)
}
