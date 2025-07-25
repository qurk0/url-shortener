package validator

import (
	"context"
	"errors"

	"github.com/go-playground/validator"
)

var global *validator.Validate

const (
	ErrInvalidFormat     = "invalid format"
	ErrFieldRequired     = "field required"
	ErrUnknownValidation = "unknown validation error"
)

func init() {
	SetValidator(New())
}

func New() *validator.Validate {
	v := validator.New()
	// _ = v.RegisterValidation("")

	return v
}

func SetValidator(v *validator.Validate) {
	global = v
}

func Validator() *validator.Validate {
	return global
}

func Validate(ctx context.Context, structure any) error {
	return parseValidationError(Validator().StructCtx(ctx, structure))
}

func parseValidationError(err error) error {
	if err == nil {
		return nil
	}

	vErrs, ok := err.(validator.ValidationErrors)
	if !ok || len(vErrs) == 0 {
		return nil
	}

	validationError := vErrs[0]
	var validationErrorDescription string
	switch validationError.Tag() {
	case "tag":
		validationErrorDescription = ErrInvalidFormat
	case "required":
		validationErrorDescription = ErrFieldRequired
	default:
		validationErrorDescription = ErrUnknownValidation
	}

	return errors.New(validationErrorDescription + ": " + validationError.Namespace())
}
