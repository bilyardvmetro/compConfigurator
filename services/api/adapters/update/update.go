package update

import (
	"context"
	"log/slog"

	updatepb "comp_config.com/services/proto/update"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log    *slog.Logger
	client updatepb.UpdateClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		client: updatepb.NewUpdateClient(conn),
		log:    log,
	}, nil
}

func (c Client) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, nil)
	return err
}

// GPU methods
func (c Client) AddGPU(ctx context.Context, req *updatepb.GpuRequest) error {
	_, err := c.client.AddGpu(ctx, req)
	return err
}

func (c Client) UpdateGPU(ctx context.Context, req *updatepb.UpdateGpuRequest) error {
	_, err := c.client.UpdateGpu(ctx, req)
	return err
}

func (c Client) DeleteGPU(ctx context.Context, id int64) error {
	_, err := c.client.DeleteGpu(ctx, &updatepb.IdRequest{Id: id})
	return err
}

// CPU Cooler methods
func (c Client) AddCPUCooler(ctx context.Context, req *updatepb.CpuCoolerRequest) error {
	_, err := c.client.AddCpuCooler(ctx, req)
	return err
}

func (c Client) UpdateCPUCooler(ctx context.Context, req *updatepb.UpdateCpuCoolerRequest) error {
	_, err := c.client.UpdateCpuCooler(ctx, req)
	return err
}

func (c Client) DeleteCPUCooler(ctx context.Context, id int64) error {
	_, err := c.client.DeleteCpuCooler(ctx, &updatepb.IdRequest{Id: id})
	return err
}

// Case methods
func (c Client) AddCase(ctx context.Context, req *updatepb.CaseRequest) error {
	_, err := c.client.AddCase(ctx, req)
	return err
}

func (c Client) UpdateCase(ctx context.Context, req *updatepb.UpdateCaseRequest) error {
	_, err := c.client.UpdateCase(ctx, req)
	return err
}

func (c Client) DeleteCase(ctx context.Context, id int64) error {
	_, err := c.client.DeleteCase(ctx, &updatepb.IdRequest{Id: id})
	return err
}

// PSU methods
func (c Client) AddPSU(ctx context.Context, req *updatepb.PsuRequest) error {
	_, err := c.client.AddPsu(ctx, req)
	return err
}

func (c Client) UpdatePSU(ctx context.Context, req *updatepb.UpdatePsuRequest) error {
	_, err := c.client.UpdatePsu(ctx, req)
	return err
}

func (c Client) DeletePSU(ctx context.Context, id int64) error {
	_, err := c.client.DeletePsu(ctx, &updatepb.IdRequest{Id: id})
	return err
}

// CPU Socket methods
func (c Client) AddCPUSocket(ctx context.Context, req *updatepb.CpuSocketRequest) error {
	_, err := c.client.AddCpuSocket(ctx, req)
	return err
}

func (c Client) UpdateCPUSocket(ctx context.Context, req *updatepb.UpdateCpuSocketRequest) error {
	_, err := c.client.UpdateCpuSocket(ctx, req)
	return err
}

func (c Client) DeleteCPUSocket(ctx context.Context, code string) error {
	_, err := c.client.DeleteCpuSocket(ctx, &updatepb.StringRequest{Value: code})
	return err
}

// Motherboard Form Factor methods
func (c Client) AddMotherboardFormFactor(ctx context.Context, req *updatepb.MotherboardFormFactorRequest) error {
	_, err := c.client.AddMotherboardFormFactor(ctx, req)
	return err
}

func (c Client) UpdateMotherboardFormFactor(ctx context.Context, req *updatepb.UpdateMotherboardFormFactorRequest) error {
	_, err := c.client.UpdateMotherboardFormFactor(ctx, req)
	return err
}

func (c Client) DeleteMotherboardFormFactor(ctx context.Context, code string) error {
	_, err := c.client.DeleteMotherboardFormFactor(ctx, &updatepb.StringRequest{Value: code})
	return err
}

// CPU methods
func (c Client) AddCPU(ctx context.Context, req *updatepb.CpuRequest) error {
	_, err := c.client.AddCpu(ctx, req)
	return err
}

func (c Client) UpdateCPU(ctx context.Context, req *updatepb.UpdateCpuRequest) error {
	_, err := c.client.UpdateCpu(ctx, req)
	return err
}

func (c Client) DeleteCPU(ctx context.Context, id int64) error {
	_, err := c.client.DeleteCpu(ctx, &updatepb.IdRequest{Id: id})
	return err
}

