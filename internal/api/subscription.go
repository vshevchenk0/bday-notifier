package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/vshevchenk0/bday-greeter/internal/middleware"
	"github.com/vshevchenk0/bday-greeter/internal/service"
)

type SubscriptionHandler struct {
	subscriptionService service.SubscriptionService
	authMiddleware      middleware.AuthMiddleware
	validator           *validator.Validate
	router              chi.Router
}

type createSubscriptionRequestBody struct {
	UserId           string `json:"user_id" validate:"uuid4,required"`
	NotifyBeforeDays int    `json:"notify_before_days" validate:"min=0,max=7,required"`
}

type deleteSubscriptionRequestBody struct {
	UserId string `json:"user_id" validate:"uuid4,required"`
}

func NewSubscriptionHandler(
	subscriptionService service.SubscriptionService,
	authMiddleware middleware.AuthMiddleware,
) *SubscriptionHandler {
	handler := &SubscriptionHandler{
		subscriptionService: subscriptionService,
		authMiddleware:      authMiddleware,
		validator:           validator.New(validator.WithRequiredStructEnabled()),
		router:              chi.NewRouter(),
	}
	handler.initRoutes()
	return handler
}

func (h *SubscriptionHandler) createSubscription(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var body createSubscriptionRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	err := h.validator.Struct(body)
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
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	err := h.validator.Struct(body)
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
