package core

import (
	"context"
	"log/slog"
)

type Service struct {
	log *slog.Logger
	db  DB
}

func NewService(log *slog.Logger, db DB) (*Service, error) {
	return &Service{
		log: log,
		db:  db,
	}, nil
}

// GPU methods
func (s *Service) GetGPU(ctx context.Context, id int) (*GPU, error) {
	s.log.Debug("Getting GPU", "id", id)
	return s.db.GetGPU(ctx, id)
}

func (s *Service) ListGPUs(ctx context.Context, opts *ListOptions) (*ListResult[GPU], error) {
	s.log.Debug("Listing GPUs", "page", opts.Page, "page_size", opts.PageSize)
	return s.db.ListGPUs(ctx, opts)
}

// CPU Cooler methods
func (s *Service) GetCpuCooler(ctx context.Context, id int) (*CpuCooler, error) {
	s.log.Debug("Getting CPU cooler", "id", id)
	return s.db.GetCpuCooler(ctx, id)
}

func (s *Service) ListCpuCoolers(ctx context.Context, opts *ListOptions) (*ListResult[CpuCooler], error) {
	s.log.Debug("Listing CPU coolers", "page", opts.Page, "page_size", opts.PageSize)
	return s.db.ListCpuCoolers(ctx, opts)
}

// Case methods
func (s *Service) GetCase(ctx context.Context, id int) (*Case, error) {
	s.log.Debug("Getting case", "id", id)
	return s.db.GetCase(ctx, id)
}

func (s *Service) ListCases(ctx context.Context, opts *ListOptions) (*ListResult[Case], error) {
	s.log.Debug("Listing cases", "page", opts.Page, "page_size", opts.PageSize)
	return s.db.ListCases(ctx, opts)
}

// PSU methods
func (s *Service) GetPSU(ctx context.Context, id int) (*PSU, error) {
	s.log.Debug("Getting PSU", "id", id)
	return s.db.GetPSU(ctx, id)
}

func (s *Service) ListPSUs(ctx context.Context, opts *ListOptions) (*ListResult[PSU], error) {
	s.log.Debug("Listing PSUs", "page", opts.Page, "page_size", opts.PageSize)
	return s.db.ListPSUs(ctx, opts)
}

// CPU Socket methods
func (s *Service) GetCpuSocket(ctx context.Context, code string) (*CpuSocket, error) {
	s.log.Debug("Getting CPU socket", "code", code)
	return s.db.GetCpuSocket(ctx, code)
}

func (s *Service) ListCpuSockets(ctx context.Context, opts *ListOptions) (*ListResult[CpuSocket], error) {
	s.log.Debug("Listing CPU sockets", "page", opts.Page, "page_size", opts.PageSize)
	return s.db.ListCpuSockets(ctx, opts)
}

// Motherboard Form Factor methods
func (s *Service) GetMotherboardFormFactor(ctx context.Context, code string) (*MotherboardFormFactor, error) {
	s.log.Debug("Getting motherboard form factor", "code", code)
	return s.db.GetMotherboardFormFactor(ctx, code)
}

func (s *Service) ListMotherboardFormFactors(ctx context.Context, opts *ListOptions) (*ListResult[MotherboardFormFactor], error) {
	s.log.Debug("Listing motherboard form factors", "page", opts.Page, "page_size", opts.PageSize)
	return s.db.ListMotherboardFormFactors(ctx, opts)
}

// CPU methods
func (s *Service) GetCPU(ctx context.Context, id int) (*CPU, error) {
	s.log.Debug("Getting CPU", "id", id)
	return s.db.GetCPU(ctx, id)
}

func (s *Service) ListCPUs(ctx context.Context, opts *ListOptions) (*ListResult[CPU], error) {
	s.log.Debug("Listing CPUs", "page", opts.Page, "page_size", opts.PageSize)
	return s.db.ListCPUs(ctx, opts)
}

// Motherboard methods
func (s *Service) GetMotherboard(ctx context.Context, id int) (*Motherboard, error) {
	s.log.Debug("Getting motherboard", "id", id)
	return s.db.GetMotherboard(ctx, id)
}

func (s *Service) ListMotherboards(ctx context.Context, opts *ListOptions) (*ListResult[Motherboard], error) {
	s.log.Debug("Listing motherboards", "page", opts.Page, "page_size", opts.PageSize)
	return s.db.ListMotherboards(ctx, opts)
}

// RAM Kit methods
func (s *Service) GetRAMKit(ctx context.Context, id int) (*RAMKit, error) {
	s.log.Debug("Getting RAM kit", "id", id)
	return s.db.GetRAMKit(ctx, id)
}

