package middleware

import (
	domainErr "compConfigurator/internal/domain/errors"
	"compConfigurator/internal/httpapi/response"
	"compConfigurator/internal/service"
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const (
	UserIDKey ctxKey = "user_id"
	RoleKey   ctxKey = "role"
)

func RequireAuth(auth *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				response.Fail(w, http.StatusUnauthorized, "missing bearer token")
				return
			}

			tokenString := strings.TrimPrefix(h, "Bearer ")
			claims, err := auth.ParseToken(tokenString)
			if err != nil {
				response.Fail(w, http.StatusUnauthorized, domainErr.ErrUnauthorized.Error())
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) (int64, bool) {
	v := r.Context().Value(UserIDKey)
	id, ok := v.(int64)
	return id, ok
}

func GetRole(r *http.Request) (string, bool) {
	v := r.Context().Value(RoleKey)
	role, ok := v.(string)
	return role, ok
}
