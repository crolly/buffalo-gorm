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
	FizzType     string
	Primary      bool
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

	if !model.HasNulls && p.Nullable {
		model.HasNulls = true
		model.Imports = append(model.Imports, "github.com/gobuffalo/pop/nulls")
	}

	var got string
	if p.Nullable {
		// append nulls. for nullable
		got = colType(fmt.Sprintf("%s.%s", "nulls", p.Type))
	} else {
		got = colType(p.Type)
	}
	ft := fizzColType(p.Type)
	a := Attribute{
		Name:         p.Name,
		OriginalType: p.Type,
		GoType:       got,
		FizzType:     ft,
		Nullable:     p.Nullable,
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
	case "nulls.uuid":
		return "nulls.UUID"
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

func fizzColType(s string) string {
	switch strings.ToLower(s) {
	case "int":
		return "integer"
	case "time", "datetime":
		return "timestamp"
	case "uuid.uuid", "uuid":
		return "uuid"
	case "nulls.float32", "nulls.float64":
		return "float"
	case "slices.string", "slices.uuid", "[]string":
		return "varchar[]"
	case "slices.float", "[]float", "[]float32", "[]float64":
		return "numeric[]"
	case "slices.int":
		return "int[]"
	case "slices.map":
		return "jsonb"
	case "float32", "float64", "float":
		return "decimal"
	case "blob", "[]byte":
		return "blob"
	default:
		if strings.HasPrefix(s, "nulls.") {
			return fizzColType(strings.Replace(s, "nulls.", "", -1))
		}
		return strings.ToLower(s)
	}
}
