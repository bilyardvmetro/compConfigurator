package core

type GPU struct {
	ID                int     `db:"gpu_id" json:"id"`
	Name              string  `db:"name" json:"name"`
	Brand             string  `db:"brand" json:"brand"`
	Interface         string  `db:"interface" json:"interface"`
	PcieVersion       int     `db:"pcie_version" json:"pcie_version"`
	PcieLanesRequired int     `db:"pcie_lanes_required" json:"pcie_lanes_required"`
	TdpWatt           int     `db:"tdp_watt" json:"tdp_watt"`
	PowerConnectors   string  `db:"power_connectors" json:"power_connectors"`
	LengthMm          int     `db:"length_mm" json:"length_mm"`
	WidthSlots        float64 `db:"width_slots" json:"width_slots"`
	HeightMm          int     `db:"height_mm" json:"height_mm"`
}

type CpuCooler struct {
	ID             int    `db:"cooler_id" json:"id"`
	Name           string `db:"name" json:"name"`
	Brand          string `db:"brand" json:"brand"`
	CoolingType    string `db:"cooling_type" json:"cooling_type"`
	TdpLimitWatt   int    `db:"tdp_limit_watt" json:"tdp_limit_watt"`
	HeightMm       int    `db:"height_mm" json:"height_mm"`
	FanCount       int    `db:"fan_count" json:"fan_count"`
	FanControlType string `db:"fan_control_type" json:"fan_control_type"`
}

type Case struct {
	ID                int      `db:"case_id" json:"id"`
	Name              string   `db:"name" json:"name"`
	Brand             string   `db:"brand" json:"brand"`
	PsuFormFactor     string   `db:"psu_form_factor" json:"psu_form_factor"`
	MaxGpuLengthMm    *int     `db:"max_gpu_length_mm" json:"max_gpu_length_mm,omitempty"`
	MaxGpuWidthSlots  *float64 `db:"max_gpu_width_slots" json:"max_gpu_width_slots,omitempty"`
	MaxCoolerHeightMm *int     `db:"max_cooler_height_mm" json:"max_cooler_height_mm,omitempty"`
	MaxPsuLengthMm    *int     `db:"max_psu_length_mm" json:"max_psu_length_mm,omitempty"`
	DriveBays3_5Count *int     `db:"drive_bays_3_5_count" json:"drive_bays_3_5_count,omitempty"`
	DriveBays2_5Count *int     `db:"drive_bays_2_5_count" json:"drive_bays_2_5_count,omitempty"`
}

type PSU struct {
	ID                        int    `db:"psu_id" json:"id"`
	Name                      string `db:"name" json:"name"`
	Brand                     string `db:"brand" json:"brand"`
	PowerWatt                 int    `db:"power_watt" json:"power_watt"`
	FormFactor                string `db:"form_factor" json:"form_factor"`
	EfficiencyRating          string `db:"efficiency_rating" json:"efficiency_rating"`
	PcieConnectors6_8pinCount int    `db:"pcie_connectors_6_8pin_count" json:"pcie_connectors_6_8pin_count"`
	Cpu8pinConnectorsCount    int    `db:"cpu_8pin_connectors_count" json:"cpu_8pin_connectors_count"`
	SataConnectorsCount       int    `db:"sata_connectors_count" json:"sata_connectors_count"`
	MolexConnectorsCount      int    `db:"molex_connectors_count" json:"molex_connectors_count"`
	LengthMm                  *int   `db:"length_mm" json:"length_mm,omitempty"`
}

type CpuSocket struct {
	Code        string `db:"socket_code" json:"code"`
	Description string `db:"description" json:"description"`
}

type MotherboardFormFactor struct {
	Code        string `db:"form_factor_code" json:"code"`
	Description string `db:"description" json:"description"`
	WidthMm     *int   `db:"width_mm" json:"width_mm,omitempty"`
	HeightMm    *int   `db:"height_mm" json:"height_mm,omitempty"`
}