// Motherboard methods
func (c Client) AddMotherboard(ctx context.Context, req *updatepb.MotherboardRequest) error {
	_, err := c.client.AddMotherboard(ctx, req)
	return err
}

func (c Client) UpdateMotherboard(ctx context.Context, req *updatepb.UpdateMotherboardRequest) error {
	_, err := c.client.UpdateMotherboard(ctx, req)
	return err
}

func (c Client) DeleteMotherboard(ctx context.Context, id int64) error {
	_, err := c.client.DeleteMotherboard(ctx, &updatepb.IdRequest{Id: id})
	return err
}

// RAM Kit methods
func (c Client) AddRAMKit(ctx context.Context, req *updatepb.RAMKitRequest) error {
	_, err := c.client.AddRAMKit(ctx, req)
	return err
}

func (c Client) UpdateRAMKit(ctx context.Context, req *updatepb.UpdateRAMKitRequest) error {
	_, err := c.client.UpdateRAMKit(ctx, req)
	return err
}

func (c Client) DeleteRAMKit(ctx context.Context, id int64) error {
	_, err := c.client.DeleteRAMKit(ctx, &updatepb.IdRequest{Id: id})
	return err
}

// Storage Drive methods
func (c Client) AddStorageDrive(ctx context.Context, req *updatepb.StorageDriveRequest) error {
	_, err := c.client.AddStorageDrive(ctx, req)
	return err
}

func (c Client) UpdateStorageDrive(ctx context.Context, req *updatepb.UpdateStorageDriveRequest) error {
	_, err := c.client.UpdateStorageDrive(ctx, req)
	return err
}

func (c Client) DeleteStorageDrive(ctx context.Context, id int64) error {
	_, err := c.client.DeleteStorageDrive(ctx, &updatepb.IdRequest{Id: id})
	return err
}

// User methods
func (c Client) AddUser(ctx context.Context, req *updatepb.UserRequest) error {
	_, err := c.client.AddUser(ctx, req)
	return err
}

func (c Client) UpdateUser(ctx context.Context, req *updatepb.UpdateUserRequest) error {
	_, err := c.client.UpdateUser(ctx, req)
	return err
}

func (c Client) DeleteUser(ctx context.Context, id int64) error {
	_, err := c.client.DeleteUser(ctx, &updatepb.IdRequest{Id: id})
	return err
}

// Shop methods
func (c Client) AddShop(ctx context.Context, req *updatepb.ShopRequest) error {
	_, err := c.client.AddShop(ctx, req)
	return err
}

func (c Client) UpdateShop(ctx context.Context, req *updatepb.UpdateShopRequest) error {
	_, err := c.client.UpdateShop(ctx, req)
	return err
}

func (c Client) DeleteShop(ctx context.Context, id int64) error {
	_, err := c.client.DeleteShop(ctx, &updatepb.IdRequest{Id: id})
	return err
}

// Product Offer methods
func (c Client) AddProductOffer(ctx context.Context, req *updatepb.ProductOfferRequest) error {
	_, err := c.client.AddProductOffer(ctx, req)
	return err
}

func (c Client) UpdateProductOffer(ctx context.Context, req *updatepb.UpdateProductOfferRequest) error {
	_, err := c.client.UpdateProductOffer(ctx, req)
	return err
}

func (c Client) DeleteProductOffer(ctx context.Context, id int64) error {
	_, err := c.client.DeleteProductOffer(ctx, &updatepb.IdRequest{Id: id})
	return err
}

// Assembly methods
func (c Client) AddAssembly(ctx context.Context, req *updatepb.AssemblyRequest) error {
	_, err := c.client.AddAssembly(ctx, req)
	return err
}

func (c Client) UpdateAssembly(ctx context.Context, req *updatepb.UpdateAssemblyRequest) error {
	_, err := c.client.UpdateAssembly(ctx, req)
	return err
}

func (c Client) DeleteAssembly(ctx context.Context, id int64) error {
	_, err := c.client.DeleteAssembly(ctx, &updatepb.IdRequest{Id: id})
	return err
}

// Cooler Socket methods (many-to-many)
func (c Client) AddCoolerSocket(ctx context.Context, req *updatepb.CoolerSocketRequest) error {
	_, err := c.client.AddCoolerSocket(ctx, req)
	return err
}

func (c Client) DeleteCoolerSocket(ctx context.Context, req *updatepb.DeleteCoolerSocketRequest) error {
	_, err := c.client.DeleteCoolerSocket(ctx, req)
	return err
}

