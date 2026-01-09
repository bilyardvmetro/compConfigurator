package core

import (
	"context"
	"log/slog"
)

type Service struct {
	log *slog.Logger
	db  DB
}

func NewSerivce(log *slog.Logger, db DB) (*Service, error) {
	return &Service{
		log: log,
		db:  db,
	}, nil
}

func chooseErr(err error, methodType string) error {
	if err == ErrNotFound || err == nil {
		return err
	}
	switch methodType {
	case "add":
		return ErrFailedAdding
	case "update":
		return ErrFailedUpdate
	case "delete":
		return ErrFailedDelete
	default:
		return nil
	}
}

// GPU methods
func (s *Service) AddGpu(ctx context.Context, element GPU) error {
	return chooseErr(s.db.AddGpu(ctx, element), "add")
}

func (s *Service) UpdateGpu(ctx context.Context, element GPU) error {
	return chooseErr(s.db.UpdateGpu(ctx, element), "update")
}

func (s *Service) DeleteGpu(ctx context.Context, id int) error {
	return chooseErr(s.db.DeleteGpu(ctx, id), "delete")
}

// CPU Cooler methods
func (s *Service) AddCpuCooler(ctx context.Context, element CpuCooler) error {
	return chooseErr(s.db.AddCpuCooler(ctx, element), "add")
}

func (s *Service) UpdateCpuCooler(ctx context.Context, element CpuCooler) error {
	return chooseErr(s.db.UpdateCpuCooler(ctx, element), "update")
}

func (s *Service) DeleteCpuCooler(ctx context.Context, id int) error {
	return chooseErr(s.db.DeleteCpuCooler(ctx, id), "delete")
}

// Case methods
func (s *Service) AddCase(ctx context.Context, element Case) error {
	return chooseErr(s.db.AddCase(ctx, element), "add")
}

func (s *Service) UpdateCase(ctx context.Context, element Case) error {
	return chooseErr(s.db.UpdateCase(ctx, element), "update")
}

func (s *Service) DeleteCase(ctx context.Context, id int) error {
	return chooseErr(s.db.DeleteCase(ctx, id), "delete")
}

// PSU methods
func (s *Service) AddPsu(ctx context.Context, element PSU) error {
	return chooseErr(s.db.AddPsu(ctx, element), "add")
}

func (s *Service) UpdatePsu(ctx context.Context, element PSU) error {
	return chooseErr(s.db.UpdatePsu(ctx, element), "update")
}

func (s *Service) DeletePsu(ctx context.Context, id int) error {
	return chooseErr(s.db.DeletePsu(ctx, id), "delete")
}

// CPU Socket methods
func (s *Service) AddCpuSocket(ctx context.Context, element CpuSocket) error {
	return chooseErr(s.db.AddCpuSocket(ctx, element), "add")
}

func (s *Service) UpdateCpuSocket(ctx context.Context, element CpuSocket) error {
	return chooseErr(s.db.UpdateCpuSocket(ctx, element), "update")
}

func (s *Service) DeleteCpuSocket(ctx context.Context, code string) error {
	return chooseErr(s.db.DeleteCpuSocket(ctx, code), "delete")
}

// Motherboard Form Factor methods
func (s *Service) AddMotherboardFormFactor(ctx context.Context, element MotherboardFormFactor) error {
	return chooseErr(s.db.AddMotherboardFormFactor(ctx, element), "add")
}

func (s *Service) UpdateMotherboardFormFactor(ctx context.Context, element MotherboardFormFactor) error {
	return chooseErr(s.db.UpdateMotherboardFormFactor(ctx, element), "update")
}

func (s *Service) DeleteMotherboardFormFactor(ctx context.Context, code string) error {
	return chooseErr(s.db.DeleteMotherboardFormFactor(ctx, code), "delete")
}

// CPU methods
func (s *Service) AddCpu(ctx context.Context, element CPU) error {
	return chooseErr(s.db.AddCpu(ctx, element), "add")
}

func (s *Service) UpdateCpu(ctx context.Context, element CPU) error {
	return chooseErr(s.db.UpdateCpu(ctx, element), "update")
}

func (s *Service) DeleteCpu(ctx context.Context, id int) error {
	return chooseErr(s.db.DeleteCpu(ctx, id), "delete")
}

// Motherboard methods
func (s *Service) AddMotherboard(ctx context.Context, element Motherboard) error {
	return chooseErr(s.db.AddMotherboard(ctx, element), "add")
}

func (s *Service) UpdateMotherboard(ctx context.Context, element Motherboard) error {
	return chooseErr(s.db.UpdateMotherboard(ctx, element), "update")
}

func (s *Service) DeleteMotherboard(ctx context.Context, id int) error {
	return chooseErr(s.db.DeleteMotherboard(ctx, id), "delete")
}

// RAM Kit methods
func (s *Service) AddRAMKit(ctx context.Context, element RAMKit) error {
	return chooseErr(s.db.AddRAMKit(ctx, element), "add")
}

func (s *Service) UpdateRAMKit(ctx context.Context, element RAMKit) error {
	return chooseErr(s.db.UpdateRAMKit(ctx, element), "update")
}

func (s *Service) DeleteRAMKit(ctx context.Context, id int) error {
	return chooseErr(s.db.DeleteRAMKit(ctx, id), "delete")
}

