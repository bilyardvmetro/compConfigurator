package grpc

import (
	"context"

	searchpb "comp_config.com/services/proto/search"
	"comp_config.com/services/search/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func NewServer(service core.Searcher) *Server {
	return &Server{service: service}
}

type Server struct {
	searchpb.UnimplementedSearchServer
	service core.Searcher
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// GPU methods
func (s *Server) GetGPU(ctx context.Context, req *searchpb.GetRequest) (*searchpb.GPU, error) {
	gpu, err := s.service.GetGPU(ctx, int(req.Id))
	if err != nil {
		if err == core.ErrNotFound {
			return nil, status.Error(codes.NotFound, "GPU not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertGPUToProto(gpu), nil
}

func (s *Server) ListGPUs(ctx context.Context, req *searchpb.ListRequest) (*searchpb.GPUList, error) {
	opts := convertListRequestToOptions(req)
	result, err := s.service.ListGPUs(ctx, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertGPUListToProto(result), nil
}

// CPU methods
func (s *Server) GetCPU(ctx context.Context, req *searchpb.GetRequest) (*searchpb.CPU, error) {
	cpu, err := s.service.GetCPU(ctx, int(req.Id))
	if err != nil {
		if err == core.ErrNotFound {
			return nil, status.Error(codes.NotFound, "CPU not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertCPUToProto(cpu), nil
}

func (s *Server) ListCPUs(ctx context.Context, req *searchpb.ListRequest) (*searchpb.CPUList, error) {
	opts := convertListRequestToOptions(req)
	result, err := s.service.ListCPUs(ctx, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertCPUListToProto(result), nil
}

// Motherboard methods
func (s *Server) GetMotherboard(ctx context.Context, req *searchpb.GetRequest) (*searchpb.Motherboard, error) {
	mb, err := s.service.GetMotherboard(ctx, int(req.Id))
	if err != nil {
		if err == core.ErrNotFound {
			return nil, status.Error(codes.NotFound, "Motherboard not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertMotherboardToProto(mb), nil
}

func (s *Server) ListMotherboards(ctx context.Context, req *searchpb.ListRequest) (*searchpb.MotherboardList, error) {
	opts := convertListRequestToOptions(req)
	result, err := s.service.ListMotherboards(ctx, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertMotherboardListToProto(result), nil
}

// RAM Kit methods
func (s *Server) GetRAMKit(ctx context.Context, req *searchpb.GetRequest) (*searchpb.RAMKit, error) {
	ram, err := s.service.GetRAMKit(ctx, int(req.Id))
	if err != nil {
		if err == core.ErrNotFound {
			return nil, status.Error(codes.NotFound, "RAM kit not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertRAMKitToProto(ram), nil
}

func (s *Server) ListRAMKits(ctx context.Context, req *searchpb.ListRequest) (*searchpb.RAMKitList, error) {
	opts := convertListRequestToOptions(req)
	result, err := s.service.ListRAMKits(ctx, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertRAMKitListToProto(result), nil
}

// PSU methods
func (s *Server) GetPSU(ctx context.Context, req *searchpb.GetRequest) (*searchpb.PSU, error) {
	psu, err := s.service.GetPSU(ctx, int(req.Id))
	if err != nil {
		if err == core.ErrNotFound {
			return nil, status.Error(codes.NotFound, "PSU not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertPSUToProto(psu), nil
}

func (s *Server) ListPSUs(ctx context.Context, req *searchpb.ListRequest) (*searchpb.PSUList, error) {
	opts := convertListRequestToOptions(req)
	result, err := s.service.ListPSUs(ctx, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertPSUListToProto(result), nil
}

// Case methods
func (s *Server) GetCase(ctx context.Context, req *searchpb.GetRequest) (*searchpb.Case, error) {
	c, err := s.service.GetCase(ctx, int(req.Id))
	if err != nil {
		if err == core.ErrNotFound {
			return nil, status.Error(codes.NotFound, "Case not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertCaseToProto(c), nil
}

func (s *Server) ListCases(ctx context.Context, req *searchpb.ListRequest) (*searchpb.CaseList, error) {
	opts := convertListRequestToOptions(req)
	result, err := s.service.ListCases(ctx, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertCaseListToProto(result), nil
}

// CPU Cooler methods
func (s *Server) GetCpuCooler(ctx context.Context, req *searchpb.GetRequest) (*searchpb.CpuCooler, error) {
	cooler, err := s.service.GetCpuCooler(ctx, int(req.Id))
	if err != nil {
		if err == core.ErrNotFound {
			return nil, status.Error(codes.NotFound, "CPU cooler not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertCpuCoolerToProto(cooler), nil
}

func (s *Server) ListCpuCoolers(ctx context.Context, req *searchpb.ListRequest) (*searchpb.CpuCoolerList, error) {
	opts := convertListRequestToOptions(req)
	result, err := s.service.ListCpuCoolers(ctx, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertCpuCoolerListToProto(result), nil
}

// Storage Drive methods
func (s *Server) GetStorageDrive(ctx context.Context, req *searchpb.GetRequest) (*searchpb.StorageDrive, error) {
	drive, err := s.service.GetStorageDrive(ctx, int(req.Id))
	if err != nil {
		if err == core.ErrNotFound {
			return nil, status.Error(codes.NotFound, "Storage drive not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertStorageDriveToProto(drive), nil
}

func (s *Server) ListStorageDrives(ctx context.Context, req *searchpb.ListRequest) (*searchpb.StorageDriveList, error) {
	opts := convertListRequestToOptions(req)
	result, err := s.service.ListStorageDrives(ctx, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertStorageDriveListToProto(result), nil
}

// CPU Socket methods
func (s *Server) GetCpuSocket(ctx context.Context, req *searchpb.GetByCodeRequest) (*searchpb.CpuSocket, error) {
	socket, err := s.service.GetCpuSocket(ctx, req.Code)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, status.Error(codes.NotFound, "CPU socket not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertCpuSocketToProto(socket), nil
}

func (s *Server) ListCpuSockets(ctx context.Context, req *searchpb.ListRequest) (*searchpb.CpuSocketList, error) {
	opts := convertListRequestToOptions(req)
	result, err := s.service.ListCpuSockets(ctx, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertCpuSocketListToProto(result), nil
}

// Motherboard Form Factor methods
func (s *Server) GetMotherboardFormFactor(ctx context.Context, req *searchpb.GetByCodeRequest) (*searchpb.MotherboardFormFactor, error) {
	ff, err := s.service.GetMotherboardFormFactor(ctx, req.Code)
	if err != nil {
		if err == core.ErrNotFound {
			return nil, status.Error(codes.NotFound, "Motherboard form factor not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertMotherboardFormFactorToProto(ff), nil
}

func (s *Server) ListMotherboardFormFactors(ctx context.Context, req *searchpb.ListRequest) (*searchpb.MotherboardFormFactorList, error) {
	opts := convertListRequestToOptions(req)
	result, err := s.service.ListMotherboardFormFactors(ctx, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertMotherboardFormFactorListToProto(result), nil
}

// User methods
func (s *Server) GetUser(ctx context.Context, req *searchpb.GetRequest) (*searchpb.User, error) {
	user, err := s.service.GetUser(ctx, int(req.Id))
	if err != nil {
		if err == core.ErrNotFound {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertUserToProto(user), nil
}

func (s *Server) ListUsers(ctx context.Context, req *searchpb.ListRequest) (*searchpb.UserList, error) {
	opts := convertListRequestToOptions(req)
	result, err := s.service.ListUsers(ctx, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertUserListToProto(result), nil
}

// Shop methods
func (s *Server) GetShop(ctx context.Context, req *searchpb.GetRequest) (*searchpb.Shop, error) {
	shop, err := s.service.GetShop(ctx, int(req.Id))
	if err != nil {
		if err == core.ErrNotFound {
			return nil, status.Error(codes.NotFound, "Shop not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertShopToProto(shop), nil
}

func (s *Server) ListShops(ctx context.Context, req *searchpb.ListRequest) (*searchpb.ShopList, error) {
	opts := convertListRequestToOptions(req)
	result, err := s.service.ListShops(ctx, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertShopListToProto(result), nil
}

// Product Offer methods
func (s *Server) GetProductOffer(ctx context.Context, req *searchpb.GetRequest) (*searchpb.ProductOffer, error) {
	offer, err := s.service.GetProductOffer(ctx, int(req.Id))
	if err != nil {
		if err == core.ErrNotFound {
			return nil, status.Error(codes.NotFound, "Product offer not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertProductOfferToProto(offer), nil
}

func (s *Server) ListProductOffers(ctx context.Context, req *searchpb.ListRequest) (*searchpb.ProductOfferList, error) {
	opts := convertListRequestToOptions(req)
	result, err := s.service.ListProductOffers(ctx, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertProductOfferListToProto(result), nil
}

// Assembly methods
func (s *Server) GetAssembly(ctx context.Context, req *searchpb.GetRequest) (*searchpb.Assembly, error) {
	assembly, err := s.service.GetAssembly(ctx, int(req.Id))
	if err != nil {
		if err == core.ErrNotFound {
			return nil, status.Error(codes.NotFound, "Assembly not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertAssemblyToProto(assembly), nil
}

func (s *Server) ListAssemblies(ctx context.Context, req *searchpb.ListRequest) (*searchpb.AssemblyList, error) {
	opts := convertListRequestToOptions(req)
	result, err := s.service.ListAssemblies(ctx, opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertAssemblyListToProto(result), nil
}

func (s *Server) ListAssembliesByUser(ctx context.Context, req *searchpb.ListByUserRequest) (*searchpb.AssemblyList, error) {
	opts := convertListRequestToOptions(req.ListRequest)
	result, err := s.service.ListAssembliesByUser(ctx, int(req.UserId), opts)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertAssemblyListToProto(result), nil
}

// Compatibility methods
func (s *Server) CheckCompatibility(ctx context.Context, req *searchpb.CheckCompatibilityRequest) (*searchpb.ValidationResult, error) {
	result, err := s.service.CheckAssemblyCompatibility(ctx, int(req.AssemblyId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertValidationResultToProto(result), nil
}

// Search methods
func (s *Server) SearchComponents(ctx context.Context, req *searchpb.SearchRequest) (*searchpb.SearchResponse, error) {
	results, totalCount, err := s.service.SearchComponents(ctx, req.Query, req.ComponentTypes, int(req.Limit), int(req.Offset), req.Filters)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertSearchResultsToProto(results, totalCount, int(req.Limit), int(req.Offset)), nil
}

func (s *Server) GetComponentPrices(ctx context.Context, req *searchpb.ComponentPricesRequest) (*searchpb.PriceList, error) {
	prices, err := s.service.GetComponentPrices(ctx, req.ComponentType, int(req.ComponentId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertPriceListToProto(req.ComponentType, int(req.ComponentId), prices), nil
}

func (s *Server) GetBestOffers(ctx context.Context, req *searchpb.OffersRequest) (*searchpb.OfferList, error) {
	offers, err := s.service.GetBestOffers(ctx, req.ComponentType, int(req.ComponentId), int(req.Limit), req.CheapestFirst)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return convertOfferListToProto(offers), nil
}