func (c Client) DeleteCoolerSocketsByCooler(ctx context.Context, coolerID int64) error {
	_, err := c.client.DeleteCoolerSocketsByCooler(ctx, &updatepb.DeleteByCoolerRequest{CoolerId: coolerID})
	return err
}

// Case Form Factor Support methods (many-to-many)
func (c Client) AddCaseFormFactorSupport(ctx context.Context, req *updatepb.CaseFormFactorSupportRequest) error {
	_, err := c.client.AddCaseFormFactorSupport(ctx, req)
	return err
}

func (c Client) DeleteCaseFormFactorSupport(ctx context.Context, req *updatepb.DeleteCaseFormFactorSupportRequest) error {
	_, err := c.client.DeleteCaseFormFactorSupport(ctx, req)
	return err
}

func (c Client) DeleteCaseFormFactorSupportsByCase(ctx context.Context, caseID int64) error {
	_, err := c.client.DeleteCaseFormFactorSupportsByCase(ctx, &updatepb.DeleteByCaseRequest{CaseId: caseID})
	return err
}

// Assembly RAM Kit methods (many-to-many)
func (c Client) AddAssemblyRAMKit(ctx context.Context, req *updatepb.AssemblyRAMKitRequest) error {
	_, err := c.client.AddAssemblyRAMKit(ctx, req)
	return err
}

func (c Client) DeleteAssemblyRAMKit(ctx context.Context, req *updatepb.DeleteAssemblyRAMKitRequest) error {
	_, err := c.client.DeleteAssemblyRAMKit(ctx, req)
	return err
}

func (c Client) DeleteAssemblyRAMKitsByAssembly(ctx context.Context, assemblyID int64) error {
	_, err := c.client.DeleteAssemblyRAMKitsByAssembly(ctx, &updatepb.DeleteByAssemblyRequest{AssemblyId: assemblyID})
	return err
}

// Assembly Drive methods (many-to-many)
func (c Client) AddAssemblyDrive(ctx context.Context, req *updatepb.AssemblyDriveRequest) error {
	_, err := c.client.AddAssemblyDrive(ctx, req)
	return err
}

func (c Client) DeleteAssemblyDrive(ctx context.Context, req *updatepb.DeleteAssemblyDriveRequest) error {
	_, err := c.client.DeleteAssemblyDrive(ctx, req)
	return err
}

func (c Client) DeleteAssemblyDrivesByAssembly(ctx context.Context, assemblyID int64) error {
	_, err := c.client.DeleteAssemblyDrivesByAssembly(ctx, &updatepb.DeleteByAssemblyRequest{AssemblyId: assemblyID})
	return err
}

// Вспомогательные методы для создания запросов
func (c Client) NewGPURequest() *updatepb.GpuRequest {
	return &updatepb.GpuRequest{}
}

func (c Client) NewCPUCoolerRequest() *updatepb.CpuCoolerRequest {
	return &updatepb.CpuCoolerRequest{}
}

func (c Client) NewCaseRequest() *updatepb.CaseRequest {
	return &updatepb.CaseRequest{}
}

func (c Client) NewPSURequest() *updatepb.PsuRequest {
	return &updatepb.PsuRequest{}
}

func (c Client) NewCPUSocketRequest() *updatepb.CpuSocketRequest {
	return &updatepb.CpuSocketRequest{}
}

func (c Client) NewMotherboardFormFactorRequest() *updatepb.MotherboardFormFactorRequest {
	return &updatepb.MotherboardFormFactorRequest{}
}

func (c Client) NewCPURequest() *updatepb.CpuRequest {
	return &updatepb.CpuRequest{}
}

func (c Client) NewMotherboardRequest() *updatepb.MotherboardRequest {
	return &updatepb.MotherboardRequest{}
}

func (c Client) NewRAMKitRequest() *updatepb.RAMKitRequest {
	return &updatepb.RAMKitRequest{}
}

func (c Client) NewStorageDriveRequest() *updatepb.StorageDriveRequest {
	return &updatepb.StorageDriveRequest{}
}

func (c Client) NewUserRequest() *updatepb.UserRequest {
	return &updatepb.UserRequest{}
}

func (c Client) NewShopRequest() *updatepb.ShopRequest {
	return &updatepb.ShopRequest{}
}

func (c Client) NewProductOfferRequest() *updatepb.ProductOfferRequest {
	return &updatepb.ProductOfferRequest{}
}

func (c Client) NewAssemblyRequest() *updatepb.AssemblyRequest {
	return &updatepb.AssemblyRequest{}
}
