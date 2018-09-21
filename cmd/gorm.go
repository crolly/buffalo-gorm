package cmd

import (
	"github.com/spf13/cobra"
)

// gormCmd represents the gorm command
var gormCmd = &cobra.Command{
	Use:   "gorm",
	Short: "description about this plugin",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(gormCmd)
}
