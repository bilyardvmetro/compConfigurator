package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const RequestIDkey ctxKey = "request_id"
const HeaderRequestID = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(HeaderRequestID)
		if rid == "" {
			rid = uuid.NewString()
		}

		w.Header().Set(HeaderRequestID, rid)
		ctx := context.WithValue(r.Context(), RequestIDkey, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
