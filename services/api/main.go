package api

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"comp_config.com/services/api/adapters/aaa"
	"comp_config.com/services/api/adapters/rest"
	"comp_config.com/services/api/adapters/rest/middleware"
	"comp_config.com/services/api/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := mustMakeLogger(cfg.LogLevel)

	log.Info("starting server")
	log.Debug("debug messages are enabled")

	auth, err := aaa.New(cfg.TokenTTL, log)
	if err != nil {
		log.Error("cannot init auth adapter", "error", err)
		return
	}

	h, err := rest.NewClientUpdate(cfg.UpdateAddress, log)
	if err != nil {
		log.Error("cannot init update adapter", "error", err)
		return
	}

	hs, err := rest.NewSearchHandler(cfg.SearchAddress, log)
	if err != nil {
		log.Error("cannot init search adapter", "error", err)
		return
	}

	mux := http.NewServeMux()
	mux.Handle("POST /api/login", rest.NewLoginHandler(log, auth))
	// Ping
	mux.HandleFunc("GET /ping", h.Ping)

	// GPU routes
	setHandler(mux, "POST /api/gpus", h.AddGPU, auth)
	setHandler(mux, "PUT /api/gpus/{id}", h.UpdateGPU, auth)
	setHandler(mux, "DELETE /api/gpus/{id}", h.DeleteGPU, auth)

	// CPU Cooler routes
	setHandler(mux, "POST /api/cpu-coolers", h.AddCPUCooler, auth)
	setHandler(mux, "PUT /api/cpu-coolers/{id}", h.UpdateCPUCooler, auth)
	setHandler(mux, "DELETE /api/cpu-coolers/{id}", h.DeleteCPUCooler, auth)

	// Case routes
	setHandler(mux, "POST /api/cases", h.AddCase, auth)
	setHandler(mux, "PUT /api/cases/{id}", h.UpdateCase, auth)
	setHandler(mux, "DELETE /api/cases/{id}", h.DeleteCase, auth)

	// PSU routes
	setHandler(mux, "POST /api/psus", h.AddPSU, auth)
	setHandler(mux, "PUT /api/psus/{id}", h.UpdatePSU, auth)
	setHandler(mux, "DELETE /api/psus/{id}", h.DeletePSU, auth)

	// CPU Socket routes
	setHandler(mux, "POST /api/cpu-sockets", h.AddCPUSocket, auth)
	setHandler(mux, "PUT /api/cpu-sockets/{code}", h.UpdateCPUSocket, auth)
	setHandler(mux, "DELETE /api/cpu-sockets/{code}", h.DeleteCPUSocket, auth)

	// Motherboard Form Factor routes
	setHandler(mux, "POST /api/motherboard-form-factors", h.AddMotherboardFormFactor, auth)
	setHandler(mux, "PUT /api/motherboard-form-factors/{code}", h.UpdateMotherboardFormFactor, auth)
	setHandler(mux, "DELETE /api/motherboard-form-factors/{code}", h.DeleteMotherboardFormFactor, auth)

	// CPU routes
	setHandler(mux, "POST /api/cpus", h.AddCPU, auth)
	setHandler(mux, "PUT /api/cpus/{id}", h.UpdateCPU, auth)
	setHandler(mux, "DELETE /api/cpus/{id}", h.DeleteCPU, auth)

	// Motherboard routes
	setHandler(mux, "POST /api/motherboards", h.AddMotherboard, auth)
	setHandler(mux, "PUT /api/motherboards/{id}", h.UpdateMotherboard, auth)
	setHandler(mux, "DELETE /api/motherboards/{id}", h.DeleteMotherboard, auth)

	// RAM Kit routes
	setHandler(mux, "POST /api/ram-kits", h.AddRAMKit, auth)
	setHandler(mux, "PUT /api/ram-kits/{id}", h.UpdateRAMKit, auth)
	setHandler(mux, "DELETE /api/ram-kits/{id}", h.DeleteRAMKit, auth)

	// Storage Drive routes
	setHandler(mux, "POST /api/storage-drives", h.AddStorageDrive, auth)
	setHandler(mux, "PUT /api/storage-drives/{id}", h.UpdateStorageDrive, auth)
	setHandler(mux, "DELETE /api/storage-drives/{id}", h.DeleteStorageDrive, auth)

	// User routes
	setHandler(mux, "POST /api/users", h.AddUser, auth)
	setHandler(mux, "PUT /api/users/{id}", h.UpdateUser, auth)
	setHandler(mux, "DELETE /api/users/{id}", h.DeleteUser, auth)

	// Shop routes
	setHandler(mux, "POST /api/shops", h.AddShop, auth)
	setHandler(mux, "PUT /api/shops/{id}", h.UpdateShop, auth)
	setHandler(mux, "DELETE /api/shops/{id}", h.DeleteShop, auth)

	// Product Offer routes
	setHandler(mux, "POST /api/product-offers", h.AddProductOffer, auth)
	setHandler(mux, "PUT /api/product-offers/{id}", h.UpdateProductOffer, auth)
	setHandler(mux, "DELETE /api/product-offers/{id}", h.DeleteProductOffer, auth)

	// Assembly routes
	setHandler(mux, "POST /api/assemblies", h.AddAssembly, auth)
	setHandler(mux, "PUT /api/assemblies/{id}", h.UpdateAssembly, auth)
	setHandler(mux, "DELETE /api/assemblies/{id}", h.DeleteAssembly, auth)

	// Cooler Socket routes (many-to-many)
	setHandler(mux, "POST /api/cooler-sockets", h.AddCoolerSocket, auth)
	setHandler(mux, "DELETE /api/cooler-sockets", h.DeleteCoolerSocket, auth)
	setHandler(mux, "DELETE /api/coolers/{cooler_id}/sockets", h.DeleteCoolerSocketsByCooler, auth)

	// Case Form Factor Support routes (many-to-many)
	setHandler(mux, "POST /api/case-form-factor-supports", h.AddCaseFormFactorSupport, auth)
	setHandler(mux, "DELETE /api/case-form-factor-supports", h.DeleteCaseFormFactorSupport, auth)
	setHandler(mux, "DELETE /api/cases/{case_id}/form-factor-supports", h.DeleteCaseFormFactorSupportsByCase, auth)

	// Assembly RAM Kit routes (many-to-many)
	setHandler(mux, "POST /api/assembly-ram-kits", h.AddAssemblyRAMKit, auth)
	setHandler(mux, "DELETE /api/assembly-ram-kits", h.DeleteAssemblyRAMKit, auth)
	setHandler(mux, "DELETE /api/assemblies/{assembly_id}/ram-kits", h.DeleteAssemblyRAMKitsByAssembly, auth)

	// Assembly Drive routes (many-to-many)
	setHandler(mux, "POST /api/assembly-drives", h.AddAssemblyDrive, auth)
	setHandler(mux, "DELETE /api/assembly-drives", h.DeleteAssemblyDrive, auth)
	setHandler(mux, "DELETE /api/assemblies/{assembly_id}/drives", h.DeleteAssemblyDrivesByAssembly, auth)

	// Get endpoints
	mux.HandleFunc("GET /api/search/gpus/{id}", hs.GetGPU)
	mux.HandleFunc("GET /api/search/cpus/{id}", hs.GetCPU)
	mux.HandleFunc("GET /api/search/motherboards/{id}", hs.GetMotherboard)
	mux.HandleFunc("GET /api/search/ram-kits/{id}", hs.GetRAMKit)
	mux.HandleFunc("GET /api/search/psus/{id}", hs.GetPSU)
	mux.HandleFunc("GET /api/search/cases/{id}", hs.GetCase)
	mux.HandleFunc("GET /api/search/cpu-coolers/{id}", hs.GetCpuCooler)
	mux.HandleFunc("GET /api/search/storage-drives/{id}", hs.GetStorageDrive)
	mux.HandleFunc("GET /api/search/cpu-sockets/{code}", hs.GetCpuSocket)
	mux.HandleFunc("GET /api/search/motherboard-form-factors/{code}", hs.GetMotherboardFormFactor)
	mux.HandleFunc("GET /api/search/users/{id}", hs.GetUser)
	mux.HandleFunc("GET /api/search/shops/{id}", hs.GetShop)
	mux.HandleFunc("GET /api/search/product-offers/{id}", hs.GetProductOffer)
	mux.HandleFunc("GET /api/search/assemblies/{id}", hs.GetAssembly)

	// List endpoints
	mux.HandleFunc("GET /api/search/gpus", hs.ListGPUs)
	mux.HandleFunc("GET /api/search/cpus", hs.ListCPUs)
	mux.HandleFunc("GET /api/search/motherboards", hs.ListMotherboards)
	mux.HandleFunc("GET /api/search/ram-kits", hs.ListRAMKits)
	mux.HandleFunc("GET /api/search/psus", hs.ListPSUs)
	mux.HandleFunc("GET /api/search/cases", hs.ListCases)
	mux.HandleFunc("GET /api/search/cpu-coolers", hs.ListCpuCoolers)
	mux.HandleFunc("GET /api/search/storage-drives", hs.ListStorageDrives)
	mux.HandleFunc("GET /api/search/cpu-sockets", hs.ListCpuSockets)
	mux.HandleFunc("GET /api/search/motherboard-form-factors", hs.ListMotherboardFormFactors)
	mux.HandleFunc("GET /api/search/users", hs.ListUsers)
	mux.HandleFunc("GET /api/search/shops", hs.ListShops)
	mux.HandleFunc("GET /api/search/product-offers", hs.ListProductOffers)
	mux.HandleFunc("GET /api/search/assemblies", hs.ListAssemblies)
	mux.HandleFunc("GET /api/users/{user_id}/assemblies", hs.ListAssembliesByUser)

	// Search endpoints
	mux.HandleFunc("POST /api/search/components", hs.SearchComponents)
	mux.HandleFunc("POST /api/search/assemblies/search", hs.SearchAssemblies)

	// Compatibility endpoints
	mux.HandleFunc("GET /api/assemblies/{assembly_id}/compatibility", hs.CheckCompatibility)
	mux.HandleFunc("POST /api/search/compatible-components", hs.GetCompatibleComponents)
	mux.HandleFunc("POST /api/search/recommendations", hs.RecommendComponents)

	// Price endpoints
	mux.HandleFunc("GET /api/components/{component_type}/{component_id}/prices", hs.GetComponentPrices)
	mux.HandleFunc("GET /api/components/{component_type}/{component_id}/best-offers", hs.GetBestOffers)

	server := http.Server{
		Addr:        cfg.HTTPConfig.Address,
		ReadTimeout: cfg.HTTPConfig.Timeout,
		Handler:     mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Debug("shutting down server")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Error("erroneous shutdown", "error", err)
		}
	}()

	log.Info("Running HTTP server", "address", cfg.HTTPConfig.Address)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error("server closed unexpectedly", "error", err)
			return
		}
	}
}

func setHandler(mux *http.ServeMux, endpoint string, f func(w http.ResponseWriter, r *http.Request), verifire middleware.TokenVerifier) {
	mux.Handle(endpoint, middleware.Auth(func(w http.ResponseWriter, r *http.Request) { f(w, r) }, verifire))
}

func mustMakeLogger(logLevel string) *slog.Logger {
	var level slog.Level
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "ERROR":
		level = slog.LevelError
	default:
		panic("unknown log level: " + logLevel)
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
