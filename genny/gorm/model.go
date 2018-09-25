package gorm

import (
	"github.com/markbates/inflect"
)

type Model struct {
	Package               string
	Imports               []string
	Name                  inflect.Name
	Attributes            []Attribute
	ValidatableAttributes []Attribute

	HasNulls  bool
	HasUUID   bool
	HasSlices bool
	HasID     bool
}

func NewModel(name string) Model {
	m := Model{
		Package: "models",
		Imports: []string{"time", "encoding/json", "github.com/gobuffalo/uuid", "github.com/jinzhu/gorm", "github.com/gobuffalo/validate", "github.com/pkg/errors"},
		Name:    inflect.Name(name),
		Attributes: []Attribute{
			{Name: inflect.Name("id"), OriginalType: "uuid", GoType: "uuid.UUID", FizzType: "uuid", Primary: true},
			{Name: inflect.Name("created_at"), OriginalType: "time.Time", GoType: "time.Time", FizzType: "timestamp"},
			{Name: inflect.Name("updated_at"), OriginalType: "time.Time", GoType: "time.Time", FizzType: "timestamp"},
			{Name: inflect.Name("deleted_at"), OriginalType: "*time.Time", GoType: "*time.Time", FizzType: "timestamp", Nullable: true},
		},
		ValidatableAttributes: []Attribute{},
	}
	return m
}

func (m *Model) addAttribute(a Attribute) {
	if a.Name == "id" {
		// No need to create a default ID
		m.HasID = true
		// Ensure ID is the first attribute
		m.Attributes = append([]Attribute{a}, m.Attributes...)
	} else {
		m.Attributes = append(m.Attributes, a)
	}

	if a.Nullable {
		return
	}

	if a.IsValidable() {
		if a.GoType == "time.Time" {
			a.GoType = "Time"
		}
		m.ValidatableAttributes = append(m.ValidatableAttributes, a)
	}
}

func (m *Model) ParseAttributes(attrs ...Prop) {
	for _, a := range attrs {
		m.addAttribute(newAttribute(a, m))
	}
}
