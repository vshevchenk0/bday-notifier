package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/vshevchenk0/bday-greeter/internal/middleware"
	"github.com/vshevchenk0/bday-greeter/internal/service"
)

type UserHandler struct {
	userService    service.UserService
	authMiddleware middleware.AuthMiddleware
	router         chi.Router
}

func NewUserHandler(
	userService service.UserService,
	authMiddleware middleware.AuthMiddleware,
) *UserHandler {
	handler := &UserHandler{
		userService:    userService,
		authMiddleware: authMiddleware,
		router:         chi.NewRouter(),
	}
	handler.initRoutes()
	return handler
}

func (h *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	userId := r.Context().Value(h.authMiddleware.GetUserIdContextKey()).(string)
	users, err := h.userService.FindAllUsers(r.Context(), userId)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	if len(users) == 0 {
		errText := fmt.Errorf("no users found")
		WriteErrorResponse(w, http.StatusNotFound, errText)
		return
	}

	response, err := json.Marshal(users)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *UserHandler) getUsersSubscribedTo(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	userId := r.Context().Value(h.authMiddleware.GetUserIdContextKey()).(string)
	users, err := h.userService.FindUsersSubscribedTo(r.Context(), userId)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	if len(users) == 0 {
		errText := fmt.Errorf("no subscriptions found")
		WriteErrorResponse(w, http.StatusNotFound, errText)
		return
	}

	response, err := json.Marshal(users)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
}

func (h *UserHandler) initRoutes() {
	h.router.With(h.authMiddleware.Auth).Get("/", h.getUsers)
	h.router.With(h.authMiddleware.Auth).Get("/subscriptions", h.getUsersSubscribedTo)
}
