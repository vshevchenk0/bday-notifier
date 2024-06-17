package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/vshevchenk0/bday-notifier/internal/service"
	"github.com/vshevchenk0/bday-notifier/pkg/validatorext"
)

type AuthHandler struct {
	authService service.AuthService
	validate    *validator.Validate
	router      chi.Router
}

type signUpRequestBody struct {
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=7"`
	Name         string `json:"name" validate:"required,min=1"`
	Surname      string `json:"surname" validate:"required,min=1"`
	BirthdayDate string `json:"birthday_date" validate:"required"`
}

type signInRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=7"`
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	handler := &AuthHandler{
		authService: authService,
		validate:    validatorext.NewValidator(),
		router:      chi.NewRouter(),
	}
	handler.initRoutes()
	return handler
}

func (h *AuthHandler) signUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var body signUpRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if errors.Is(err, io.EOF) {
		errText := fmt.Errorf("body is required")
		WriteErrorResponse(w, http.StatusBadRequest, errText)
		return
	}
	if err != nil {
		errText := fmt.Errorf("failed to decode request body")
		WriteErrorResponse(w, http.StatusInternalServerError, errText)
		return
	}

	err = h.validate.Struct(body)
	if _, ok := err.(*validator.InvalidValidationError); ok {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		WriteValidationErrorResponse(w, http.StatusBadRequest, validationErrors)
		return
	}

	birthdayDate, err := time.Parse(time.DateOnly, body.BirthdayDate)
	if err != nil {
		errText := fmt.Errorf("wrong date format, please use this format: %s", time.DateOnly)
		WriteErrorResponse(w, http.StatusBadRequest, errText)
		return
	}

	token, err := h.authService.SignUp(r.Context(), body.Email, body.Password, body.Name, body.Surname, birthdayDate)
	if errors.Is(err, service.ErrDuplicateUser) {
		errText := fmt.Errorf("user with this email already exists")
		WriteErrorResponse(w, http.StatusBadRequest, errText)
		return
	}
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(token)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(response)
}

func (h *AuthHandler) signIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var body signInRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if errors.Is(err, io.EOF) {
		errText := fmt.Errorf("body is required")
		WriteErrorResponse(w, http.StatusBadRequest, errText)
		return
	}
	if err != nil {
		errText := fmt.Errorf("failed to decode request body")
		WriteErrorResponse(w, http.StatusInternalServerError, errText)
		return
	}

	err = h.validate.Struct(body)
	if _, ok := err.(*validator.InvalidValidationError); ok {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		WriteValidationErrorResponse(w, http.StatusBadRequest, validationErrors)
		return
	}

	token, err := h.authService.SignIn(r.Context(), body.Email, body.Password)
	if errors.Is(err, service.ErrUserNotFound) {
		errText := fmt.Errorf("user with this email was not found")
		WriteErrorResponse(w, http.StatusNotFound, errText)
		return
	}
	if errors.Is(err, service.ErrInvalidPassword) {
		errText := fmt.Errorf("invalid password")
		WriteErrorResponse(w, http.StatusUnauthorized, errText)
		return
	}
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	response, err := json.Marshal(token)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(response)
}

func (h *AuthHandler) initRoutes() {
	h.router.Post("/signup", h.signUp)
	h.router.Post("/signin", h.signIn)
}
