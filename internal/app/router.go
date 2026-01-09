package app

import (
	"compConfigurator/internal/config"
	"compConfigurator/internal/db"
	"compConfigurator/internal/httpapi/handlers/assemblies"
	"compConfigurator/internal/httpapi/handlers/auth"
	"compConfigurator/internal/httpapi/handlers/components/cpus"
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

	// deps
	userRepo := repo.NewUsersRepo(pool)
	authSvc := service.NewAuthService(userRepo, cfg.Auth.JWTSecret, cfg.Auth.JWTTTL)

	cpusRepo := repo.NewCPUsRepo(pool)
	cpusHandler := cpus.New(cpusRepo)

	assembliesRepo := repo.NewAssembliesRepo(pool)
	assemblySvc := service.NewAssemblyService(assembliesRepo)
	assembliesHandler := assemblies.New(assemblySvc)

	// handlers
	authHandler := auth.New(authSvc)
	meHandler := me.New(userRepo)

	// public routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	r.Route("/components", func(r chi.Router) {
		r.Route("/cpus", func(r chi.Router) {
			r.Get("/", cpusHandler.List)
			r.Get("/{id}", cpusHandler.Get)
		})
	})

	// protected
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(authSvc))
		r.Get("/me", meHandler.Me)

		r.Route("/assemblies", func(r chi.Router) {
			r.Post("/", assembliesHandler.Create)
			r.Get("/", assembliesHandler.List)
			r.Get("/{id}", assembliesHandler.Get)
		})
	})

	return r
}
