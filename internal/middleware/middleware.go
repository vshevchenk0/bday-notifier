package middleware

import "net/http"

type key string

type AuthMiddleware interface {
	Auth(h http.Handler) http.Handler
	GetUserIdContextKey() key
}
