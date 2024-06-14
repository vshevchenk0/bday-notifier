package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
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
	response := make([]ValidationError, len(errors))
	for idx, err := range errors {
		response[idx] = ValidationError{Field: err.Field(), Error: err.Error()}
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(responseBytes)
}