// Storage Drive methods
func (s *Service) AddStorageDrive(ctx context.Context, element StorageDrive) error {
	return chooseErr(s.db.AddStorageDrive(ctx, element), "add")
}

func (s *Service) UpdateStorageDrive(ctx context.Context, element StorageDrive) error {
	return chooseErr(s.db.UpdateStorageDrive(ctx, element), "update")
}

func (s *Service) DeleteStorageDrive(ctx context.Context, id int) error {
	return chooseErr(s.db.DeleteStorageDrive(ctx, id), "delete")
}

// User methods
func (s *Service) AddUser(ctx context.Context, element User) error {
	return chooseErr(s.db.AddUser(ctx, element), "add")
}

func (s *Service) UpdateUser(ctx context.Context, element User) error {
	return chooseErr(s.db.UpdateUser(ctx, element), "update")
}

func (s *Service) DeleteUser(ctx context.Context, id int) error {
	return chooseErr(s.db.DeleteUser(ctx, id), "delete")
}

// Shop methods
func (s *Service) AddShop(ctx context.Context, element Shop) error {
	return chooseErr(s.db.AddShop(ctx, element), "add")
}

func (s *Service) UpdateShop(ctx context.Context, element Shop) error {
	return chooseErr(s.db.UpdateShop(ctx, element), "update")
}

func (s *Service) DeleteShop(ctx context.Context, id int) error {
	return chooseErr(s.db.DeleteShop(ctx, id), "delete")
}

// Product Offer methods
func (s *Service) AddProductOffer(ctx context.Context, element ProductOffer) error {
	return chooseErr(s.db.AddProductOffer(ctx, element), "add")
}

func (s *Service) UpdateProductOffer(ctx context.Context, element ProductOffer) error {
	return chooseErr(s.db.UpdateProductOffer(ctx, element), "update")
}

func (s *Service) DeleteProductOffer(ctx context.Context, id int) error {
	return chooseErr(s.db.DeleteProductOffer(ctx, id), "delete")
}

// Assembly methods
func (s *Service) AddAssembly(ctx context.Context, element Assembly) error {
	return chooseErr(s.db.AddAssembly(ctx, element), "add")
}

func (s *Service) UpdateAssembly(ctx context.Context, element Assembly) error {
	return chooseErr(s.db.UpdateAssembly(ctx, element), "update")
}

func (s *Service) DeleteAssembly(ctx context.Context, id int) error {
	return chooseErr(s.db.DeleteAssembly(ctx, id), "delete")
}

// Cooler Socket methods (many-to-many)
func (s *Service) AddCoolerSocket(ctx context.Context, element CoolerSocket) error {
	return chooseErr(s.db.AddCoolerSocket(ctx, element), "add")
}

func (s *Service) DeleteCoolerSocket(ctx context.Context, coolerID int, socketCode string) error {
	return chooseErr(s.db.DeleteCoolerSocket(ctx, coolerID, socketCode), "delete")
}

func (s *Service) DeleteCoolerSocketsByCooler(ctx context.Context, coolerID int) error {
	return chooseErr(s.db.DeleteCoolerSocketsByCooler(ctx, coolerID), "delete")
}

// Case Form Factor Support methods (many-to-many)
func (s *Service) AddCaseFormFactorSupport(ctx context.Context, element CaseFormFactorSupport) error {
	return chooseErr(s.db.AddCaseFormFactorSupport(ctx, element), "add")
}

func (s *Service) DeleteCaseFormFactorSupport(ctx context.Context, caseID int, formFactorCode string) error {
	return chooseErr(s.db.DeleteCaseFormFactorSupport(ctx, caseID, formFactorCode), "delete")
}

func (s *Service) DeleteCaseFormFactorSupportsByCase(ctx context.Context, caseID int) error {
	return chooseErr(s.db.DeleteCaseFormFactorSupportsByCase(ctx, caseID), "delete")
}

// Assembly RAM Kit methods (many-to-many)
func (s *Service) AddAssemblyRAMKit(ctx context.Context, element AssemblyRAMKit) error {
	return chooseErr(s.db.AddAssemblyRAMKit(ctx, element), "add")
}

func (s *Service) DeleteAssemblyRAMKit(ctx context.Context, assemblyID int, ramKitID int) error {
	return chooseErr(s.db.DeleteAssemblyRAMKit(ctx, assemblyID, ramKitID), "delete")
}

func (s *Service) DeleteAssemblyRAMKitsByAssembly(ctx context.Context, assemblyID int) error {
	return chooseErr(s.db.DeleteAssemblyRAMKitsByAssembly(ctx, assemblyID), "delete")
}

// Assembly Drive methods (many-to-many)
func (s *Service) AddAssemblyDrive(ctx context.Context, element AssemblyDrive) error {
	return chooseErr(s.db.AddAssemblyDrive(ctx, element), "add")
}

func (s *Service) DeleteAssemblyDrive(ctx context.Context, assemblyID int, driveID int) error {
	return chooseErr(s.db.DeleteAssemblyDrive(ctx, assemblyID, driveID), "delete")
}

func (s *Service) DeleteAssemblyDrivesByAssembly(ctx context.Context, assemblyID int) error {
	return chooseErr(s.db.DeleteAssemblyDrivesByAssembly(ctx, assemblyID), "delete")
}
