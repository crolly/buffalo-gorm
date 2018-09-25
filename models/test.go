package models

import (
	"time"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/validate/validators"
	)

type test struct {
		ID uuid.UUID `json:"id" gorm:"column:id;type:char(36);"`
		CreatedAt time.Time `json:"created_at" gorm:"column:created_at;"`
		UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;"`
		DeletedAt *time.Time `json:"deleted_at" gorm:"column:deleted_at;"`
		Name string `json:"name" gorm:"column:name;"`
		Birthday time.Time `json:"birthday" gorm:"column:birthday;"`
		AddressID nulls.UUID `json:"address_id" gorm:"column:address_id;"`
	}

// String is not required by pop and may be deleted
func (t test) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Tests is not required by pop and may be deleted
type Tests []test

// String is not required by pop and may be deleted
func (t Tests) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *test) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.Name, Name: "Name"},
		&validators.TimeIsPresent{Field: t.Birthday, Name: "Birthday"},
		), nil
	}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *test) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *test) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}