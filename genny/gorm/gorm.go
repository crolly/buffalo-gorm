package gorm

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/genny/movinglater/attrs"
	"github.com/gobuffalo/genny/movinglater/gotools"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr"
	"github.com/markbates/inflect"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	pwd, _ := os.Getwd()
	opts.App = meta.New(pwd)

	if len(opts.Args) > 0 {
		opts.Name = inflect.Name(opts.Args[0])
		opts.Model = inflect.Name(opts.Args[0])
	}

	var err error
	opts.NamedAttrs, err = attrs.ParseNamedArgs(opts.Args...)
	if err != nil {
		return g, errors.WithStack(err)
	}

	opts.Props = modelPropertiesFromArgs(opts.Args)

	opts.FilesPath = opts.Name.PluralUnder()
	opts.ActionsPath = opts.FilesPath
	if strings.Contains(string(opts.Name), "/") {
		parts := strings.Split(string(opts.Name), "/")
		opts.Model = inflect.Name(parts[len(parts)-1])
		opts.ActionsPath = inflect.Underscore(opts.Name.Resource())
	}

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	if err := g.Box(packr.NewBox("../gorm/templates")); err != nil {
		return g, errors.WithStack(err)
	}

	// define actions for resource
	opts.Actions = []string{"List", "Show", "New", "Create", "Edit", "Update", "Destroy"}

	// transform templates
	g.Transformer(gotools.TemplateTransformer(opts, template.FuncMap{}))

	// rename migrations
	g.Transformer(genny.NewTransformer(".fizz", func(f genny.File) (genny.File, error) {
		if !strings.Contains(f.Name(), ".fizz") {
			return f, nil
		}
		t := time.Now()
		p := opts.Model.PluralUnder()
		fN := strings.Replace(f.Name(), "migrations/migration", fmt.Sprintf("%s_create_%s", t.UTC().Format("20060102150405"), p), -1)
		return genny.NewFile(filepath.Join("migrations", fN), f), nil
	}))

	// rename resource actions
	g.Transformer(genny.NewTransformer(".go", func(f genny.File) (genny.File, error) {
		if !strings.Contains(f.Name(), "resource") {
			return f, nil
		}

		fN := strings.Replace(f.Name(), "resource", opts.FilesPath, -1)
		return genny.NewFile(fN, f), nil
	}))

	// rename view templates
	g.Transformer(genny.NewTransformer(".html", func(f genny.File) (genny.File, error) {
		if !strings.Contains(f.Name(), "model-view-") {
			return f, nil
		}

		fN := strings.Replace(f.Name(), "model-view-", fmt.Sprintf("%s/", opts.FilesPath), -1)
		return genny.NewFile(fN, f), nil
	}))

	// rename locales
	g.Transformer(genny.NewTransformer(".yaml", func(f genny.File) (genny.File, error) {
		if !strings.Contains(f.Name(), "resource") {
			return f, nil
		}

		fN := strings.Replace(f.Name(), "resource", opts.FilesPath, -1)
		return genny.NewFile(fN, f), nil
	}))

	return g, nil
}
