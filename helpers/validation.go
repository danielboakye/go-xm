package helpers

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validation struct {
	validate *validator.Validate
}

// NewValidation returns a Validator instance
func NewValidation() (*Validation, error) {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &Validation{validate}, nil
}

// ValidateForm function return all validation errors in the associated form fields and struct
func (v *Validation) ValidateForm(i interface{}) error {

	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	if e, ok := err.(validator.ValidationErrors); ok {
		return errors.New(e[0].Field())
	}

	return nil
}
