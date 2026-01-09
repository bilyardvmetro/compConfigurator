package app

import (
	"compConfigurator/internal/config"
	"compConfigurator/internal/db"
	"compConfigurator/internal/httpapi/handlers/auth"
	"compConfigurator/internal/httpapi/handlers/health"
	"compConfigurator/internal/httpapi/handlers/me"
	"compConfigurator/internal/httpapi/middleware"
	"compConfigurator/internal/repo"
	"compConfigurator/internal/service"

	"github.com/go-chi/chi/v5"
)

func NewRouter(cfg config.Config, pool *db.Pool) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recover)
	r.Use(middleware.Logger)

	// health
	healthHandler := health.New(pool, cfg.DB.HealthTimeout)
	r.Get("/health", healthHandler.Health)

	userRepo := repo.NewUsersRepo(pool)
	authSvc := service.NewAuthService(userRepo, cfg.Auth.JWTSecret, cfg.Auth.JWTTTL)

	// handlers
	authHandler := auth.New(authSvc)
	meHandler := me.New(userRepo)

	// routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// protected
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(authSvc))
		r.Get("/me", meHandler.Me)
	})

	return r
}
