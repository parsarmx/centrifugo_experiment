package pkg

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	validate := validator.New()
	return &Validator{validate: validate}
}

func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}
