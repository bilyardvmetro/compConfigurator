package grpc

import (
	"strconv"

	searchpb "comp_config.com/services/proto/search"
	"comp_config.com/services/search/core"
)

func convertListRequestToOptions(req *searchpb.ListRequest) *core.ListOptions {
	if req == nil {
		return &core.ListOptions{
			Page:      1,
			PageSize:  50,
			SortBy:    "id",
			SortOrder: "asc",
		}
	}

	opts := &core.ListOptions{
		Page:      int(req.Page),
		PageSize:  int(req.PageSize),
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	if len(req.Filters) > 0 {
		opts.Filters = make([]core.Filter, len(req.Filters))
		for i, f := range req.Filters {
			opts.Filters[i] = core.Filter{
				Field:    f.Field,
				Operator: f.Operator,
				Value:    f.Value,
			}
		}
	}

	return opts
}

// Конвертация GPU
func convertGPUToProto(gpu *core.GPU) *searchpb.GPU {
	if gpu == nil {
		return nil
	}
	return &searchpb.GPU{
		Id:                int32(gpu.ID),
		Name:              gpu.Name,
		Brand:             gpu.Brand,
		Interface:         gpu.Interface,
		PcieVersion:       strconv.Itoa(gpu.PcieVersion),
		PcieLanesRequired: int32(gpu.PcieLanesRequired),
		TdpWatt:           optionalInt32(&gpu.TdpWatt),
		PowerConnectors:   optionalString(&gpu.PowerConnectors),
		LengthMm:          optionalInt32(&gpu.LengthMm),
		WidthSlots:        optionalFloat(&gpu.WidthSlots),
		HeightMm:          optionalInt32(&gpu.HeightMm),
	}
}

func convertGPUListToProto(result *core.ListResult[core.GPU]) *searchpb.GPUList {
	if result == nil {
		return &searchpb.GPUList{}
	}

	items := make([]*searchpb.GPU, len(result.Items))
	for i, gpu := range result.Items {
		items[i] = convertGPUToProto(&gpu)
	}

	return &searchpb.GPUList{
		Items:      items,
		TotalCount: int32(result.TotalCount),
		Page:       int32(result.Page),
		PageSize:   int32(result.PageSize),
		TotalPages: int32(result.TotalPages),
	}
}

// Конвертация CPU
func convertCPUToProto(cpu *core.CPU) *searchpb.CPU {
	if cpu == nil {
		return nil
	}
	return &searchpb.CPU{
		Id:                     int32(cpu.ID),
		Name:                   cpu.Name,
		Brand:                  cpu.Brand,
		SocketCode:             cpu.SocketCode,
		Architecture:           cpu.Architecture,
		CoreCount:              int32(cpu.CoreCount),
		ThreadCount:            int32(cpu.ThreadCount),
		BaseClockMhz:           int32(cpu.BaseClockMhz),
		BoostClockMhz:          optionalInt32(cpu.BoostClockMhz),
		TdpWatt:                optionalInt32(&cpu.TdpWatt),
		HasIntegratedGpu:       cpu.HasIntegratedGpu,
		SupportedRamType:       cpu.SupportedRamType,
		SupportedRamFreqMaxMhz: optionalInt32(cpu.SupportedRamFreqMaxMhz),
		MemoryChannels:         int32(cpu.MemoryChannels),
		PcieVersion:            strconv.Itoa(cpu.PcieVersion),
		PcieLanesTotal:         int32(cpu.PcieLanesTotal),
	}
}

func convertCPUListToProto(result *core.ListResult[core.CPU]) *searchpb.CPUList {
	if result == nil {
		return &searchpb.CPUList{}
	}

	items := make([]*searchpb.CPU, len(result.Items))
	for i, cpu := range result.Items {
		items[i] = convertCPUToProto(&cpu)
	}

	return &searchpb.CPUList{
		Items:      items,
		TotalCount: int32(result.TotalCount),
		Page:       int32(result.Page),
		PageSize:   int32(result.PageSize),
		TotalPages: int32(result.TotalPages),
	}
}

// Конвертация Motherboard
func convertMotherboardToProto(mb *core.Motherboard) *searchpb.Motherboard {
	if mb == nil {
		return nil
	}
	return &searchpb.Motherboard{
		Id:                    int32(mb.ID),
		Name:                  mb.Name,
		Brand:                 mb.Brand,
		SocketCode:            mb.SocketCode,
		FormFactorCode:        mb.FormFactorCode,
		Chipset:               mb.Chipset,
		RamType:               mb.RamType,
		RamSlots:              int32(mb.RamSlots),
		RamCapacityMaxGb:      int32(*mb.RamCapacityMaxGb),
		RamFreqMaxMhz:         int32(*mb.RamFreqMaxMhz),
		PcieX16SlotsCount:     int32(mb.PcieX16SlotsCount),
		PcieVersionMax:        strconv.Itoa(mb.PcieVersionMax),
		M2SlotsCount:          int32(mb.M2SlotsCount),
		SataPortsCount:        int32(mb.SataPortsCount),
		PsuMainConnectorType:  mb.PsuMainConnectorType,
		CpuPowerConnectorType: mb.CpuPowerConnectorType,
	}
}

func convertMotherboardListToProto(result *core.ListResult[core.Motherboard]) *searchpb.MotherboardList {
	if result == nil {
		return &searchpb.MotherboardList{}
	}

	items := make([]*searchpb.Motherboard, len(result.Items))
	for i, mb := range result.Items {
		items[i] = convertMotherboardToProto(&mb)
	}

	return &searchpb.MotherboardList{
		Items:      items,
		TotalCount: int32(result.TotalCount),
		Page:       int32(result.Page),
		PageSize:   int32(result.PageSize),
		TotalPages: int32(result.TotalPages),
	}
}

// Конвертация Case
func convertCaseToProto(c *core.Case) *searchpb.Case {
	if c == nil {
		return nil
	}
	return &searchpb.Case{
		Id:                 int32(c.ID),
		Name:               c.Name,
		Brand:              c.Brand,
		PsuFormFactor:      c.PsuFormFactor,
		MaxGpuLengthMm:     optionalInt32(c.MaxGpuLengthMm),
		MaxGpuWidthSlots:   optionalFloat(c.MaxGpuWidthSlots),
		MaxCoolerHeightMm:  optionalInt32(c.MaxCoolerHeightMm),
		MaxPsuLengthMm:     optionalInt32(c.MaxPsuLengthMm),
		DriveBays_3_5Count: int32(*c.DriveBays3_5Count),
		DriveBays_2_5Count: int32(*c.DriveBays2_5Count),
	}
}

func convertCaseListToProto(result *core.ListResult[core.Case]) *searchpb.CaseList {
	if result == nil {
		return &searchpb.CaseList{}
	}

	items := make([]*searchpb.Case, len(result.Items))
	for i, c := range result.Items {
		items[i] = convertCaseToProto(&c)
	}

	return &searchpb.CaseList{
		Items:      items,
		TotalCount: int32(result.TotalCount),
		Page:       int32(result.Page),
		PageSize:   int32(result.PageSize),
		TotalPages: int32(result.TotalPages),
	}
}

// Конвертация CPU Cooler
func convertCpuCoolerToProto(cooler *core.CpuCooler) *searchpb.CpuCooler {
	if cooler == nil {
		return nil
	}
	return &searchpb.CpuCooler{
		Id:             int32(cooler.ID),
		Name:           cooler.Name,
		Brand:          cooler.Brand,
		CoolingType:    cooler.CoolingType,
		TdpLimitWatt:   optionalInt32(&cooler.TdpLimitWatt),
		HeightMm:       optionalInt32(&cooler.HeightMm),
		FanCount:       int32(cooler.FanCount),
		FanControlType: cooler.FanControlType,
	}
}

func convertCpuCoolerListToProto(result *core.ListResult[core.CpuCooler]) *searchpb.CpuCoolerList {
	if result == nil {
		return &searchpb.CpuCoolerList{}
	}

	items := make([]*searchpb.CpuCooler, len(result.Items))
	for i, cooler := range result.Items {
		items[i] = convertCpuCoolerToProto(&cooler)
	}

	return &searchpb.CpuCoolerList{
		Items:      items,
		TotalCount: int32(result.TotalCount),
		Page:       int32(result.Page),
		PageSize:   int32(result.PageSize),
		TotalPages: int32(result.TotalPages),
	}
}

// Конвертация PSU
func convertPSUToProto(psu *core.PSU) *searchpb.PSU {
	if psu == nil {
		return nil
	}
	return &searchpb.PSU{
		Id:                         int32(psu.ID),
		Name:                       psu.Name,
		Brand:                      psu.Brand,
		PowerWatt:                  int32(psu.PowerWatt),
		FormFactor:                 psu.FormFactor,
		EfficiencyRating:           psu.EfficiencyRating,
		PcieConnectors_6_8PinCount: int32(psu.PcieConnectors6_8pinCount),
		Cpu_8PinConnectorsCount:    int32(psu.Cpu8pinConnectorsCount),
		SataConnectorsCount:        int32(psu.SataConnectorsCount),
		MolexConnectorsCount:       int32(psu.MolexConnectorsCount),
		LengthMm:                   optionalInt32(psu.LengthMm),
	}
}

func convertPSUListToProto(result *core.ListResult[core.PSU]) *searchpb.PSUList {
	if result == nil {
		return &searchpb.PSUList{}
	}

	items := make([]*searchpb.PSU, len(result.Items))
	for i, psu := range result.Items {
		items[i] = convertPSUToProto(&psu)
	}

	return &searchpb.PSUList{
		Items:      items,
		TotalCount: int32(result.TotalCount),
		Page:       int32(result.Page),
		PageSize:   int32(result.PageSize),
		TotalPages: int32(result.TotalPages),
	}
}

// Конвертация RAM Kit
func convertRAMKitToProto(ram *core.RAMKit) *searchpb.RAMKit {
	if ram == nil {
		return nil
	}
	return &searchpb.RAMKit{
		Id:               int32(ram.ID),
		Name:             ram.Name,
		Brand:            ram.Brand,
		RamType:          ram.RamType,
		ModuleCapacityGb: int32(ram.ModuleCapacityGb),
		ModuleCount:      int32(ram.ModuleCount),
		TotalCapacityGb:  int32(ram.TotalCapacityGb),
		FreqMhz:          int32(ram.FreqMhz),
		Timings:          ram.Timings,
		VoltageV:         float32(*ram.VoltageV),
		FormFactor:       ram.FormFactor,
	}
}

func convertRAMKitListToProto(result *core.ListResult[core.RAMKit]) *searchpb.RAMKitList {
	if result == nil {
		return &searchpb.RAMKitList{}
	}

	items := make([]*searchpb.RAMKit, len(result.Items))
	for i, ram := range result.Items {
		items[i] = convertRAMKitToProto(&ram)
	}

	return &searchpb.RAMKitList{
		Items:      items,
		TotalCount: int32(result.TotalCount),
		Page:       int32(result.Page),
		PageSize:   int32(result.PageSize),
		TotalPages: int32(result.TotalPages),
	}
}

// Конвертация Storage Drive
func convertStorageDriveToProto(drive *core.StorageDrive) *searchpb.StorageDrive {
	if drive == nil {
		return nil
	}
	return &searchpb.StorageDrive{
		Id:         int32(drive.ID),
		Name:       drive.Name,
		Brand:      drive.Brand,
		DriveType:  drive.DriveType,
		FormFactor: drive.FormFactor,
		Interface:  drive.Interface,
		CapacityGb: int32(drive.CapacityGb),
	}
}

func convertStorageDriveListToProto(result *core.ListResult[core.StorageDrive]) *searchpb.StorageDriveList {
	if result == nil {
		return &searchpb.StorageDriveList{}
	}

	items := make([]*searchpb.StorageDrive, len(result.Items))
	for i, drive := range result.Items {
		items[i] = convertStorageDriveToProto(&drive)
	}

	return &searchpb.StorageDriveList{
		Items:      items,
		TotalCount: int32(result.TotalCount),
		Page:       int32(result.Page),
		PageSize:   int32(result.PageSize),
		TotalPages: int32(result.TotalPages),
	}
}

// Конвертация CPU Socket
func convertCpuSocketToProto(socket *core.CpuSocket) *searchpb.CpuSocket {
	if socket == nil {
		return nil
	}
	return &searchpb.CpuSocket{
		Code:        socket.Code,
		Description: socket.Description,
	}
}

func convertCpuSocketListToProto(result *core.ListResult[core.CpuSocket]) *searchpb.CpuSocketList {
	if result == nil {
		return &searchpb.CpuSocketList{}
	}

	items := make([]*searchpb.CpuSocket, len(result.Items))
	for i, socket := range result.Items {
		items[i] = convertCpuSocketToProto(&socket)
	}

	return &searchpb.CpuSocketList{
		Items:      items,
		TotalCount: int32(result.TotalCount),
		Page:       int32(result.Page),
		PageSize:   int32(result.PageSize),
		TotalPages: int32(result.TotalPages),
	}
}

// Конвертация Motherboard Form Factor
func convertMotherboardFormFactorToProto(ff *core.MotherboardFormFactor) *searchpb.MotherboardFormFactor {
	if ff == nil {
		return nil
	}
	return &searchpb.MotherboardFormFactor{
		Code:        ff.Code,
		Description: ff.Description,
		WidthMm:     int32(*ff.WidthMm),
		HeightMm:    int32(*ff.HeightMm),
	}
}

func convertMotherboardFormFactorListToProto(result *core.ListResult[core.MotherboardFormFactor]) *searchpb.MotherboardFormFactorList {
	if result == nil {
		return &searchpb.MotherboardFormFactorList{}
	}

	items := make([]*searchpb.MotherboardFormFactor, len(result.Items))
	for i, ff := range result.Items {
		items[i] = convertMotherboardFormFactorToProto(&ff)
	}

	return &searchpb.MotherboardFormFactorList{
		Items:      items,
		TotalCount: int32(result.TotalCount),
		Page:       int32(result.Page),
		PageSize:   int32(result.PageSize),
		TotalPages: int32(result.TotalPages),
	}
}

// Конвертация User
func convertUserToProto(user *core.User) *searchpb.User {
	if user == nil {
		return nil
	}
	return &searchpb.User{
		Id:           int32(user.ID),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Nickname:     *user.Nickname,
		AvatarUrl:    optionalString(user.AvatarUrl),
		IsAdmin:      user.IsAdmin,
	}
}

func convertUserListToProto(result *core.ListResult[core.User]) *searchpb.UserList {
	if result == nil {
		return &searchpb.UserList{}
	}

	items := make([]*searchpb.User, len(result.Items))
	for i, user := range result.Items {
		items[i] = convertUserToProto(&user)
	}

	return &searchpb.UserList{
		Items:      items,
		TotalCount: int32(result.TotalCount),
		Page:       int32(result.Page),
		PageSize:   int32(result.PageSize),
		TotalPages: int32(result.TotalPages),
	}
}

// Конвертация Shop
func convertShopToProto(shop *core.Shop) *searchpb.Shop {
	if shop == nil {
		return nil
	}
	return &searchpb.Shop{
		Id:   int32(shop.ID),
		Name: shop.Name,
		Url:  *shop.Url,
	}
}

func convertShopListToProto(result *core.ListResult[core.Shop]) *searchpb.ShopList {
	if result == nil {
		return &searchpb.ShopList{}
	}

	items := make([]*searchpb.Shop, len(result.Items))
	for i, shop := range result.Items {
		items[i] = convertShopToProto(&shop)
	}

	return &searchpb.ShopList{
		Items:      items,
		TotalCount: int32(result.TotalCount),
		Page:       int32(result.Page),
		PageSize:   int32(result.PageSize),
		TotalPages: int32(result.TotalPages),
	}
}

// Конвертация Product Offer
func convertProductOfferToProto(offer *core.ProductOffer) *searchpb.ProductOffer {
	if offer == nil {
		return nil
	}
	return &searchpb.ProductOffer{
		Id:            int32(offer.ID),
		ShopId:        int32(offer.ShopID),
		ComponentType: offer.ComponentType,
		ComponentId:   int32(offer.ComponentID),
		Price:         float64(offer.Price),
		Available:     offer.Available,
	}
}

func convertProductOfferListToProto(result *core.ListResult[core.ProductOffer]) *searchpb.ProductOfferList {
	if result == nil {
		return &searchpb.ProductOfferList{}
	}

	items := make([]*searchpb.ProductOffer, len(result.Items))
	for i, offer := range result.Items {
		items[i] = convertProductOfferToProto(&offer)
	}

	return &searchpb.ProductOfferList{
		Items:      items,
		TotalCount: int32(result.TotalCount),
		Page:       int32(result.Page),
		PageSize:   int32(result.PageSize),
		TotalPages: int32(result.TotalPages),
	}
}

// Конвертация Assembly
func convertAssemblyToProto(assembly *core.Assembly) *searchpb.Assembly {
	if assembly == nil {
		return nil
	}

	protoAssembly := &searchpb.Assembly{
		Id:               int32(assembly.ID),
		UserId:           int32(assembly.UserID),
		Name:             assembly.Name,
		IsPublic:         assembly.IsPublic,
		TotalPriceCached: assembly.TotalPriceCached,
		CpuId:            optionalInt32(&assembly.CpuID),
		GpuId:            optionalInt32(assembly.GpuID),
		MotherboardId:    optionalInt32(&assembly.MotherboardID),
		PsuId:            optionalInt32(&assembly.PsuID),
		CaseId:           optionalInt32(&assembly.CaseID),
		CoolerId:         optionalInt32(assembly.CoolerID),
	}

	return protoAssembly
}

func convertAssemblyListToProto(result *core.ListResult[core.Assembly]) *searchpb.AssemblyList {
	if result == nil {
		return &searchpb.AssemblyList{}
	}

	items := make([]*searchpb.Assembly, len(result.Items))
	for i, assembly := range result.Items {
		items[i] = convertAssemblyToProto(&assembly)
	}

	return &searchpb.AssemblyList{
		Items:      items,
		TotalCount: int32(result.TotalCount),
		Page:       int32(result.Page),
		PageSize:   int32(result.PageSize),
		TotalPages: int32(result.TotalPages),
	}
}

// Конвертация Validation Result
func convertValidationResultToProto(result *core.ValidationResult) *searchpb.ValidationResult {
	if result == nil {
		return &searchpb.ValidationResult{}
	}

	return &searchpb.ValidationResult{
		IsValid:  result.IsValid,
		Errors:   result.Errors,
		Warnings: result.Warnings,
	}
}

// Конвертация Search Results
func convertSearchResultsToProto(results []core.SearchResult, totalCount, limit, offset int) *searchpb.SearchResponse {
	protoResults := make([]*searchpb.SearchResponse_SearchResultItem, len(results))
	for i, result := range results {
		protoResults[i] = &searchpb.SearchResponse_SearchResultItem{
			ComponentType:  result.ComponentType,
			ComponentId:    int32(result.ComponentID),
			Name:           result.Name,
			Brand:          result.Brand,
			RelevanceScore: float64(result.RelevanceScore),
		}
	}

	return &searchpb.SearchResponse{
		Results:    protoResults,
		TotalCount: int32(totalCount),
		PageSize:   int32(limit),
		Page:       int32(offset/limit + 1),
	}
}

// Конвертация Price List
func convertPriceListToProto(componentType string, componentID int, prices []core.PriceInfo) *searchpb.PriceList {
	protoPrices := make([]*searchpb.PriceList_PriceInfo, len(prices))
	for i, price := range prices {
		protoPrices[i] = &searchpb.PriceList_PriceInfo{
			ShopId:      int32(price.ShopID),
			ShopName:    price.ShopName,
			Price:       price.Price,
			Available:   price.Available,
			Url:         price.URL,
			LastUpdated: price.LastUpdated,
		}
	}

	// Рассчитываем статистику цен
	var minPrice, maxPrice, totalPrice float64
	if len(prices) > 0 {
		minPrice = prices[0].Price
		maxPrice = prices[0].Price
		for _, price := range prices {
			totalPrice += price.Price
			if price.Price < minPrice {
				minPrice = price.Price
			}
			if price.Price > maxPrice {
				maxPrice = price.Price
			}
		}
	}

	return &searchpb.PriceList{
		ComponentType: componentType,
		ComponentId:   int32(componentID),
		Prices:        protoPrices,
		MinPrice:      minPrice,
		MaxPrice:      maxPrice,
		AveragePrice:  totalPrice / float64(len(prices)),
	}
}

// Конвертация Offer List
func convertOfferListToProto(offers []core.Offer) *searchpb.OfferList {
	protoOffers := make([]*searchpb.OfferList_Offer, len(offers))
	for i, offer := range offers {
		protoOffers[i] = &searchpb.OfferList_Offer{
			OfferId:      int32(offer.OfferID),
			ShopId:       int32(offer.ShopID),
			ShopName:     offer.ShopName,
			Price:        offer.Price,
			Available:    offer.Available,
			Url:          offer.URL,
			InStock:      offer.InStock,
			DeliveryTime: offer.DeliveryTime,
		}
	}

	return &searchpb.OfferList{
		Offers:     protoOffers,
		TotalCount: int32(len(offers)),
	}
}

// Вспомогательные функции для опциональных значений
func optionalInt32(val *int) *int32 {
	if val == nil {
		return nil
	}
	v := int32(*val)
	return &v
}

func optionalString(val *string) *string {
	if val == nil {
		return nil
	}
	return val
}

func optionalFloat(val *float64) *float32 {
	if val == nil {
		return nil
	}
	v := float32(*val)
	return &v
}
