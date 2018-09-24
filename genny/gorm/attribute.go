package gorm

import (
	"fmt"
	"strings"

	"github.com/markbates/inflect"
)

type Attribute struct {
	Name         inflect.Name
	OriginalType string
	GoType       string
	Nullable     bool
}

func (a Attribute) String() string {
	t := ""
	if a.GoType == "uuid.UUID" {
		t = "type:char(36);"
	}

	return fmt.Sprintf("\t%s %s `%s:\"%s\" gorm:\"column:%s;%s\"`", a.Name.Camel(), a.GoType, "json", a.Name.Underscore(), a.Name.Underscore(), t)
}

func (a Attribute) IsValidable() bool {
	return a.GoType == "string" || a.GoType == "time.Time" || a.GoType == "int"
}

func newAttribute(p Prop, model *Model) Attribute {
	// col := strings.Split(base, ":")
	// if len(col) == 1 {
	// 	col = append(col, "string")
	// }

	nullable := strings.HasPrefix(p.Type, "nulls.")
	if !model.HasNulls && nullable {
		model.HasNulls = true
		model.Imports = append(model.Imports, "github.com/gobuffalo/pop/nulls")
	}

	got := colType(p.Type)
	// if len(col) > 2 {
	// 	got = col[2]
	// }
	a := Attribute{
		Name:         p.Name,
		OriginalType: p.Type,
		GoType:       got,
		Nullable:     nullable,
	}

	return a
}

func colType(s string) string {
	switch strings.ToLower(s) {
	case "text":
		return "string"
	case "time", "timestamp", "datetime":
		return "time.Time"
	case "nulls.text":
		return "nulls.String"
	case "uuid":
		return "uuid.UUID"
	case "json", "jsonb":
		return "slices.Map"
	case "[]string":
		return "slices.String"
	case "[]int":
		return "slices.Int"
	case "slices.float", "[]float", "[]float32", "[]float64":
		return "slices.Float"
	case "decimal", "float":
		return "float64"
	case "[]byte", "blob":
		return "[]byte"
	default:
		return s
	}
}
