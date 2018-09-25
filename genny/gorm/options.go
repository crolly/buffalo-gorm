package gorm

import (
	"errors"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/genny/movinglater/attrs"
	"github.com/markbates/inflect"
)

type Options struct {
	App       meta.App     `json:"app"`
	Name      inflect.Name `json:"name"`
	ModelName inflect.Name `json:"model_name"`
	Model     Model        `json:"model"`
	Char      string       `json:"char"`
	// SkipMigration bool         `json:"skip_migration"`
	// SkipModel     bool         `json:"skip_model"`
	// SkipTemplates bool         `json:"skip_templates"`
	// UseModel      bool         `json:"use_model"`
	FilesPath   string           `json:"files_path"`
	ActionsPath string           `json:"actions_path"`
	Props       []Prop           `json:"props"`
	NamedAttrs  attrs.NamedAttrs `json:"named_attrs"`
	Actions     []string         `json:"actions"`
	Args        []string         `json:"args"`

	Init bool `json:"init"`
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts == nil || opts.ModelName == "" {
		return errors.New("you must specify a resource name")
	}
	return nil
}
