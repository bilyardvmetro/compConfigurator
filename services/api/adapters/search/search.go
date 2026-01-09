package search

import (
	"context"
	"log/slog"

	searchpb "comp_config.com/services/proto/search"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	log    *slog.Logger
	client searchpb.SearchClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		client: searchpb.NewSearchClient(conn),
		log:    log,
	}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, nil)
	return err
}

// Get methods
func (c *Client) GetGPU(ctx context.Context, id int32) (*searchpb.GPU, error) {
	return c.client.GetGPU(ctx, &searchpb.GetRequest{Id: id})
}

func (c *Client) GetCPU(ctx context.Context, id int32) (*searchpb.CPU, error) {
	return c.client.GetCPU(ctx, &searchpb.GetRequest{Id: id})
}

func (c *Client) GetMotherboard(ctx context.Context, id int32) (*searchpb.Motherboard, error) {
	return c.client.GetMotherboard(ctx, &searchpb.GetRequest{Id: id})
}

func (c *Client) GetRAMKit(ctx context.Context, id int32) (*searchpb.RAMKit, error) {
	return c.client.GetRAMKit(ctx, &searchpb.GetRequest{Id: id})
}

func (c *Client) GetPSU(ctx context.Context, id int32) (*searchpb.PSU, error) {
	return c.client.GetPSU(ctx, &searchpb.GetRequest{Id: id})
}

func (c *Client) GetCase(ctx context.Context, id int32) (*searchpb.Case, error) {
	return c.client.GetCase(ctx, &searchpb.GetRequest{Id: id})
}

func (c *Client) GetCpuCooler(ctx context.Context, id int32) (*searchpb.CpuCooler, error) {
	return c.client.GetCpuCooler(ctx, &searchpb.GetRequest{Id: id})
}

func (c *Client) GetStorageDrive(ctx context.Context, id int32) (*searchpb.StorageDrive, error) {
	return c.client.GetStorageDrive(ctx, &searchpb.GetRequest{Id: id})
}

func (c *Client) GetCpuSocket(ctx context.Context, code string) (*searchpb.CpuSocket, error) {
	return c.client.GetCpuSocket(ctx, &searchpb.GetByCodeRequest{Code: code})
}

func (c *Client) GetMotherboardFormFactor(ctx context.Context, code string) (*searchpb.MotherboardFormFactor, error) {
	return c.client.GetMotherboardFormFactor(ctx, &searchpb.GetByCodeRequest{Code: code})
}

func (c *Client) GetUser(ctx context.Context, id int32) (*searchpb.User, error) {
	return c.client.GetUser(ctx, &searchpb.GetRequest{Id: id})
}

func (c *Client) GetShop(ctx context.Context, id int32) (*searchpb.Shop, error) {
	return c.client.GetShop(ctx, &searchpb.GetRequest{Id: id})
}

func (c *Client) GetProductOffer(ctx context.Context, id int32) (*searchpb.ProductOffer, error) {
	return c.client.GetProductOffer(ctx, &searchpb.GetRequest{Id: id})
}

func (c *Client) GetAssembly(ctx context.Context, id int32) (*searchpb.Assembly, error) {
	return c.client.GetAssembly(ctx, &searchpb.GetRequest{Id: id})
}

// List methods
func (c *Client) ListGPUs(ctx context.Context, req *searchpb.ListRequest) (*searchpb.GPUList, error) {
	return c.client.ListGPUs(ctx, req)
}

func (c *Client) ListCPUs(ctx context.Context, req *searchpb.ListRequest) (*searchpb.CPUList, error) {
	return c.client.ListCPUs(ctx, req)
}

func (c *Client) ListMotherboards(ctx context.Context, req *searchpb.ListRequest) (*searchpb.MotherboardList, error) {
	return c.client.ListMotherboards(ctx, req)
}

func (c *Client) ListRAMKits(ctx context.Context, req *searchpb.ListRequest) (*searchpb.RAMKitList, error) {
	return c.client.ListRAMKits(ctx, req)
}

func (c *Client) ListPSUs(ctx context.Context, req *searchpb.ListRequest) (*searchpb.PSUList, error) {
	return c.client.ListPSUs(ctx, req)
}

func (c *Client) ListCases(ctx context.Context, req *searchpb.ListRequest) (*searchpb.CaseList, error) {
	return c.client.ListCases(ctx, req)
}

func (c *Client) ListCpuCoolers(ctx context.Context, req *searchpb.ListRequest) (*searchpb.CpuCoolerList, error) {
	return c.client.ListCpuCoolers(ctx, req)
}

func (c *Client) ListStorageDrives(ctx context.Context, req *searchpb.ListRequest) (*searchpb.StorageDriveList, error) {
	return c.client.ListStorageDrives(ctx, req)
}

func (c *Client) ListCpuSockets(ctx context.Context, req *searchpb.ListRequest) (*searchpb.CpuSocketList, error) {
	return c.client.ListCpuSockets(ctx, req)
}