func (s *Service) ListRAMKits(ctx context.Context, opts *ListOptions) (*ListResult[RAMKit], error) {
	s.log.Debug("Listing RAM kits", "page", opts.Page, "page_size", opts.PageSize)
	return s.db.ListRAMKits(ctx, opts)
}

// Storage Drive methods
func (s *Service) GetStorageDrive(ctx context.Context, id int) (*StorageDrive, error) {
	s.log.Debug("Getting storage drive", "id", id)
	return s.db.GetStorageDrive(ctx, id)
}

func (s *Service) ListStorageDrives(ctx context.Context, opts *ListOptions) (*ListResult[StorageDrive], error) {
	s.log.Debug("Listing storage drives", "page", opts.Page, "page_size", opts.PageSize)
	return s.db.ListStorageDrives(ctx, opts)
}

// User methods
func (s *Service) GetUser(ctx context.Context, id int) (*User, error) {
	s.log.Debug("Getting user", "id", id)
	return s.db.GetUser(ctx, id)
}

func (s *Service) ListUsers(ctx context.Context, opts *ListOptions) (*ListResult[User], error) {
	s.log.Debug("Listing users", "page", opts.Page, "page_size", opts.PageSize)
	return s.db.ListUsers(ctx, opts)
}

// Shop methods
func (s *Service) GetShop(ctx context.Context, id int) (*Shop, error) {
	s.log.Debug("Getting shop", "id", id)
	return s.db.GetShop(ctx, id)
}

func (s *Service) ListShops(ctx context.Context, opts *ListOptions) (*ListResult[Shop], error) {
	s.log.Debug("Listing shops", "page", opts.Page, "page_size", opts.PageSize)
	return s.db.ListShops(ctx, opts)
}

// Product Offer methods
func (s *Service) GetProductOffer(ctx context.Context, id int) (*ProductOffer, error) {
	s.log.Debug("Getting product offer", "id", id)
	return s.db.GetProductOffer(ctx, id)
}

func (s *Service) ListProductOffers(ctx context.Context, opts *ListOptions) (*ListResult[ProductOffer], error) {
	s.log.Debug("Listing product offers", "page", opts.Page, "page_size", opts.PageSize)
	return s.db.ListProductOffers(ctx, opts)
}

// Assembly methods
func (s *Service) GetAssembly(ctx context.Context, id int) (*Assembly, error) {
	s.log.Debug("Getting assembly", "id", id)
	return s.db.GetAssembly(ctx, id)
}

func (s *Service) ListAssemblies(ctx context.Context, opts *ListOptions) (*ListResult[Assembly], error) {
	s.log.Debug("Listing assemblies", "page", opts.Page, "page_size", opts.PageSize)
	return s.db.ListAssemblies(ctx, opts)
}

func (s *Service) ListAssembliesByUser(ctx context.Context, userID int, opts *ListOptions) (*ListResult[Assembly], error) {
	s.log.Debug("Listing assemblies by user", "user_id", userID, "page", opts.Page)
	return s.db.ListAssembliesByUser(ctx, userID, opts)
}

// Compatibility methods
func (s *Service) CheckAssemblyCompatibility(ctx context.Context, assemblyID int) (*ValidationResult, error) {
	s.log.Info("Checking assembly compatibility", "assembly_id", assemblyID)
	return s.db.CheckAssemblyCompatibility(ctx, assemblyID)
}

func (s *Service) ValidateAssemblyComponents(ctx context.Context, components AssemblyComponents) (*ValidationResult, error) {
	s.log.Info("Validating assembly components")
	return s.db.ValidateAssemblyComponents(ctx, components)
}

// Search methods
func (s *Service) SearchComponents(ctx context.Context, query string, componentTypes []string, limit, offset int, filters map[string]string) ([]SearchResult, int, error) {
	s.log.Debug("Searching components", "query", query, "types", componentTypes)
	return s.db.SearchComponents(ctx, query, componentTypes, limit, offset, filters)
}

func (s *Service) GetComponentPrices(ctx context.Context, componentType string, componentID int) ([]PriceInfo, error) {
	s.log.Debug("Getting component prices", "type", componentType, "id", componentID)
	return s.db.GetComponentPrices(ctx, componentType, componentID)
}

func (s *Service) GetBestOffers(ctx context.Context, componentType string, componentID int, limit int, cheapestFirst bool) ([]Offer, error) {
	s.log.Debug("Getting best offers", "type", componentType, "id", componentID)
	return s.db.GetBestOffers(ctx, componentType, componentID, limit, cheapestFirst)
}

// GetComponentByID универсальный метод получения компонента
func (s *Service) GetComponentByID(ctx context.Context, componentType string, id int) (interface{}, error) {
	s.log.Debug("Getting component by ID", "type", componentType, "id", id)
	return s.db.GetComponentByID(ctx, componentType, id)
}
