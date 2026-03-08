package pkg

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	if len(v) == 1 {
		return fmt.Sprintf("%s: %s", v[0].Field, v[0].Message)
	}
	return fmt.Sprintf("%d validation errors", len(v))
}

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	errors := make(ValidationErrors, 0, len(validationErrs))
	for _, e := range validationErrs {
		errors = append(errors, ValidationError{
			Field:   e.Field(),
			Message: formatValidationError(e),
		})
	}

	return errors
}

func formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return fmt.Sprintf("must be at least %s characters", e.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", e.Param())
	case "len":
		return fmt.Sprintf("must be %s characters", e.Param())
	case "numeric":
		return "must be numeric"
	default:
		return fmt.Sprintf("failed %s validation", e.Tag())
	}
}