func (c *Client) ListMotherboardFormFactors(ctx context.Context, req *searchpb.ListRequest) (*searchpb.MotherboardFormFactorList, error) {
	return c.client.ListMotherboardFormFactors(ctx, req)
}

func (c *Client) ListUsers(ctx context.Context, req *searchpb.ListRequest) (*searchpb.UserList, error) {
	return c.client.ListUsers(ctx, req)
}

func (c *Client) ListShops(ctx context.Context, req *searchpb.ListRequest) (*searchpb.ShopList, error) {
	return c.client.ListShops(ctx, req)
}

func (c *Client) ListProductOffers(ctx context.Context, req *searchpb.ListRequest) (*searchpb.ProductOfferList, error) {
	return c.client.ListProductOffers(ctx, req)
}

func (c *Client) ListAssemblies(ctx context.Context, req *searchpb.ListRequest) (*searchpb.AssemblyList, error) {
	return c.client.ListAssemblies(ctx, req)
}

func (c *Client) ListAssembliesByUser(ctx context.Context, userID int32, req *searchpb.ListRequest) (*searchpb.AssemblyList, error) {
	return c.client.ListAssembliesByUser(ctx, &searchpb.ListByUserRequest{
		UserId:      userID,
		ListRequest: req,
	})
}

// Search methods
func (c *Client) SearchComponents(ctx context.Context, req *searchpb.SearchRequest) (*searchpb.SearchResponse, error) {
	return c.client.SearchComponents(ctx, req)
}

func (c *Client) SearchAssemblies(ctx context.Context, req *searchpb.SearchAssembliesRequest) (*searchpb.AssemblyList, error) {
	return c.client.SearchAssemblies(ctx, req)
}

// Compatibility methods
func (c *Client) CheckCompatibility(ctx context.Context, assemblyID int32) (*searchpb.ValidationResult, error) {
	return c.client.CheckCompatibility(ctx, &searchpb.CheckCompatibilityRequest{AssemblyId: assemblyID})
}

func (c *Client) GetCompatibleComponents(ctx context.Context, req *searchpb.CompatibilityRequest) (*searchpb.CompatibleComponents, error) {
	return c.client.GetCompatibleComponents(ctx, req)
}

func (c *Client) RecommendComponents(ctx context.Context, req *searchpb.RecommendationRequest) (*searchpb.RecommendationResponse, error) {
	return c.client.RecommendComponents(ctx, req)
}

// Price methods
func (c *Client) GetComponentPrices(ctx context.Context, componentType string, componentID int32, shopID ...int32) (*searchpb.PriceList, error) {
	req := &searchpb.ComponentPricesRequest{
		ComponentType: componentType,
		ComponentId:   componentID,
	}
	if len(shopID) > 0 {
		shopIDVal := int32(shopID[0])
		req.ShopId = &shopIDVal
	}
	return c.client.GetComponentPrices(ctx, req)
}

func (c *Client) GetBestOffers(ctx context.Context, componentType string, componentID int32, limit int32, cheapestFirst bool) (*searchpb.OfferList, error) {
	return c.client.GetBestOffers(ctx, &searchpb.OffersRequest{
		ComponentType: componentType,
		ComponentId:   componentID,
		Limit:         limit,
		CheapestFirst: cheapestFirst,
	})
}

// Helper methods to create requests
func (c *Client) NewListRequest(page, pageSize int32, sortBy, sortOrder string, filters []*searchpb.Filter) *searchpb.ListRequest {
	return &searchpb.ListRequest{
		Page:      page,
		PageSize:  pageSize,
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Filters:   filters,
	}
}

func (c *Client) NewFilter(field, operator, value string) *searchpb.Filter {
	return &searchpb.Filter{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}

func (c *Client) NewSearchRequest(query string, componentTypes []string, limit, offset int32, filters map[string]string) *searchpb.SearchRequest {
	return &searchpb.SearchRequest{
		Query:          query,
		ComponentTypes: componentTypes,
		Limit:          limit,
		Offset:         offset,
		Filters:        filters,
	}
}

func (c *Client) NewSearchAssembliesRequest(query string, userID *int32, onlyPublic bool, pagination *searchpb.ListRequest) *searchpb.SearchAssembliesRequest {
	req := &searchpb.SearchAssembliesRequest{
		Query:      query,
		OnlyPublic: onlyPublic,
		Pagination: pagination,
	}
	if userID != nil {
		req.UserId = userID
	}
	return req
}

func (c *Client) NewCompatibilityRequest() *searchpb.CompatibilityRequest {
	return &searchpb.CompatibilityRequest{}
}

func (c *Client) NewRecommendationRequest(purpose string, budget int32, priorities []string) *searchpb.RecommendationRequest {
	return &searchpb.RecommendationRequest{
		Purpose:    purpose,
		Budget:     budget,
		Priorities: priorities,
	}
}
