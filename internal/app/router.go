package app

import (
	"compConfigurator/internal/config"
	"compConfigurator/internal/db"
	"compConfigurator/internal/httpapi/handlers/health"
	"compConfigurator/internal/httpapi/middleware"

	"github.com/go-chi/chi/v5"
)

func NewRouter(cfg config.Config, pool *db.Pool) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recover)
	r.Use(middleware.Logger)

	healthHandler := health.New(pool, cfg.DB.HealthTimeout)

	r.Get("/health", healthHandler.Health)

	return r
}
