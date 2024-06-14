package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

func initDocsFilesServer() http.Handler {
	return http.FileServer(http.Dir("./api/app"))
}

func NewRouter(
	authHandler *AuthHandler,
	subscriptionHandler *SubscriptionHandler,
	userHandler *UserHandler,
) http.Handler {
	r := chi.NewRouter()
	r.Mount("/auth", authHandler.router)
	r.Mount("/api/subscription", subscriptionHandler.router)
	r.Mount("/api/users", userHandler.router)
	r.Handle("/docs/*", http.StripPrefix("/docs/", initDocsFilesServer()))
	return r
}
