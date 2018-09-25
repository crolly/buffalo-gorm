package cmd

import (
	"context"

	"github.com/crolly/buffalo-gorm/genny/gorm"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var generateOptions = struct {
	*gorm.Options
	dryRun bool
}{
	Options: &gorm.Options{},
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "gorm [name]",
	Short: "A collection of gorm generators",
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := generateOptions
		opts.Args = args

		r := genny.WetRunner(context.Background())

		if generateOptions.dryRun {
			r = genny.DryRunner(context.Background())
		}

		g, err := gorm.New(opts.Options)
		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		g, err = gotools.GoFmt(r.Root)
		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		return r.Run()
	},
}

func init() {
	generateCmd.Flags().BoolVarP(&generateOptions.Init, "init", "i", false, "initialize gorm for the project")
	generateCmd.Flags().BoolVarP(&generateOptions.dryRun, "dry-run", "d", false, "run the generator without creating files or running commands")
	gormCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(generateCmd)
}
