package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/vshevchenk0/bday-notifier/internal/middleware"
	"github.com/vshevchenk0/bday-notifier/internal/service"
	"github.com/vshevchenk0/bday-notifier/pkg/validatorext"
)

type SubscriptionHandler struct {
	subscriptionService service.SubscriptionService
	authMiddleware      middleware.AuthMiddleware
	validate            *validator.Validate
	router              chi.Router
}

type createSubscriptionRequestBody struct {
	UserId           string `json:"user_id" validate:"required,uuid4"`
	NotifyBeforeDays int    `json:"notify_before_days" validate:"required,min=1,max=7"`
}

type deleteSubscriptionRequestBody struct {
	UserId string `json:"user_id" validate:"required,uuid4"`
}

func NewSubscriptionHandler(
	subscriptionService service.SubscriptionService,
	authMiddleware middleware.AuthMiddleware,
) *SubscriptionHandler {
	handler := &SubscriptionHandler{
		subscriptionService: subscriptionService,
		authMiddleware:      authMiddleware,
		validate:            validatorext.NewValidator(),
		router:              chi.NewRouter(),
	}
	handler.initRoutes()
	return handler
}

func (h *SubscriptionHandler) createSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var body createSubscriptionRequestBody
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

	subscriberId := r.Context().Value(h.authMiddleware.GetUserIdContextKey()).(string)

	err = h.subscriptionService.CreateSubscription(r.Context(), body.UserId, subscriberId, body.NotifyBeforeDays)
	if errors.Is(err, service.ErrUserNotFound) {
		errText := fmt.Errorf("user you are trying to subscribe to was not found")
		WriteErrorResponse(w, http.StatusNotFound, errText)
		return
	}
	if errors.Is(err, service.ErrDuplicateSubscription) {
		errText := fmt.Errorf("subscription already exists")
		WriteErrorResponse(w, http.StatusBadRequest, errText)
		return
	}
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("subscription created"))
}

func (h *SubscriptionHandler) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var body deleteSubscriptionRequestBody
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

	subscriberId := r.Context().Value(h.authMiddleware.GetUserIdContextKey()).(string)
	err = h.subscriptionService.DeleteSubscription(r.Context(), body.UserId, subscriberId)
	if errors.Is(err, service.ErrOperationResultUnknown) {
		errText := fmt.Errorf("deletion result unknown. check your subscriptions and try again if needed")
		WriteErrorResponse(w, http.StatusInternalServerError, errText)
		return
	}
	if errors.Is(err, service.ErrSubscriptionNotFound) {
		errText := fmt.Errorf("subscription you are trying to delete was not found")
		WriteErrorResponse(w, http.StatusNotFound, errText)
		return
	}
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("subscription deleted"))
}

func (h *SubscriptionHandler) initRoutes() {
	h.router.With(h.authMiddleware.Auth).Post("/", h.createSubscription)
	h.router.With(h.authMiddleware.Auth).Delete("/", h.deleteSubscription)
}
