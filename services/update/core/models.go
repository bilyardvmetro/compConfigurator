package core

type GPU struct {
	ID                int
	Name              string
	Brand             string
	Interface         string
	PcieVersion       int
	PcieLanesRequired int
	TdpWatt           int
	PowerConnectors   string
	LengthMm          int
	WidthSlots        int
	HeightMm          int
}

type CpuCooler struct {
	ID             int
	Name           string
	Brand          string
	CoolingType    string
	TdpLimitWatt   int
	HeightMm       int
	FanCount       int
	FanControlType string
}

type Case struct {
	ID                int
	Name              string
	Brand             string
	PsuFormFactor     string
	MaxGpuLengthMm    *int
	MaxGpuWidthSlots  *float64
	MaxCoolerHeightMm *int
	MaxPsuLengthMm    *int
	DriveBays3_5Count *int
	DriveBays2_5Count *int
}

type PSU struct {
	ID                        int
	Name                      string
	Brand                     string
	PowerWatt                 int
	FormFactor                string
	EfficiencyRating          string
	PcieConnectors6_8pinCount int
	Cpu8pinConnectorsCount    int
	SataConnectorsCount       int
	MolexConnectorsCount      int
	LengthMm                  *int
}

type CpuSocket struct {
	Code        string
	Description string
}

type MotherboardFormFactor struct {
	Code        string
	Description string
	WidthMm     *int
	HeightMm    *int
}

type CPU struct {
	ID                     int
	Name                   string
	Brand                  string
	SocketCode             string
	Architecture           string
	CoreCount              int
	ThreadCount            int
	BaseClockMhz           int
	BoostClockMhz          *int
	TdpWatt                int
	HasIntegratedGpu       bool
	SupportedRamType       string
	SupportedRamFreqMaxMhz *int
	MemoryChannels         int
	PcieVersion            int
	PcieLanesTotal         int
}

type Motherboard struct {
	ID                    int
	Name                  string
	Brand                 string
	SocketCode            string
	FormFactorCode        string
	Chipset               string
	RamType               string
	RamSlots              int
	RamCapacityMaxGb      *int
	RamFreqMaxMhz         *int
	PcieX16SlotsCount     int
	PcieVersionMax        int
	M2SlotsCount          int
	SataPortsCount        int
	PsuMainConnectorType  string
	CpuPowerConnectorType string
}

type RAMKit struct {
	ID               int
	Name             string
	Brand            string
	RamType          string
	ModuleCapacityGb int
	ModuleCount      int
	TotalCapacityGb  int
	FreqMhz          int
	Timings          string
	VoltageV         *float64
	FormFactor       string
}

type StorageDrive struct {
	ID         int
	Name       string
	Brand      string
	DriveType  string
	FormFactor string
	Interface  string
	CapacityGb int
}

type User struct {
	ID           int
	Email        string
	PasswordHash string
	Nickname     *string
	AvatarUrl    *string
	CreatedAt    string
	IsAdmin      bool
}

type Shop struct {
	ID   int
	Name string
	Url  *string
}

type ProductOffer struct {
	ID            int
	ShopID        int
	ComponentType string
	ComponentID   int
	Price         float64
	Available     bool
}

type Assembly struct {
	ID               int
	UserID           int
	Name             string
	CreatedAt        string
	UpdatedAt        string
	IsPublic         bool
	TotalPriceCached *float64
	CpuID            int
	GpuID            *int
	MotherboardID    int
	PsuID            int
	CaseID           int
	CoolerID         *int
}

type CoolerSocket struct {
	CoolerID   int
	SocketCode string
	Notes      *string
}

type CaseFormFactorSupport struct {
	CaseID         int
	FormFactorCode string
}

type AssemblyRAMKit struct {
	AssemblyID int
	RAMKitID   int
}

type AssemblyDrive struct {
	AssemblyID int
	DriveID    int
	MountType  *string
}
