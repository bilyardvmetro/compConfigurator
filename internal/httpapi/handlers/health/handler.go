package health

import (
	"compConfigurator/internal/db"
	"context"
	"net/http"
	"time"
)

type Handler struct {
	pool          *db.Pool
	healthTimeout time.Duration
}

func New(pool *db.Pool, healthTimeout time.Duration) *Handler {
	return &Handler{pool: pool, healthTimeout: healthTimeout}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.healthTimeout)
	defer cancel()

	if err := h.pool.Ping(ctx); err != nil {
		http.Error(w, "db not ready", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