type CPU struct {
	ID                     int    `db:"cpu_id" json:"id"`
	Name                   string `db:"name" json:"name"`
	Brand                  string `db:"brand" json:"brand"`
	SocketCode             string `db:"socket_code" json:"socket_code"`
	Architecture           string `db:"architecture" json:"architecture"`
	CoreCount              int    `db:"core_count" json:"core_count"`
	ThreadCount            int    `db:"thread_count" json:"thread_count"`
	BaseClockMhz           int    `db:"base_clock_mhz" json:"base_clock_mhz"`
	BoostClockMhz          *int   `db:"boost_clock_mhz" json:"boost_clock_mhz,omitempty"`
	TdpWatt                int    `db:"tdp_watt" json:"tdp_watt"`
	HasIntegratedGpu       bool   `db:"has_integrated_gpu" json:"has_integrated_gpu"`
	SupportedRamType       string `db:"supported_ram_type" json:"supported_ram_type"`
	SupportedRamFreqMaxMhz *int   `db:"supported_ram_freq_max_mhz" json:"supported_ram_freq_max_mhz,omitempty"`
	MemoryChannels         int    `db:"memory_channels" json:"memory_channels"`
	PcieVersion            int    `db:"pcie_version" json:"pcie_version"`
	PcieLanesTotal         int    `db:"pcie_lanes_total" json:"pcie_lanes_total"`
}

type Motherboard struct {
	ID                    int    `db:"motherboard_id" json:"id"`
	Name                  string `db:"name" json:"name"`
	Brand                 string `db:"brand" json:"brand"`
	SocketCode            string `db:"socket_code" json:"socket_code"`
	FormFactorCode        string `db:"form_factor_code" json:"form_factor_code"`
	Chipset               string `db:"chipset" json:"chipset"`
	RamType               string `db:"ram_type" json:"ram_type"`
	RamSlots              int    `db:"ram_slots" json:"ram_slots"`
	RamCapacityMaxGb      *int   `db:"ram_capacity_max_gb" json:"ram_capacity_max_gb,omitempty"`
	RamFreqMaxMhz         *int   `db:"ram_freq_max_mhz" json:"ram_freq_max_mhz,omitempty"`
	PcieX16SlotsCount     int    `db:"pcie_x16_slots_count" json:"pcie_x16_slots_count"`
	PcieVersionMax        int    `db:"pcie_version_max" json:"pcie_version_max"`
	M2SlotsCount          int    `db:"m2_slots_count" json:"m2_slots_count"`
	SataPortsCount        int    `db:"sata_ports_count" json:"sata_ports_count"`
	PsuMainConnectorType  string `db:"psu_main_connector_type" json:"psu_main_connector_type"`
	CpuPowerConnectorType string `db:"cpu_power_connector_type" json:"cpu_power_connector_type"`
}

type RAMKit struct {
	ID               int      `db:"ram_kit_id" json:"id"`
	Name             string   `db:"name" json:"name"`
	Brand            string   `db:"brand" json:"brand"`
	RamType          string   `db:"ram_type" json:"ram_type"`
	ModuleCapacityGb int      `db:"module_capacity_gb" json:"module_capacity_gb"`
	ModuleCount      int      `db:"module_count" json:"module_count"`
	TotalCapacityGb  int      `db:"total_capacity_gb" json:"total_capacity_gb"`
	FreqMhz          int      `db:"freq_mhz" json:"freq_mhz"`
	Timings          string   `db:"timings" json:"timings"`
	VoltageV         *float64 `db:"voltage_v" json:"voltage_v,omitempty"`
	FormFactor       string   `db:"form_factor" json:"form_factor"`
}

type StorageDrive struct {
	ID         int    `db:"drive_id" json:"id"`
	Name       string `db:"name" json:"name"`
	Brand      string `db:"brand" json:"brand"`
	DriveType  string `db:"drive_type" json:"drive_type"`
	FormFactor string `db:"form_factor" json:"form_factor"`
	Interface  string `db:"interface" json:"interface"`
	CapacityGb int    `db:"capacity_gb" json:"capacity_gb"`
}

type User struct {
	ID           int     `db:"user_id" json:"id"`
	Email        string  `db:"email" json:"email"`
	PasswordHash string  `db:"password_hash" json:"password_hash"`
	Nickname     *string `db:"nickname" json:"nickname,omitempty"`
	AvatarUrl    *string `db:"avatar_url" json:"avatar_url,omitempty"`
	CreatedAt    string  `db:"created_at" json:"created_at"`
	IsAdmin      bool    `db:"is_admin" json:"is_admin"`
}

type Shop struct {
	ID   int     `db:"shop_id" json:"id"`
	Name string  `db:"name" json:"name"`
	Url  *string `db:"url" json:"url,omitempty"`
}

type ProductOffer struct {
	ID            int     `db:"offer_id" json:"id"`
	ShopID        int     `db:"shop_id" json:"shop_id"`
	ComponentType string  `db:"component_type" json:"component_type"`
	ComponentID   int     `db:"component_id" json:"component_id"`
	Price         float64 `db:"price" json:"price"`
	Available     bool    `db:"available" json:"available"`
}

