package grpc

import (
	"context"

	updatepb "comp_config.com/services/proto/update"
	"comp_config.com/services/update/core"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewServer(service core.Updater) *Server {
	return &Server{service: service}
}

type Server struct {
	updatepb.UnimplementedUpdateServer
	service core.Updater
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// GPU methods
func (s *Server) AddGpu(ctx context.Context, in *updatepb.GpuRequest) (*emptypb.Empty, error) {
	if err := s.service.AddGpu(ctx, *toCoreGpuRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateGpu(ctx context.Context, in *updatepb.UpdateGpuRequest) (*emptypb.Empty, error) {
	gpu := toCoreGpuRequest(in.Gpu)
	gpu.ID = int(in.Id)

	if err := s.service.UpdateGpu(ctx, *gpu); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteGpu(ctx context.Context, in *updatepb.IdRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteGpu(ctx, int(in.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// CPU Cooler methods
func (s *Server) AddCpuCooler(ctx context.Context, in *updatepb.CpuCoolerRequest) (*emptypb.Empty, error) {
	if err := s.service.AddCpuCooler(ctx, *toCoreCpuCoolerRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateCpuCooler(ctx context.Context, in *updatepb.UpdateCpuCoolerRequest) (*emptypb.Empty, error) {
	cooler := toCoreCpuCoolerRequest(in.CpuCooler)
	cooler.ID = int(in.Id)

	if err := s.service.UpdateCpuCooler(ctx, *cooler); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteCpuCooler(ctx context.Context, in *updatepb.IdRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteCpuCooler(ctx, int(in.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Case methods
func (s *Server) AddCase(ctx context.Context, in *updatepb.CaseRequest) (*emptypb.Empty, error) {
	if err := s.service.AddCase(ctx, *toCoreCaseRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateCase(ctx context.Context, in *updatepb.UpdateCaseRequest) (*emptypb.Empty, error) {
	case_ := toCoreCaseRequest(in.Case)
	case_.ID = int(in.Id)

	if err := s.service.UpdateCase(ctx, *case_); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteCase(ctx context.Context, in *updatepb.IdRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteCase(ctx, int(in.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// PSU methods
func (s *Server) AddPsu(ctx context.Context, in *updatepb.PsuRequest) (*emptypb.Empty, error) {
	if err := s.service.AddPsu(ctx, *toCorePsuRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdatePsu(ctx context.Context, in *updatepb.UpdatePsuRequest) (*emptypb.Empty, error) {
	psu := toCorePsuRequest(in.Psu)
	psu.ID = int(in.Id)

	if err := s.service.UpdatePsu(ctx, *psu); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeletePsu(ctx context.Context, in *updatepb.IdRequest) (*emptypb.Empty, error) {
	if err := s.service.DeletePsu(ctx, int(in.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// CPU Socket methods
func (s *Server) AddCpuSocket(ctx context.Context, in *updatepb.CpuSocketRequest) (*emptypb.Empty, error) {
	if err := s.service.AddCpuSocket(ctx, *toCoreCpuSocketRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateCpuSocket(ctx context.Context, in *updatepb.UpdateCpuSocketRequest) (*emptypb.Empty, error) {
	socket := toCoreCpuSocketRequest(in.CpuSocket)
	socket.Code = in.SocketCode

	if err := s.service.UpdateCpuSocket(ctx, *socket); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteCpuSocket(ctx context.Context, in *updatepb.StringRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteCpuSocket(ctx, in.Value); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Motherboard Form Factor methods
func (s *Server) AddMotherboardFormFactor(ctx context.Context, in *updatepb.MotherboardFormFactorRequest) (*emptypb.Empty, error) {
	if err := s.service.AddMotherboardFormFactor(ctx, *toCoreMotherboardFormFactorRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateMotherboardFormFactor(ctx context.Context, in *updatepb.UpdateMotherboardFormFactorRequest) (*emptypb.Empty, error) {
	formFactor := toCoreMotherboardFormFactorRequest(in.MotherboardFormFactor)
	formFactor.Code = in.FormFactorCode

	if err := s.service.UpdateMotherboardFormFactor(ctx, *formFactor); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteMotherboardFormFactor(ctx context.Context, in *updatepb.StringRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteMotherboardFormFactor(ctx, in.Value); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// CPU methods
func (s *Server) AddCpu(ctx context.Context, in *updatepb.CpuRequest) (*emptypb.Empty, error) {
	if err := s.service.AddCpu(ctx, *toCoreCpuRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateCpu(ctx context.Context, in *updatepb.UpdateCpuRequest) (*emptypb.Empty, error) {
	cpu := toCoreCpuRequest(in.Cpu)
	cpu.ID = int(in.Id)

	if err := s.service.UpdateCpu(ctx, *cpu); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteCpu(ctx context.Context, in *updatepb.IdRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteCpu(ctx, int(in.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Motherboard methods
func (s *Server) AddMotherboard(ctx context.Context, in *updatepb.MotherboardRequest) (*emptypb.Empty, error) {
	if err := s.service.AddMotherboard(ctx, *toCoreMotherboardRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateMotherboard(ctx context.Context, in *updatepb.UpdateMotherboardRequest) (*emptypb.Empty, error) {
	mb := toCoreMotherboardRequest(in.Motherboard)
	mb.ID = int(in.Id)

	if err := s.service.UpdateMotherboard(ctx, *mb); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteMotherboard(ctx context.Context, in *updatepb.IdRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteMotherboard(ctx, int(in.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Конвертеры
func toCoreGpuRequest(in *updatepb.GpuRequest) *core.GPU {
	return &core.GPU{
		Name:              in.Name,
		Brand:             in.Brand,
		Interface:         in.Interface,
		PcieVersion:       int(in.PcieVersion),
		PcieLanesRequired: int(in.PcieLanesRequired),
		TdpWatt:           int(in.TdpWatt),
		PowerConnectors:   in.PowerConnectors,
		LengthMm:          int(in.LengthMm),
		WidthSlots:        int(in.WidthSlots),
		HeightMm:          int(in.HeightMm),
	}
}

func toCoreCpuCoolerRequest(in *updatepb.CpuCoolerRequest) *core.CpuCooler {
	return &core.CpuCooler{
		Name:           in.Name,
		Brand:          in.Brand,
		CoolingType:    in.CoolingType,
		TdpLimitWatt:   int(in.TdpLimitWatt),
		HeightMm:       int(in.HeightMm),
		FanCount:       int(in.FanCount),
		FanControlType: in.FanControlType,
	}
}

func toCoreCaseRequest(in *updatepb.CaseRequest) *core.Case {
	return &core.Case{
		Name:              in.Name,
		Brand:             in.Brand,
		PsuFormFactor:     in.PsuFormFactor,
		MaxGpuLengthMm:    optionalInt(in.MaxGpuLengthMm),
		MaxGpuWidthSlots:  optionalFloat(in.MaxGpuWidthSlots),
		MaxCoolerHeightMm: optionalInt(in.MaxCoolerHeightMm),
		MaxPsuLengthMm:    optionalInt(in.MaxPsuLengthMm),
		DriveBays3_5Count: optionalInt(in.DriveBays_3_5Count),
		DriveBays2_5Count: optionalInt(in.DriveBays_2_5Count),
	}
}

func toCorePsuRequest(in *updatepb.PsuRequest) *core.PSU {
	return &core.PSU{
		Name:                      in.Name,
		Brand:                     in.Brand,
		PowerWatt:                 int(in.PowerWatt),
		FormFactor:                in.FormFactor,
		EfficiencyRating:          in.EfficiencyRating,
		PcieConnectors6_8pinCount: int(in.PcieConnectors_6_8PinCount),
		Cpu8pinConnectorsCount:    int(in.Cpu_8PinConnectorsCount),
		SataConnectorsCount:       int(in.SataConnectorsCount),
		MolexConnectorsCount:      int(in.MolexConnectorsCount),
		LengthMm:                  optionalInt(in.LengthMm),
	}
}

func toCoreCpuSocketRequest(in *updatepb.CpuSocketRequest) *core.CpuSocket {
	return &core.CpuSocket{
		Code:        in.SocketCode,
		Description: in.Description,
	}
}

func toCoreMotherboardFormFactorRequest(in *updatepb.MotherboardFormFactorRequest) *core.MotherboardFormFactor {
	return &core.MotherboardFormFactor{
		Code:        in.FormFactorCode,
		Description: in.Description,
		WidthMm:     optionalInt(in.WidthMm),
		HeightMm:    optionalInt(in.HeightMm),
	}
}

func toCoreCpuRequest(in *updatepb.CpuRequest) *core.CPU {
	return &core.CPU{
		Name:                   in.Name,
		Brand:                  in.Brand,
		SocketCode:             in.SocketCode,
		Architecture:           in.Architecture,
		CoreCount:              int(in.CoreCount),
		ThreadCount:            int(in.ThreadCount),
		BaseClockMhz:           int(in.BaseClockMhz),
		BoostClockMhz:          optionalInt(in.BoostClockMhz),
		TdpWatt:                int(in.TdpWatt),
		HasIntegratedGpu:       in.HasIntegratedGpu,
		SupportedRamType:       in.SupportedRamType,
		SupportedRamFreqMaxMhz: optionalInt(in.SupportedRamFreqMaxMhz),
		MemoryChannels:         int(in.MemoryChannels),
		PcieVersion:            int(in.PcieVersion),
		PcieLanesTotal:         int(in.PcieLanesTotal),
	}
}

func toCoreMotherboardRequest(in *updatepb.MotherboardRequest) *core.Motherboard {
	return &core.Motherboard{
		Name:                  in.Name,
		Brand:                 in.Brand,
		SocketCode:            in.SocketCode,
		FormFactorCode:        in.FormFactorCode,
		Chipset:               in.Chipset,
		RamType:               in.RamType,
		RamSlots:              int(in.RamSlots),
		RamCapacityMaxGb:      optionalInt(in.RamCapacityMaxGb),
		RamFreqMaxMhz:         optionalInt(in.RamFreqMaxMhz),
		PcieX16SlotsCount:     int(in.PcieX16SlotsCount),
		PcieVersionMax:        int(in.PcieVersionMax),
		M2SlotsCount:          int(in.M2SlotsCount),
		SataPortsCount:        int(in.SataPortsCount),
		PsuMainConnectorType:  in.PsuMainConnectorType,
		CpuPowerConnectorType: in.CpuPowerConnectorType,
	}
}

func optionalInt(val *int64) *int {
	if val == nil {
		return nil
	}
	v := int(*val)
	return &v
}

func optionalFloat(val *float64) *float64 {
	if val == nil {
		return nil
	}
	return val
}

// RAM Kit methods
func (s *Server) AddRAMKit(ctx context.Context, in *updatepb.RAMKitRequest) (*emptypb.Empty, error) {
	if err := s.service.AddRAMKit(ctx, *toCoreRAMKitRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateRAMKit(ctx context.Context, in *updatepb.UpdateRAMKitRequest) (*emptypb.Empty, error) {
	ramKit := toCoreRAMKitRequest(in.RamKit)
	ramKit.ID = int(in.Id)

	if err := s.service.UpdateRAMKit(ctx, *ramKit); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteRAMKit(ctx context.Context, in *updatepb.IdRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteRAMKit(ctx, int(in.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Storage Drive methods
func (s *Server) AddStorageDrive(ctx context.Context, in *updatepb.StorageDriveRequest) (*emptypb.Empty, error) {
	if err := s.service.AddStorageDrive(ctx, *toCoreStorageDriveRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateStorageDrive(ctx context.Context, in *updatepb.UpdateStorageDriveRequest) (*emptypb.Empty, error) {
	drive := toCoreStorageDriveRequest(in.StorageDrive)
	drive.ID = int(in.Id)

	if err := s.service.UpdateStorageDrive(ctx, *drive); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteStorageDrive(ctx context.Context, in *updatepb.IdRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteStorageDrive(ctx, int(in.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// User methods
func (s *Server) AddUser(ctx context.Context, in *updatepb.UserRequest) (*emptypb.Empty, error) {
	if err := s.service.AddUser(ctx, *toCoreUserRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateUser(ctx context.Context, in *updatepb.UpdateUserRequest) (*emptypb.Empty, error) {
	user := toCoreUserRequest(in.User)
	user.ID = int(in.Id)

	if err := s.service.UpdateUser(ctx, *user); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteUser(ctx context.Context, in *updatepb.IdRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteUser(ctx, int(in.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Shop methods
func (s *Server) AddShop(ctx context.Context, in *updatepb.ShopRequest) (*emptypb.Empty, error) {
	if err := s.service.AddShop(ctx, *toCoreShopRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateShop(ctx context.Context, in *updatepb.UpdateShopRequest) (*emptypb.Empty, error) {
	shop := toCoreShopRequest(in.Shop)
	shop.ID = int(in.Id)

	if err := s.service.UpdateShop(ctx, *shop); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteShop(ctx context.Context, in *updatepb.IdRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteShop(ctx, int(in.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Product Offer methods
func (s *Server) AddProductOffer(ctx context.Context, in *updatepb.ProductOfferRequest) (*emptypb.Empty, error) {
	if err := s.service.AddProductOffer(ctx, *toCoreProductOfferRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateProductOffer(ctx context.Context, in *updatepb.UpdateProductOfferRequest) (*emptypb.Empty, error) {
	offer := toCoreProductOfferRequest(in.ProductOffer)
	offer.ID = int(in.Id)

	if err := s.service.UpdateProductOffer(ctx, *offer); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteProductOffer(ctx context.Context, in *updatepb.IdRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteProductOffer(ctx, int(in.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Assembly methods
func (s *Server) AddAssembly(ctx context.Context, in *updatepb.AssemblyRequest) (*emptypb.Empty, error) {
	if err := s.service.AddAssembly(ctx, *toCoreAssemblyRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) UpdateAssembly(ctx context.Context, in *updatepb.UpdateAssemblyRequest) (*emptypb.Empty, error) {
	assembly := toCoreAssemblyRequest(in.Assembly)
	assembly.ID = int(in.Id)

	if err := s.service.UpdateAssembly(ctx, *assembly); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteAssembly(ctx context.Context, in *updatepb.IdRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteAssembly(ctx, int(in.Id)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Cooler Socket methods (many-to-many)
func (s *Server) AddCoolerSocket(ctx context.Context, in *updatepb.CoolerSocketRequest) (*emptypb.Empty, error) {
	if err := s.service.AddCoolerSocket(ctx, *toCoreCoolerSocketRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteCoolerSocket(ctx context.Context, in *updatepb.DeleteCoolerSocketRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteCoolerSocket(ctx, int(in.CoolerId), in.SocketCode); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteCoolerSocketsByCooler(ctx context.Context, in *updatepb.DeleteByCoolerRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteCoolerSocketsByCooler(ctx, int(in.CoolerId)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Case Form Factor Support methods (many-to-many)
func (s *Server) AddCaseFormFactorSupport(ctx context.Context, in *updatepb.CaseFormFactorSupportRequest) (*emptypb.Empty, error) {
	if err := s.service.AddCaseFormFactorSupport(ctx, *toCoreCaseFormFactorSupportRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteCaseFormFactorSupport(ctx context.Context, in *updatepb.DeleteCaseFormFactorSupportRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteCaseFormFactorSupport(ctx, int(in.CaseId), in.FormFactorCode); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteCaseFormFactorSupportsByCase(ctx context.Context, in *updatepb.DeleteByCaseRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteCaseFormFactorSupportsByCase(ctx, int(in.CaseId)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Assembly RAM Kit methods (many-to-many)
func (s *Server) AddAssemblyRAMKit(ctx context.Context, in *updatepb.AssemblyRAMKitRequest) (*emptypb.Empty, error) {
	if err := s.service.AddAssemblyRAMKit(ctx, *toCoreAssemblyRAMKitRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteAssemblyRAMKit(ctx context.Context, in *updatepb.DeleteAssemblyRAMKitRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteAssemblyRAMKit(ctx, int(in.AssemblyId), int(in.RamKitId)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteAssemblyRAMKitsByAssembly(ctx context.Context, in *updatepb.DeleteByAssemblyRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteAssemblyRAMKitsByAssembly(ctx, int(in.AssemblyId)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

// Assembly Drive methods (many-to-many)
func (s *Server) AddAssemblyDrive(ctx context.Context, in *updatepb.AssemblyDriveRequest) (*emptypb.Empty, error) {
	if err := s.service.AddAssemblyDrive(ctx, *toCoreAssemblyDriveRequest(in)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteAssemblyDrive(ctx context.Context, in *updatepb.DeleteAssemblyDriveRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteAssemblyDrive(ctx, int(in.AssemblyId), int(in.DriveId)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) DeleteAssemblyDrivesByAssembly(ctx context.Context, in *updatepb.DeleteByAssemblyRequest) (*emptypb.Empty, error) {
	if err := s.service.DeleteAssemblyDrivesByAssembly(ctx, int(in.AssemblyId)); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func toCoreRAMKitRequest(in *updatepb.RAMKitRequest) *core.RAMKit {
	return &core.RAMKit{
		Name:             in.Name,
		Brand:            in.Brand,
		RamType:          in.RamType,
		ModuleCapacityGb: int(in.ModuleCapacityGb),
		ModuleCount:      int(in.ModuleCount),
		TotalCapacityGb:  int(in.TotalCapacityGb),
		FreqMhz:          int(in.FreqMhz),
		Timings:          in.Timings,
		VoltageV:         optionalDouble(in.VoltageV),
		FormFactor:       in.FormFactor,
	}
}

func toCoreStorageDriveRequest(in *updatepb.StorageDriveRequest) *core.StorageDrive {
	return &core.StorageDrive{
		Name:       in.Name,
		Brand:      in.Brand,
		DriveType:  in.DriveType,
		FormFactor: in.FormFactor,
		Interface:  in.Interface,
		CapacityGb: int(in.CapacityGb),
	}
}

func toCoreUserRequest(in *updatepb.UserRequest) *core.User {
	return &core.User{
		Email:        in.Email,
		PasswordHash: in.PasswordHash,
		Nickname:     optionalString(in.Nickname),
		AvatarUrl:    optionalString(in.AvatarUrl),
		IsAdmin:      in.IsAdmin,
	}
}

func toCoreShopRequest(in *updatepb.ShopRequest) *core.Shop {
	return &core.Shop{
		Name: in.Name,
		Url:  optionalString(in.Url),
	}
}

func toCoreProductOfferRequest(in *updatepb.ProductOfferRequest) *core.ProductOffer {
	return &core.ProductOffer{
		ShopID:        int(in.ShopId),
		ComponentType: in.ComponentType,
		ComponentID:   int(in.ComponentId),
		Price:         in.Price,
		Available:     in.Available,
	}
}

func toCoreAssemblyRequest(in *updatepb.AssemblyRequest) *core.Assembly {
	return &core.Assembly{
		UserID:           int(in.UserId),
		Name:             in.Name,
		IsPublic:         in.IsPublic,
		TotalPriceCached: optionalDouble(in.TotalPriceCached),
		CpuID:            int(in.CpuId),
		GpuID:            optionalInt(in.GpuId),
		MotherboardID:    int(in.MotherboardId),
		PsuID:            int(in.PsuId),
		CaseID:           int(in.CaseId),
		CoolerID:         optionalInt(in.CoolerId),
	}
}

func toCoreCoolerSocketRequest(in *updatepb.CoolerSocketRequest) *core.CoolerSocket {
	return &core.CoolerSocket{
		CoolerID:   int(in.CoolerId),
		SocketCode: in.SocketCode,
		Notes:      optionalString(in.Notes),
	}
}

func toCoreCaseFormFactorSupportRequest(in *updatepb.CaseFormFactorSupportRequest) *core.CaseFormFactorSupport {
	return &core.CaseFormFactorSupport{
		CaseID:         int(in.CaseId),
		FormFactorCode: in.FormFactorCode,
	}
}

func toCoreAssemblyRAMKitRequest(in *updatepb.AssemblyRAMKitRequest) *core.AssemblyRAMKit {
	return &core.AssemblyRAMKit{
		AssemblyID: int(in.AssemblyId),
		RAMKitID:   int(in.RamKitId),
	}
}

func toCoreAssemblyDriveRequest(in *updatepb.AssemblyDriveRequest) *core.AssemblyDrive {
	return &core.AssemblyDrive{
		AssemblyID: int(in.AssemblyId),
		DriveID:    int(in.DriveId),
		MountType:  optionalString(in.MountType),
	}
}

func optionalString(val *string) *string {
	if val == nil || *val == "" {
		return nil
	}
	return val
}

func optionalDouble(val *float64) *float64 {
	if val == nil {
		return nil
	}
	return val
}
