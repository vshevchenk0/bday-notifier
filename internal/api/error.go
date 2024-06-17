package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/vshevchenk0/bday-notifier/pkg/validatorext"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidationErrorResponse struct {
	Errors []validatorext.FormattedError `json:"errors"`
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, error error) {
	response := ErrorResponse{Error: error.Error()}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(responseBytes)
}

func WriteValidationErrorResponse(w http.ResponseWriter, statusCode int, errors validator.ValidationErrors) {
	formattedErrors := validatorext.FormatErrors(errors)
	response := ValidationErrorResponse{Errors: formattedErrors}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(responseBytes)
}