type Assembly struct {
	ID               int      `db:"assembly_id" json:"id"`
	UserID           int      `db:"user_id" json:"user_id"`
	Name             string   `db:"name" json:"name"`
	CreatedAt        string   `db:"created_at" json:"created_at"`
	UpdatedAt        string   `db:"updated_at" json:"updated_at"`
	IsPublic         bool     `db:"is_public" json:"is_public"`
	TotalPriceCached *float64 `db:"total_price_cached" json:"total_price_cached,omitempty"`
	CpuID            int      `db:"cpu_id" json:"cpu_id"`
	GpuID            *int     `db:"gpu_id" json:"gpu_id,omitempty"`
	MotherboardID    int      `db:"motherboard_id" json:"motherboard_id"`
	PsuID            int      `db:"psu_id" json:"psu_id"`
	CaseID           int      `db:"case_id" json:"case_id"`
	CoolerID         *int     `db:"cooler_id" json:"cooler_id,omitempty"`
}

type CoolerSocket struct {
	CoolerID   int     `db:"cooler_id" json:"cooler_id"`
	SocketCode string  `db:"socket_code" json:"socket_code"`
	Notes      *string `db:"notes" json:"notes,omitempty"`
}

type CaseFormFactorSupport struct {
	CaseID         int    `db:"case_id" json:"case_id"`
	FormFactorCode string `db:"form_factor_code" json:"form_factor_code"`
}

type AssemblyRAMKit struct {
	AssemblyID int `db:"assembly_id" json:"assembly_id"`
	RAMKitID   int `db:"ram_kit_id" json:"ram_kit_id"`
}

type AssemblyDrive struct {
	AssemblyID int     `db:"assembly_id" json:"assembly_id"`
	DriveID    int     `db:"drive_id" json:"drive_id"`
	MountType  *string `db:"mount_type" json:"mount_type,omitempty"`
}

type ListOptions struct {
	Page      int      `json:"page"`       // номер страницы (начиная с 1)
	PageSize  int      `json:"page_size"`  // размер страницы
	SortBy    string   `json:"sort_by"`    // поле для сортировки
	SortOrder string   `json:"sort_order"` // "asc" или "desc"
	Filters   []Filter `json:"filters"`    // фильтры
}

type Filter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"` // "=", ">", "<", "LIKE", "IN"
	Value    interface{} `json:"value"`
}

type ListResult[T any] struct {
	Items      []T `json:"items"`
	TotalCount int `json:"total_count"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalPages int `json:"total_pages"`
}

type ValidationResult struct {
	IsValid  bool     `json:"is_valid"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

type AssemblyComponents struct {
	CpuID         *int         `json:"cpu_id,omitempty"`
	GpuID         *int         `json:"gpu_id,omitempty"`
	MotherboardID *int         `json:"motherboard_id,omitempty"`
	PsuID         *int         `json:"psu_id,omitempty"`
	CaseID        *int         `json:"case_id,omitempty"`
	CoolerID      *int         `json:"cooler_id,omitempty"`
	RAMKitIDs     []int        `json:"ram_kit_ids,omitempty"`
	Drives        []DriveMount `json:"drives,omitempty"`
}

type DriveMount struct {
	DriveID   int    `json:"drive_id"`
	MountType string `json:"mount_type"` // "M.2", "SATA", "3.5", "2.5"
}

type SearchResult struct {
	ComponentType  string  `json:"component_type" db:"component_type"`
	ComponentID    int     `json:"component_id" db:"id"`
	Name           string  `json:"name" db:"name"`
	Brand          string  `json:"brand" db:"brand"`
	RelevanceScore float64 `json:"relevance_score" db:"rank"`
}

type PriceInfo struct {
	ShopID      int     `json:"shop_id" db:"shop_id"`
	ShopName    string  `json:"shop_name" db:"shop_name"`
	Price       float64 `json:"price" db:"price"`
	Available   bool    `json:"available" db:"available"`
	URL         string  `json:"url" db:"url"`
	LastUpdated string  `json:"last_updated" db:"last_updated"`
}

type Offer struct {
	OfferID      int     `json:"offer_id" db:"offer_id"`
	ShopID       int     `json:"shop_id" db:"shop_id"`
	ShopName     string  `json:"shop_name" db:"shop_name"`
	Price        float64 `json:"price" db:"price"`
	Available    bool    `json:"available" db:"available"`
	URL          string  `json:"url" db:"url"`
	InStock      bool    `json:"in_stock" db:"in_stock"`
	DeliveryTime string  `json:"delivery_time" db:"delivery_time"`
}
