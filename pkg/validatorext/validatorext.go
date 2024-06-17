package validatorext

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type FormattedError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func NewValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return validate
}

func FormatErrors(errors validator.ValidationErrors) []FormattedError {
	formattedErrors := make([]FormattedError, len(errors))
	for idx, err := range errors {
		var errorText string
		switch err.Tag() {
		case "required":
			errorText = "field is required"
		case "email":
			errorText = "must be email"
		case "uuid4":
			errorText = "must be UUIDv4"
		case "min":
			if err.Kind() == reflect.Int {
				errorText = fmt.Sprintf("should be greater than %s", err.Param())
			} else {
				errorText = fmt.Sprintf("minimum length is %s", err.Param())
			}
		case "max":
			if err.Kind() == reflect.Int {
				errorText = fmt.Sprintf("should be lower than %s", err.Param())
			} else {
				errorText = fmt.Sprintf("maximum length is %s", err.Param())
			}
		}
		formattedErrors[idx] = FormattedError{Field: err.Field(), Error: errorText}
	}
	return formattedErrors
}
