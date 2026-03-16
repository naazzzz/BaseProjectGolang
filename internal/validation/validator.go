package validation

import (
	"BaseProjectGolang/internal/config"

	"github.com/go-playground/validator/v10"
	"github.com/rotisserie/eris"
	"github.com/soner3/flora"
)

type IValidator interface {
	Validate(data interface{}) (err error)
}

type XValidator struct {
	flora.Component
	validator *validator.Validate
}

func NewXValidator(
	_ *config.Config,
) *XValidator {
	valid := &XValidator{
		validator: validator.New(),
	}

	return valid
}

func (v *XValidator) Validate(data interface{}) (err error) {
	if err = v.validator.Struct(data); err != nil {
		return eris.Wrap(err, "Validation errors")
	}

	return
}
