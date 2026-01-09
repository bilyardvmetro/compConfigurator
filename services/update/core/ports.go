package core

import "context"

type DB interface {
	// GPU методы
	AddGpu(ctx context.Context, element GPU) error
	UpdateGpu(ctx context.Context, element GPU) error
	DeleteGpu(ctx context.Context, id int) error

	// CPU Cooler методы
	AddCpuCooler(ctx context.Context, element CpuCooler) error
	UpdateCpuCooler(ctx context.Context, element CpuCooler) error
	DeleteCpuCooler(ctx context.Context, id int) error

	// Case методы
	AddCase(ctx context.Context, element Case) error
	UpdateCase(ctx context.Context, element Case) error
	DeleteCase(ctx context.Context, id int) error

	// PSU методы
	AddPsu(ctx context.Context, element PSU) error
	UpdatePsu(ctx context.Context, element PSU) error
	DeletePsu(ctx context.Context, id int) error

	// CPU Socket методы
	AddCpuSocket(ctx context.Context, element CpuSocket) error
	UpdateCpuSocket(ctx context.Context, element CpuSocket) error
	DeleteCpuSocket(ctx context.Context, code string) error

	// Motherboard Form Factor методы
	AddMotherboardFormFactor(ctx context.Context, element MotherboardFormFactor) error
	UpdateMotherboardFormFactor(ctx context.Context, element MotherboardFormFactor) error
	DeleteMotherboardFormFactor(ctx context.Context, code string) error

	// CPU методы
	AddCpu(ctx context.Context, element CPU) error
	UpdateCpu(ctx context.Context, element CPU) error
	DeleteCpu(ctx context.Context, id int) error

	// Motherboard методы
	AddMotherboard(ctx context.Context, element Motherboard) error
	UpdateMotherboard(ctx context.Context, element Motherboard) error
	DeleteMotherboard(ctx context.Context, id int) error

	// RAM Kit методы
	AddRAMKit(ctx context.Context, element RAMKit) error
	UpdateRAMKit(ctx context.Context, element RAMKit) error
	DeleteRAMKit(ctx context.Context, id int) error

	// Storage Drive методы
	AddStorageDrive(ctx context.Context, element StorageDrive) error
	UpdateStorageDrive(ctx context.Context, element StorageDrive) error
	DeleteStorageDrive(ctx context.Context, id int) error

	// User методы
	AddUser(ctx context.Context, element User) error
	UpdateUser(ctx context.Context, element User) error
	DeleteUser(ctx context.Context, id int) error

	// Shop методы
	AddShop(ctx context.Context, element Shop) error
	UpdateShop(ctx context.Context, element Shop) error
	DeleteShop(ctx context.Context, id int) error

	// Product Offer методы
	AddProductOffer(ctx context.Context, element ProductOffer) error
	UpdateProductOffer(ctx context.Context, element ProductOffer) error
	DeleteProductOffer(ctx context.Context, id int) error

	// Assembly методы
	AddAssembly(ctx context.Context, element Assembly) error
	UpdateAssembly(ctx context.Context, element Assembly) error
	DeleteAssembly(ctx context.Context, id int) error

	// Cooler Socket методы (many-to-many)
	AddCoolerSocket(ctx context.Context, element CoolerSocket) error
	DeleteCoolerSocket(ctx context.Context, coolerID int, socketCode string) error
	DeleteCoolerSocketsByCooler(ctx context.Context, coolerID int) error

	// Case Form Factor Support методы (many-to-many)
	AddCaseFormFactorSupport(ctx context.Context, element CaseFormFactorSupport) error
	DeleteCaseFormFactorSupport(ctx context.Context, caseID int, formFactorCode string) error
	DeleteCaseFormFactorSupportsByCase(ctx context.Context, caseID int) error

	// Assembly RAM Kit методы (many-to-many)
	AddAssemblyRAMKit(ctx context.Context, element AssemblyRAMKit) error
	DeleteAssemblyRAMKit(ctx context.Context, assemblyID int, ramKitID int) error
	DeleteAssemblyRAMKitsByAssembly(ctx context.Context, assemblyID int) error

	// Assembly Drive методы (many-to-many)
	AddAssemblyDrive(ctx context.Context, element AssemblyDrive) error
	DeleteAssemblyDrive(ctx context.Context, assemblyID int, driveID int) error
	DeleteAssemblyDrivesByAssembly(ctx context.Context, assemblyID int) error
}

type Updater interface {
	DB
}
