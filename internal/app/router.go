package app

import (
	"compConfigurator/internal/config"
	"compConfigurator/internal/db"
	"compConfigurator/internal/httpapi/handlers/assemblies"
	"compConfigurator/internal/httpapi/handlers/auth"
	"compConfigurator/internal/httpapi/handlers/components"
	"compConfigurator/internal/httpapi/handlers/health"
	"compConfigurator/internal/httpapi/handlers/me"
	"compConfigurator/internal/httpapi/handlers/offers"
	"compConfigurator/internal/httpapi/handlers/public"
	"compConfigurator/internal/httpapi/middleware"
	"compConfigurator/internal/repo"
	"compConfigurator/internal/service"
	"compConfigurator/internal/service/pricing"

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

	assembliesRepo := repo.NewAssembliesRepo(pool)
	assemblySvc := service.NewAssemblyService(assembliesRepo)

	partsRepo := repo.NewAssemblyPartsRepo(pool)
	partsSvc := service.NewAssemblyPartsService(assembliesRepo, partsRepo)

	detailsRepo := repo.NewAssemblyDetailsRepo(pool)
	detailsSvc := service.NewAssemblyDetailsService(assembliesRepo, detailsRepo)

	cpusRepo := repo.NewCPUsRepo(pool)
	ramRepo := repo.NewRamKitsRepo(pool)
	drvRepo := repo.NewDrivesRepo(pool)
	mbRepo := repo.NewMotherboardsRepo(pool)
	casesRepo := repo.NewCasesRepo(pool)
	psusRepo := repo.NewPSUsRepo(pool)
	coolersRepo := repo.NewCPUCoolersRepo(pool)
	gpuRepo := repo.NewGPUsRepo(pool)

	publicRepo := repo.NewPublicAssembliesRepo(pool)
	cloneRepo := repo.NewCloneRepo(pool)

	offersRepo := repo.NewOffersRepo(pool)
	offersHandler := offers.New(offersRepo)

	componentsRepo := repo.NewAssemblyComponentsRepo(pool)
	pricingSvc := pricing.New(componentsRepo, detailsRepo, offersRepo)

	// handlers
	assembliesHandler := assemblies.New(assemblySvc, partsSvc, detailsSvc, pricingSvc)
	authHandler := auth.New(authSvc)
	meHandler := me.New(userRepo)
	componentsHandler := components.New(
		cpusRepo,
		ramRepo,
		drvRepo,
		mbRepo,
		casesRepo,
		psusRepo,
		coolersRepo,
		gpuRepo,
	)

	publicHandler := public.New(publicRepo, cloneRepo, detailsSvc)

	// public routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	r.Get("/offers", offersHandler.List)
	r.Get("/offers/best", offersHandler.Best)

	r.Route("/components", func(r chi.Router) {
		r.Get("/cpus", componentsHandler.ListCPUs)
		r.Get("/motherboards", componentsHandler.ListMotherboards)
		r.Get("/gpus", componentsHandler.ListGPUs)
		r.Get("/psus", componentsHandler.ListPSUs)
		r.Get("/cases", componentsHandler.ListCases)
		r.Get("/cpu-coolers", componentsHandler.ListCPUCoolers)
		r.Get("/ram-kits", componentsHandler.ListRamKits)
		r.Get("/drives", componentsHandler.ListDrives)
	})

	r.Route("/public", func(r chi.Router) {
		r.Get("/assemblies", publicHandler.ListPublicAssemblies)
		r.Get("/assemblies/{id}", publicHandler.GetPublicAssembly)
		r.Get("/assemblies/{id}/details", publicHandler.PublicDetails)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAuth(authSvc))
			r.Post("/assemblies/{id}/clone", publicHandler.ClonePublicAssembly)
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
			r.Put("/{id}", assembliesHandler.Update)

			// RAM
			r.Post("/{id}/ram/{ramKitId}", assembliesHandler.AddRamKit)
			r.Delete("/{id}/ram/{ramKitId}", assembliesHandler.RemoveRamKit)

			// Drives
			r.Post("/{id}/drives/{driveId}", assembliesHandler.AddDrive)
			r.Delete("/{id}/drives/{driveId}", assembliesHandler.RemoveDrive)

			// all assembly
			r.Get("/{id}/details", assembliesHandler.Details)

			r.Get("/{id}/pricing", assembliesHandler.Pricing)
		})
	})

	return r
}
