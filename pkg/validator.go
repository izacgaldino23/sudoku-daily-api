package pkg

import (
	"errors"
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

	var validationErrs validator.ValidationErrors
	ok := errors.As(err, &validationErrs)
	if !ok {
		return err
	}

	errorList := make(ValidationErrors, 0, len(validationErrs))
	for _, e := range validationErrs {
		errorList = append(errorList, ValidationError{
			Field:   e.Field(),
			Message: formatValidationError(e),
		})
	}

	return errorList
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
	case "oneof":
		return fmt.Sprintf("must be one of %s", e.Param())
	default:
		return fmt.Sprintf("failed %s validation", e.Tag())
	}
}
