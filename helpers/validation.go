package helpers

import (
	"errors"
	"reflect"
	"strings"

	"github.com/danielboakye/go-xm/helpers/consts"
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

	err := validate.RegisterValidation("company-types", validateCompanyType)
	if err != nil {
		return nil, err
	}

	return &Validation{validate}, nil
}

func validateCompanyType(fl validator.FieldLevel) bool {
	docType := fl.Field().String()

	contains := SliceContains(consts.CompanyTypes, docType)
	return contains
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
