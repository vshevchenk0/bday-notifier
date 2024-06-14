package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/vshevchenk0/bday-greeter/internal/service"
)

const UserIdKey key = "userId"

type authMiddleware struct {
	authService      service.AuthService
	userIdContextKey key
}

func NewAuthMiddleware(authService service.AuthService) *authMiddleware {
	return &authMiddleware{
		authService:      authService,
		userIdContextKey: UserIdKey,
	}
}

func (m *authMiddleware) Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("no token provided"))
			return
		}
		authHeader := r.Header["Authorization"][0]
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" || headerParts[1] == "" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("malformed token"))
			return
		}
		userId, err := m.authService.VerifyToken(r.Context(), headerParts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("token expired"))
			return
		}
		ctx := context.WithValue(r.Context(), m.userIdContextKey, userId)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *authMiddleware) GetUserIdContextKey() key {
	return m.userIdContextKey
}
