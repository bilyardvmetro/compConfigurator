package core

import "context"

type DB interface {
	// GPU methods
	GetGPU(ctx context.Context, id int) (*GPU, error)
	ListGPUs(ctx context.Context, opts *ListOptions) (*ListResult[GPU], error)

	// CPU Cooler methods
	GetCpuCooler(ctx context.Context, id int) (*CpuCooler, error)
	ListCpuCoolers(ctx context.Context, opts *ListOptions) (*ListResult[CpuCooler], error)

	// Case methods
	GetCase(ctx context.Context, id int) (*Case, error)
	ListCases(ctx context.Context, opts *ListOptions) (*ListResult[Case], error)

	// PSU methods
	GetPSU(ctx context.Context, id int) (*PSU, error)
	ListPSUs(ctx context.Context, opts *ListOptions) (*ListResult[PSU], error)

	// CPU Socket methods
	GetCpuSocket(ctx context.Context, code string) (*CpuSocket, error)
	ListCpuSockets(ctx context.Context, opts *ListOptions) (*ListResult[CpuSocket], error)

	// Motherboard Form Factor methods
	GetMotherboardFormFactor(ctx context.Context, code string) (*MotherboardFormFactor, error)
	ListMotherboardFormFactors(ctx context.Context, opts *ListOptions) (*ListResult[MotherboardFormFactor], error)

	// CPU methods
	GetCPU(ctx context.Context, id int) (*CPU, error)
	ListCPUs(ctx context.Context, opts *ListOptions) (*ListResult[CPU], error)

	// Motherboard methods
	GetMotherboard(ctx context.Context, id int) (*Motherboard, error)
	ListMotherboards(ctx context.Context, opts *ListOptions) (*ListResult[Motherboard], error)

	// RAM Kit methods
	GetRAMKit(ctx context.Context, id int) (*RAMKit, error)
	ListRAMKits(ctx context.Context, opts *ListOptions) (*ListResult[RAMKit], error)

	// Storage Drive methods
	GetStorageDrive(ctx context.Context, id int) (*StorageDrive, error)
	ListStorageDrives(ctx context.Context, opts *ListOptions) (*ListResult[StorageDrive], error)

	// User methods
	GetUser(ctx context.Context, id int) (*User, error)
	ListUsers(ctx context.Context, opts *ListOptions) (*ListResult[User], error)

	// Shop methods
	GetShop(ctx context.Context, id int) (*Shop, error)
	ListShops(ctx context.Context, opts *ListOptions) (*ListResult[Shop], error)

	// Product Offer methods
	GetProductOffer(ctx context.Context, id int) (*ProductOffer, error)
	ListProductOffers(ctx context.Context, opts *ListOptions) (*ListResult[ProductOffer], error)

	// Assembly methods
	GetAssembly(ctx context.Context, id int) (*Assembly, error)
	ListAssemblies(ctx context.Context, opts *ListOptions) (*ListResult[Assembly], error)
	ListAssembliesByUser(ctx context.Context, userID int, opts *ListOptions) (*ListResult[Assembly], error)

	// Compatibility methods
	CheckAssemblyCompatibility(ctx context.Context, assemblyID int) (*ValidationResult, error)
	ValidateAssemblyComponents(ctx context.Context, components AssemblyComponents) (*ValidationResult, error)

	// Search methods
	SearchComponents(ctx context.Context, query string, componentTypes []string, limit, offset int, filters map[string]string) ([]SearchResult, int, error)
	GetComponentPrices(ctx context.Context, componentType string, componentID int) ([]PriceInfo, error)
	GetBestOffers(ctx context.Context, componentType string, componentID int, limit int, cheapestFirst bool) ([]Offer, error)
	GetComponentByID(ctx context.Context, componentType string, id int) (interface{}, error)
}

type Searcher interface {
	DB
}
