package db

import (
	"context"
	"log/slog"

	"comp_config.com/services/update/core"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	log  *slog.Logger
	conn *sqlx.DB
}

func New(log *slog.Logger, address string) (*DB, error) {
	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}

	return &DB{
		log:  log,
		conn: db,
	}, nil
}

// GPU methods
func (db *DB) AddGpu(ctx context.Context, element core.GPU) error {
	query := `
		INSERT INTO gpus (name, brand, interface, pcie_version, pcie_lanes_required, tdp_watt, power_connectors, length_mm, width_slots, height_mm)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := db.conn.Exec(query, element.Name, element.Brand, element.Interface,
		element.PcieVersion, element.PcieLanesRequired, element.TdpWatt,
		element.PowerConnectors, element.LengthMm, element.WidthSlots, element.HeightMm)
	if err != nil {
		db.log.Error("Failed to add gpu to db", "name", element.Name, "error", err)
		return err
	}
	db.log.Info("Gpu added to db", "name", element.Name)
	return nil
}

func (db *DB) UpdateGpu(ctx context.Context, element core.GPU) error {
	query := `
		UPDATE gpus 
		SET name = $1, brand = $2, interface = $3, pcie_version = $4, 
		    pcie_lanes_required = $5, tdp_watt = $6, power_connectors = $7,
		    length_mm = $8, width_slots = $9, height_mm = $10
		WHERE gpu_id = $11
	`
	result, err := db.conn.Exec(query,
		element.Name, element.Brand, element.Interface,
		element.PcieVersion, element.PcieLanesRequired, element.TdpWatt,
		element.PowerConnectors, element.LengthMm, element.WidthSlots, element.HeightMm,
		element.ID)
	if err != nil {
		db.log.Error("Failed to update gpu", "id", element.ID, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("GPU not found for update", "id", element.ID)
		return core.ErrNotFound
	}

	db.log.Info("Gpu updated", "id", element.ID, "name", element.Name)
	return nil
}

func (db *DB) DeleteGpu(ctx context.Context, id int) error {
	query := `DELETE FROM gpus WHERE gpu_id = $1`
	result, err := db.conn.Exec(query, id)
	if err != nil {
		db.log.Error("Failed to delete gpu", "id", id, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("GPU not found for deletion", "id", id)
		return core.ErrNotFound
	}

	db.log.Info("Gpu deleted", "id", id)
	return nil
}

// CPU Cooler methods
func (db *DB) AddCpuCooler(ctx context.Context, element core.CpuCooler) error {
	query := `
		INSERT INTO cpu_coolers (name, brand, cooling_type, tdp_limit_watt, height_mm, fan_count, fan_control_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := db.conn.Exec(query, element.Name, element.Brand, element.CoolingType,
		element.TdpLimitWatt, element.HeightMm, element.FanCount, element.FanControlType)
	if err != nil {
		db.log.Error("Failed to add cpu cooler to db", "name", element.Name, "error", err)
		return err
	}
	db.log.Info("Cpu cooler added to db", "name", element.Name)
	return nil
}

func (db *DB) UpdateCpuCooler(ctx context.Context, element core.CpuCooler) error {
	query := `
		UPDATE cpu_coolers 
		SET name = $1, brand = $2, cooling_type = $3, tdp_limit_watt = $4,
		    height_mm = $5, fan_count = $6, fan_control_type = $7
		WHERE cooler_id = $8
	`
	result, err := db.conn.Exec(query,
		element.Name, element.Brand, element.CoolingType,
		element.TdpLimitWatt, element.HeightMm, element.FanCount, element.FanControlType,
		element.ID)
	if err != nil {
		db.log.Error("Failed to update cpu cooler", "id", element.ID, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("CPU cooler not found for update", "id", element.ID)
		return core.ErrNotFound
	}

	db.log.Info("CPU cooler updated", "id", element.ID, "name", element.Name)
	return nil
}

func (db *DB) DeleteCpuCooler(ctx context.Context, id int) error {
	query := `DELETE FROM cpu_coolers WHERE cooler_id = $1`
	result, err := db.conn.Exec(query, id)
	if err != nil {
		db.log.Error("Failed to delete cpu cooler", "id", id, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("CPU cooler not found for deletion", "id", id)
		return core.ErrNotFound
	}

	db.log.Info("CPU cooler deleted", "id", id)
	return nil
}

// Case methods
func (db *DB) AddCase(ctx context.Context, element core.Case) error {
	query := `
        INSERT INTO cases (name, brand, psu_form_factor, max_gpu_length_mm, 
                          max_gpu_width_slots, max_cooler_height_mm, max_psu_length_mm,
                          drive_bays_3_5_count, drive_bays_2_5_count)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	_, err := db.conn.Exec(query,
		element.Name, element.Brand, element.PsuFormFactor,
		element.MaxGpuLengthMm, element.MaxGpuWidthSlots, element.MaxCoolerHeightMm,
		element.MaxPsuLengthMm, element.DriveBays3_5Count, element.DriveBays2_5Count)
	if err != nil {
		db.log.Error("Failed to add case to db", "name", element.Name, "error", err)
		return err
	}
	db.log.Info("Case added to db", "name", element.Name)
	return nil
}

func (db *DB) UpdateCase(ctx context.Context, element core.Case) error {
	query := `
        UPDATE cases 
        SET name = $1, brand = $2, psu_form_factor = $3, max_gpu_length_mm = $4,
            max_gpu_width_slots = $5, max_cooler_height_mm = $6, max_psu_length_mm = $7,
            drive_bays_3_5_count = $8, drive_bays_2_5_count = $9
        WHERE case_id = $10
    `
	result, err := db.conn.Exec(query,
		element.Name, element.Brand, element.PsuFormFactor,
		element.MaxGpuLengthMm, element.MaxGpuWidthSlots, element.MaxCoolerHeightMm,
		element.MaxPsuLengthMm, element.DriveBays3_5Count, element.DriveBays2_5Count,
		element.ID)
	if err != nil {
		db.log.Error("Failed to update case", "id", element.ID, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Case not found for update", "id", element.ID)
		return core.ErrNotFound
	}

	db.log.Info("Case updated", "id", element.ID, "name", element.Name)
	return nil
}

func (db *DB) DeleteCase(ctx context.Context, id int) error {
	query := `DELETE FROM cases WHERE case_id = $1`
	result, err := db.conn.Exec(query, id)
	if err != nil {
		db.log.Error("Failed to delete case", "id", id, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Case not found for deletion", "id", id)
		return core.ErrNotFound
	}

	db.log.Info("Case deleted", "id", id)
	return nil
}

// PSU methods
func (db *DB) AddPsu(ctx context.Context, element core.PSU) error {
	query := `
        INSERT INTO psus (name, brand, power_watt, form_factor, efficiency_rating,
                         pcie_connectors_6_8pin_count, cpu_8pin_connectors_count,
                         sata_connectors_count, molex_connectors_count, length_mm)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `
	_, err := db.conn.Exec(query,
		element.Name, element.Brand, element.PowerWatt, element.FormFactor,
		element.EfficiencyRating, element.PcieConnectors6_8pinCount,
		element.Cpu8pinConnectorsCount, element.SataConnectorsCount,
		element.MolexConnectorsCount, element.LengthMm)
	if err != nil {
		db.log.Error("Failed to add PSU to db", "name", element.Name, "error", err)
		return err
	}
	db.log.Info("PSU added to db", "name", element.Name)
	return nil
}

func (db *DB) UpdatePsu(ctx context.Context, element core.PSU) error {
	query := `
        UPDATE psus 
        SET name = $1, brand = $2, power_watt = $3, form_factor = $4,
            efficiency_rating = $5, pcie_connectors_6_8pin_count = $6,
            cpu_8pin_connectors_count = $7, sata_connectors_count = $8,
            molex_connectors_count = $9, length_mm = $10
        WHERE psu_id = $11
    `
	result, err := db.conn.Exec(query,
		element.Name, element.Brand, element.PowerWatt, element.FormFactor,
		element.EfficiencyRating, element.PcieConnectors6_8pinCount,
		element.Cpu8pinConnectorsCount, element.SataConnectorsCount,
		element.MolexConnectorsCount, element.LengthMm,
		element.ID)
	if err != nil {
		db.log.Error("Failed to update PSU", "id", element.ID, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("PSU not found for update", "id", element.ID)
		return core.ErrNotFound
	}

	db.log.Info("PSU updated", "id", element.ID, "name", element.Name)
	return nil
}

func (db *DB) DeletePsu(ctx context.Context, id int) error {
	query := `DELETE FROM psus WHERE psu_id = $1`
	result, err := db.conn.Exec(query, id)
	if err != nil {
		db.log.Error("Failed to delete PSU", "id", id, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("PSU not found for deletion", "id", id)
		return core.ErrNotFound
	}

	db.log.Info("PSU deleted", "id", id)
	return nil
}

// CPU Socket methods
func (db *DB) AddCpuSocket(ctx context.Context, element core.CpuSocket) error {
	query := `
        INSERT INTO cpu_sockets (socket_code, description)
        VALUES ($1, $2)
    `
	_, err := db.conn.Exec(query, element.Code, element.Description)
	if err != nil {
		db.log.Error("Failed to add CPU socket to db", "code", element.Code, "error", err)
		return err
	}
	db.log.Info("CPU socket added to db", "code", element.Code)
	return nil
}

func (db *DB) UpdateCpuSocket(ctx context.Context, element core.CpuSocket) error {
	query := `
        UPDATE cpu_sockets 
        SET description = $1
        WHERE socket_code = $2
    `
	result, err := db.conn.Exec(query, element.Description, element.Code)
	if err != nil {
		db.log.Error("Failed to update CPU socket", "code", element.Code, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("CPU socket not found for update", "code", element.Code)
		return core.ErrNotFound
	}

	db.log.Info("CPU socket updated", "code", element.Code)
	return nil
}

func (db *DB) DeleteCpuSocket(ctx context.Context, code string) error {
	query := `DELETE FROM cpu_sockets WHERE socket_code = $1`
	result, err := db.conn.Exec(query, code)
	if err != nil {
		db.log.Error("Failed to delete CPU socket", "code", code, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("CPU socket not found for deletion", "code", code)
		return core.ErrNotFound
	}

	db.log.Info("CPU socket deleted", "code", code)
	return nil
}

// Motherboard Form Factor methods
func (db *DB) AddMotherboardFormFactor(ctx context.Context, element core.MotherboardFormFactor) error {
	query := `
        INSERT INTO motherboard_form_factors (form_factor_code, description, width_mm, height_mm)
        VALUES ($1, $2, $3, $4)
    `
	_, err := db.conn.Exec(query, element.Code, element.Description, element.WidthMm, element.HeightMm)
	if err != nil {
		db.log.Error("Failed to add motherboard form factor to db", "code", element.Code, "error", err)
		return err
	}
	db.log.Info("Motherboard form factor added to db", "code", element.Code)
	return nil
}

func (db *DB) UpdateMotherboardFormFactor(ctx context.Context, element core.MotherboardFormFactor) error {
	query := `
        UPDATE motherboard_form_factors 
        SET description = $1, width_mm = $2, height_mm = $3
        WHERE form_factor_code = $4
    `
	result, err := db.conn.Exec(query, element.Description, element.WidthMm, element.HeightMm, element.Code)
	if err != nil {
		db.log.Error("Failed to update motherboard form factor", "code", element.Code, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Motherboard form factor not found for update", "code", element.Code)
		return core.ErrNotFound
	}

	db.log.Info("Motherboard form factor updated", "code", element.Code)
	return nil
}

func (db *DB) DeleteMotherboardFormFactor(ctx context.Context, code string) error {
	query := `DELETE FROM motherboard_form_factors WHERE form_factor_code = $1`
	result, err := db.conn.Exec(query, code)
	if err != nil {
		db.log.Error("Failed to delete motherboard form factor", "code", code, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Motherboard form factor not found for deletion", "code", code)
		return core.ErrNotFound
	}

	db.log.Info("Motherboard form factor deleted", "code", code)
	return nil
}

// CPU methods
func (db *DB) AddCpu(ctx context.Context, element core.CPU) error {
	query := `
        INSERT INTO cpus (name, brand, socket_code, architecture, core_count,
                         thread_count, base_clock_mhz, boost_clock_mhz, tdp_watt,
                         has_integrated_gpu, supported_ram_type, supported_ram_freq_max_mhz,
                         memory_channels, pcie_version, pcie_lanes_total)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
    `
	_, err := db.conn.Exec(query,
		element.Name, element.Brand, element.SocketCode, element.Architecture,
		element.CoreCount, element.ThreadCount, element.BaseClockMhz, element.BoostClockMhz,
		element.TdpWatt, element.HasIntegratedGpu, element.SupportedRamType,
		element.SupportedRamFreqMaxMhz, element.MemoryChannels, element.PcieVersion,
		element.PcieLanesTotal)
	if err != nil {
		db.log.Error("Failed to add CPU to db", "name", element.Name, "error", err)
		return err
	}
	db.log.Info("CPU added to db", "name", element.Name)
	return nil
}

func (db *DB) UpdateCpu(ctx context.Context, element core.CPU) error {
	query := `
        UPDATE cpus 
        SET name = $1, brand = $2, socket_code = $3, architecture = $4,
            core_count = $5, thread_count = $6, base_clock_mhz = $7,
            boost_clock_mhz = $8, tdp_watt = $9, has_integrated_gpu = $10,
            supported_ram_type = $11, supported_ram_freq_max_mhz = $12,
            memory_channels = $13, pcie_version = $14, pcie_lanes_total = $15
        WHERE cpu_id = $16
    `
	result, err := db.conn.Exec(query,
		element.Name, element.Brand, element.SocketCode, element.Architecture,
		element.CoreCount, element.ThreadCount, element.BaseClockMhz, element.BoostClockMhz,
		element.TdpWatt, element.HasIntegratedGpu, element.SupportedRamType,
		element.SupportedRamFreqMaxMhz, element.MemoryChannels, element.PcieVersion,
		element.PcieLanesTotal,
		element.ID)
	if err != nil {
		db.log.Error("Failed to update CPU", "id", element.ID, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("CPU not found for update", "id", element.ID)
		return core.ErrNotFound
	}

	db.log.Info("CPU updated", "id", element.ID, "name", element.Name)
	return nil
}

func (db *DB) DeleteCpu(ctx context.Context, id int) error {
	query := `DELETE FROM cpus WHERE cpu_id = $1`
	result, err := db.conn.Exec(query, id)
	if err != nil {
		db.log.Error("Failed to delete CPU", "id", id, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("CPU not found for deletion", "id", id)
		return core.ErrNotFound
	}

	db.log.Info("CPU deleted", "id", id)
	return nil
}

// Motherboard methods
func (db *DB) AddMotherboard(ctx context.Context, element core.Motherboard) error {
	query := `
        INSERT INTO motherboards (name, brand, socket_code, form_factor_code,
                                 chipset, ram_type, ram_slots, ram_capacity_max_gb,
                                 ram_freq_max_mhz, pcie_x16_slots_count, pcie_version_max,
                                 m2_slots_count, sata_ports_count, psu_main_connector_type,
                                 cpu_power_connector_type)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
    `
	_, err := db.conn.Exec(query,
		element.Name, element.Brand, element.SocketCode, element.FormFactorCode,
		element.Chipset, element.RamType, element.RamSlots, element.RamCapacityMaxGb,
		element.RamFreqMaxMhz, element.PcieX16SlotsCount, element.PcieVersionMax,
		element.M2SlotsCount, element.SataPortsCount, element.PsuMainConnectorType,
		element.CpuPowerConnectorType)
	if err != nil {
		db.log.Error("Failed to add motherboard to db", "name", element.Name, "error", err)
		return err
	}
	db.log.Info("Motherboard added to db", "name", element.Name)
	return nil
}

func (db *DB) UpdateMotherboard(ctx context.Context, element core.Motherboard) error {
	query := `
        UPDATE motherboards 
        SET name = $1, brand = $2, socket_code = $3, form_factor_code = $4,
            chipset = $5, ram_type = $6, ram_slots = $7, ram_capacity_max_gb = $8,
            ram_freq_max_mhz = $9, pcie_x16_slots_count = $10, pcie_version_max = $11,
            m2_slots_count = $12, sata_ports_count = $13, psu_main_connector_type = $14,
            cpu_power_connector_type = $15
        WHERE motherboard_id = $16
    `
	result, err := db.conn.Exec(query,
		element.Name, element.Brand, element.SocketCode, element.FormFactorCode,
		element.Chipset, element.RamType, element.RamSlots, element.RamCapacityMaxGb,
		element.RamFreqMaxMhz, element.PcieX16SlotsCount, element.PcieVersionMax,
		element.M2SlotsCount, element.SataPortsCount, element.PsuMainConnectorType,
		element.CpuPowerConnectorType,
		element.ID)
	if err != nil {
		db.log.Error("Failed to update motherboard", "id", element.ID, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Motherboard not found for update", "id", element.ID)
		return core.ErrNotFound
	}

	db.log.Info("Motherboard updated", "id", element.ID, "name", element.Name)
	return nil
}

func (db *DB) DeleteMotherboard(ctx context.Context, id int) error {
	query := `DELETE FROM motherboards WHERE motherboard_id = $1`
	result, err := db.conn.Exec(query, id)
	if err != nil {
		db.log.Error("Failed to delete motherboard", "id", id, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Motherboard not found for deletion", "id", id)
		return core.ErrNotFound
	}

	db.log.Info("Motherboard deleted", "id", id)
	return nil
}

// RAM Kit methods
func (db *DB) AddRAMKit(ctx context.Context, element core.RAMKit) error {
	query := `
        INSERT INTO ram_kits (name, brand, ram_type, module_capacity_gb, module_count,
                             total_capacity_gb, freq_mhz, timings, voltage_v, form_factor)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `
	_, err := db.conn.Exec(query,
		element.Name, element.Brand, element.RamType, element.ModuleCapacityGb,
		element.ModuleCount, element.TotalCapacityGb, element.FreqMhz,
		element.Timings, element.VoltageV, element.FormFactor)
	if err != nil {
		db.log.Error("Failed to add RAM kit to db", "name", element.Name, "error", err)
		return err
	}
	db.log.Info("RAM kit added to db", "name", element.Name)
	return nil
}

func (db *DB) UpdateRAMKit(ctx context.Context, element core.RAMKit) error {
	query := `
        UPDATE ram_kits 
        SET name = $1, brand = $2, ram_type = $3, module_capacity_gb = $4,
            module_count = $5, total_capacity_gb = $6, freq_mhz = $7,
            timings = $8, voltage_v = $9, form_factor = $10
        WHERE ram_kit_id = $11
    `
	result, err := db.conn.Exec(query,
		element.Name, element.Brand, element.RamType, element.ModuleCapacityGb,
		element.ModuleCount, element.TotalCapacityGb, element.FreqMhz,
		element.Timings, element.VoltageV, element.FormFactor,
		element.ID)
	if err != nil {
		db.log.Error("Failed to update RAM kit", "id", element.ID, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("RAM kit not found for update", "id", element.ID)
		return core.ErrNotFound
	}

	db.log.Info("RAM kit updated", "id", element.ID, "name", element.Name)
	return nil
}

func (db *DB) DeleteRAMKit(ctx context.Context, id int) error {
	query := `DELETE FROM ram_kits WHERE ram_kit_id = $1`
	result, err := db.conn.Exec(query, id)
	if err != nil {
		db.log.Error("Failed to delete RAM kit", "id", id, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("RAM kit not found for deletion", "id", id)
		return core.ErrNotFound
	}

	db.log.Info("RAM kit deleted", "id", id)
	return nil
}

// Storage Drive methods
func (db *DB) AddStorageDrive(ctx context.Context, element core.StorageDrive) error {
	query := `
        INSERT INTO storage_drives (name, brand, drive_type, form_factor, interface, capacity_gb)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := db.conn.Exec(query,
		element.Name, element.Brand, element.DriveType, element.FormFactor,
		element.Interface, element.CapacityGb)
	if err != nil {
		db.log.Error("Failed to add storage drive to db", "name", element.Name, "error", err)
		return err
	}
	db.log.Info("Storage drive added to db", "name", element.Name)
	return nil
}

func (db *DB) UpdateStorageDrive(ctx context.Context, element core.StorageDrive) error {
	query := `
        UPDATE storage_drives 
        SET name = $1, brand = $2, drive_type = $3, form_factor = $4,
            interface = $5, capacity_gb = $6
        WHERE drive_id = $7
    `
	result, err := db.conn.Exec(query,
		element.Name, element.Brand, element.DriveType, element.FormFactor,
		element.Interface, element.CapacityGb,
		element.ID)
	if err != nil {
		db.log.Error("Failed to update storage drive", "id", element.ID, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Storage drive not found for update", "id", element.ID)
		return core.ErrNotFound
	}

	db.log.Info("Storage drive updated", "id", element.ID, "name", element.Name)
	return nil
}

func (db *DB) DeleteStorageDrive(ctx context.Context, id int) error {
	query := `DELETE FROM storage_drives WHERE drive_id = $1`
	result, err := db.conn.Exec(query, id)
	if err != nil {
		db.log.Error("Failed to delete storage drive", "id", id, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Storage drive not found for deletion", "id", id)
		return core.ErrNotFound
	}

	db.log.Info("Storage drive deleted", "id", id)
	return nil
}

// User methods
func (db *DB) AddUser(ctx context.Context, element core.User) error {
	query := `
        INSERT INTO users (email, password_hash, nickname, avatar_url, is_admin)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := db.conn.Exec(query,
		element.Email, element.PasswordHash, element.Nickname,
		element.AvatarUrl, element.IsAdmin)
	if err != nil {
		db.log.Error("Failed to add user to db", "email", element.Email, "error", err)
		return err
	}
	db.log.Info("User added to db", "email", element.Email)
	return nil
}

func (db *DB) UpdateUser(ctx context.Context, element core.User) error {
	query := `
        UPDATE users 
        SET email = $1, password_hash = $2, nickname = $3,
            avatar_url = $4, is_admin = $5
        WHERE user_id = $6
    `
	result, err := db.conn.Exec(query,
		element.Email, element.PasswordHash, element.Nickname,
		element.AvatarUrl, element.IsAdmin,
		element.ID)
	if err != nil {
		db.log.Error("Failed to update user", "id", element.ID, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("User not found for update", "id", element.ID)
		return core.ErrNotFound
	}

	db.log.Info("User updated", "id", element.ID, "email", element.Email)
	return nil
}

func (db *DB) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE user_id = $1`
	result, err := db.conn.Exec(query, id)
	if err != nil {
		db.log.Error("Failed to delete user", "id", id, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("User not found for deletion", "id", id)
		return core.ErrNotFound
	}

	db.log.Info("User deleted", "id", id)
	return nil
}

// Shop methods
func (db *DB) AddShop(ctx context.Context, element core.Shop) error {
	query := `
        INSERT INTO shops (name, url)
        VALUES ($1, $2)
    `
	_, err := db.conn.Exec(query, element.Name, element.Url)
	if err != nil {
		db.log.Error("Failed to add shop to db", "name", element.Name, "error", err)
		return err
	}
	db.log.Info("Shop added to db", "name", element.Name)
	return nil
}

func (db *DB) UpdateShop(ctx context.Context, element core.Shop) error {
	query := `
        UPDATE shops 
        SET name = $1, url = $2
        WHERE shop_id = $3
    `
	result, err := db.conn.Exec(query,
		element.Name, element.Url,
		element.ID)
	if err != nil {
		db.log.Error("Failed to update shop", "id", element.ID, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Shop not found for update", "id", element.ID)
		return core.ErrNotFound
	}

	db.log.Info("Shop updated", "id", element.ID, "name", element.Name)
	return nil
}

func (db *DB) DeleteShop(ctx context.Context, id int) error {
	query := `DELETE FROM shops WHERE shop_id = $1`
	result, err := db.conn.Exec(query, id)
	if err != nil {
		db.log.Error("Failed to delete shop", "id", id, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Shop not found for deletion", "id", id)
		return core.ErrNotFound
	}

	db.log.Info("Shop deleted", "id", id)
	return nil
}

// Product Offer methods
func (db *DB) AddProductOffer(ctx context.Context, element core.ProductOffer) error {
	query := `
        INSERT INTO product_offers (shop_id, component_type, component_id, price, available)
        VALUES ($1, $2, $3, $4, $5)
    `
	_, err := db.conn.Exec(query,
		element.ShopID, element.ComponentType, element.ComponentID,
		element.Price, element.Available)
	if err != nil {
		db.log.Error("Failed to add product offer to db", "shop_id", element.ShopID, "error", err)
		return err
	}
	db.log.Info("Product offer added to db", "shop_id", element.ShopID, "component_type", element.ComponentType)
	return nil
}

func (db *DB) UpdateProductOffer(ctx context.Context, element core.ProductOffer) error {
	query := `
        UPDATE product_offers 
        SET shop_id = $1, component_type = $2, component_id = $3,
            price = $4, available = $5
        WHERE offer_id = $6
    `
	result, err := db.conn.Exec(query,
		element.ShopID, element.ComponentType, element.ComponentID,
		element.Price, element.Available,
		element.ID)
	if err != nil {
		db.log.Error("Failed to update product offer", "id", element.ID, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Product offer not found for update", "id", element.ID)
		return core.ErrNotFound
	}

	db.log.Info("Product offer updated", "id", element.ID)
	return nil
}

func (db *DB) DeleteProductOffer(ctx context.Context, id int) error {
	query := `DELETE FROM product_offers WHERE offer_id = $1`
	result, err := db.conn.Exec(query, id)
	if err != nil {
		db.log.Error("Failed to delete product offer", "id", id, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Product offer not found for deletion", "id", id)
		return core.ErrNotFound
	}

	db.log.Info("Product offer deleted", "id", id)
	return nil
}

// Assembly methods
func (db *DB) AddAssembly(ctx context.Context, element core.Assembly) error {
	query := `
        INSERT INTO assemblies (user_id, name, is_public, total_price_cached,
                               cpu_id, gpu_id, motherboard_id, psu_id, case_id, cooler_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `
	_, err := db.conn.Exec(query,
		element.UserID, element.Name, element.IsPublic, element.TotalPriceCached,
		element.CpuID, element.GpuID, element.MotherboardID,
		element.PsuID, element.CaseID, element.CoolerID)
	if err != nil {
		db.log.Error("Failed to add assembly to db", "name", element.Name, "error", err)
		return err
	}
	db.log.Info("Assembly added to db", "name", element.Name)
	return nil
}

func (db *DB) UpdateAssembly(ctx context.Context, element core.Assembly) error {
	query := `
        UPDATE assemblies 
        SET user_id = $1, name = $2, is_public = $3, total_price_cached = $4,
            cpu_id = $5, gpu_id = $6, motherboard_id = $7,
            psu_id = $8, case_id = $9, cooler_id = $10,
            updated_at = NOW()
        WHERE assembly_id = $11
    `
	result, err := db.conn.Exec(query,
		element.UserID, element.Name, element.IsPublic, element.TotalPriceCached,
		element.CpuID, element.GpuID, element.MotherboardID,
		element.PsuID, element.CaseID, element.CoolerID,
		element.ID)
	if err != nil {
		db.log.Error("Failed to update assembly", "id", element.ID, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Assembly not found for update", "id", element.ID)
		return core.ErrNotFound
	}

	db.log.Info("Assembly updated", "id", element.ID, "name", element.Name)
	return nil
}

func (db *DB) DeleteAssembly(ctx context.Context, id int) error {
	query := `DELETE FROM assemblies WHERE assembly_id = $1`
	result, err := db.conn.Exec(query, id)
	if err != nil {
		db.log.Error("Failed to delete assembly", "id", id, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Assembly not found for deletion", "id", id)
		return core.ErrNotFound
	}

	db.log.Info("Assembly deleted", "id", id)
	return nil
}

// Cooler Socket methods (many-to-many)
func (db *DB) AddCoolerSocket(ctx context.Context, element core.CoolerSocket) error {
	query := `
        INSERT INTO cooler_sockets (cooler_id, socket_code, notes)
        VALUES ($1, $2, $3)
        ON CONFLICT (cooler_id, socket_code) DO UPDATE
        SET notes = EXCLUDED.notes
    `
	_, err := db.conn.Exec(query, element.CoolerID, element.SocketCode, element.Notes)
	if err != nil {
		db.log.Error("Failed to add cooler socket", "cooler_id", element.CoolerID, "socket_code", element.SocketCode, "error", err)
		return err
	}
	db.log.Info("Cooler socket added", "cooler_id", element.CoolerID, "socket_code", element.SocketCode)
	return nil
}

func (db *DB) DeleteCoolerSocket(ctx context.Context, coolerID int, socketCode string) error {
	query := `DELETE FROM cooler_sockets WHERE cooler_id = $1 AND socket_code = $2`
	result, err := db.conn.Exec(query, coolerID, socketCode)
	if err != nil {
		db.log.Error("Failed to delete cooler socket", "cooler_id", coolerID, "socket_code", socketCode, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Cooler socket not found for deletion", "cooler_id", coolerID, "socket_code", socketCode)
		return core.ErrNotFound
	}

	db.log.Info("Cooler socket deleted", "cooler_id", coolerID, "socket_code", socketCode)
	return nil
}

func (db *DB) DeleteCoolerSocketsByCooler(ctx context.Context, coolerID int) error {
	query := `DELETE FROM cooler_sockets WHERE cooler_id = $1`
	_, err := db.conn.Exec(query, coolerID)
	if err != nil {
		db.log.Error("Failed to delete cooler sockets by cooler", "cooler_id", coolerID, "error", err)
		return err
	}
	db.log.Info("Cooler sockets deleted by cooler", "cooler_id", coolerID)
	return nil
}

// Case Form Factor Support methods (many-to-many)
func (db *DB) AddCaseFormFactorSupport(ctx context.Context, element core.CaseFormFactorSupport) error {
	query := `
        INSERT INTO case_form_factor_support (case_id, form_factor_code)
        VALUES ($1, $2)
        ON CONFLICT (case_id, form_factor_code) DO NOTHING
    `
	_, err := db.conn.Exec(query, element.CaseID, element.FormFactorCode)
	if err != nil {
		db.log.Error("Failed to add case form factor support", "case_id", element.CaseID, "form_factor_code", element.FormFactorCode, "error", err)
		return err
	}
	db.log.Info("Case form factor support added", "case_id", element.CaseID, "form_factor_code", element.FormFactorCode)
	return nil
}

func (db *DB) DeleteCaseFormFactorSupport(ctx context.Context, caseID int, formFactorCode string) error {
	query := `DELETE FROM case_form_factor_support WHERE case_id = $1 AND form_factor_code = $2`
	result, err := db.conn.Exec(query, caseID, formFactorCode)
	if err != nil {
		db.log.Error("Failed to delete case form factor support", "case_id", caseID, "form_factor_code", formFactorCode, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Case form factor support not found for deletion", "case_id", caseID, "form_factor_code", formFactorCode)
		return core.ErrNotFound
	}

	db.log.Info("Case form factor support deleted", "case_id", caseID, "form_factor_code", formFactorCode)
	return nil
}

func (db *DB) DeleteCaseFormFactorSupportsByCase(ctx context.Context, caseID int) error {
	query := `DELETE FROM case_form_factor_support WHERE case_id = $1`
	_, err := db.conn.Exec(query, caseID)
	if err != nil {
		db.log.Error("Failed to delete case form factor supports by case", "case_id", caseID, "error", err)
		return err
	}
	db.log.Info("Case form factor supports deleted by case", "case_id", caseID)
	return nil
}

// Assembly RAM Kit methods (many-to-many)
func (db *DB) AddAssemblyRAMKit(ctx context.Context, element core.AssemblyRAMKit) error {
	query := `
        INSERT INTO assembly_ram_kits (assembly_id, ram_kit_id)
        VALUES ($1, $2)
        ON CONFLICT (assembly_id, ram_kit_id) DO NOTHING
    `
	_, err := db.conn.Exec(query, element.AssemblyID, element.RAMKitID)
	if err != nil {
		db.log.Error("Failed to add assembly RAM kit", "assembly_id", element.AssemblyID, "ram_kit_id", element.RAMKitID, "error", err)
		return err
	}
	db.log.Info("Assembly RAM kit added", "assembly_id", element.AssemblyID, "ram_kit_id", element.RAMKitID)
	return nil
}

func (db *DB) DeleteAssemblyRAMKit(ctx context.Context, assemblyID int, ramKitID int) error {
	query := `DELETE FROM assembly_ram_kits WHERE assembly_id = $1 AND ram_kit_id = $2`
	result, err := db.conn.Exec(query, assemblyID, ramKitID)
	if err != nil {
		db.log.Error("Failed to delete assembly RAM kit", "assembly_id", assemblyID, "ram_kit_id", ramKitID, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Assembly RAM kit not found for deletion", "assembly_id", assemblyID, "ram_kit_id", ramKitID)
		return core.ErrNotFound
	}

	db.log.Info("Assembly RAM kit deleted", "assembly_id", assemblyID, "ram_kit_id", ramKitID)
	return nil
}

func (db *DB) DeleteAssemblyRAMKitsByAssembly(ctx context.Context, assemblyID int) error {
	query := `DELETE FROM assembly_ram_kits WHERE assembly_id = $1`
	_, err := db.conn.Exec(query, assemblyID)
	if err != nil {
		db.log.Error("Failed to delete assembly RAM kits by assembly", "assembly_id", assemblyID, "error", err)
		return err
	}
	db.log.Info("Assembly RAM kits deleted by assembly", "assembly_id", assemblyID)
	return nil
}

// Assembly Drive methods (many-to-many)
func (db *DB) AddAssemblyDrive(ctx context.Context, element core.AssemblyDrive) error {
	query := `
        INSERT INTO assembly_drives (assembly_id, drive_id, mount_type)
        VALUES ($1, $2, $3)
        ON CONFLICT (assembly_id, drive_id) DO UPDATE
        SET mount_type = EXCLUDED.mount_type
    `
	_, err := db.conn.Exec(query, element.AssemblyID, element.DriveID, element.MountType)
	if err != nil {
		db.log.Error("Failed to add assembly drive", "assembly_id", element.AssemblyID, "drive_id", element.DriveID, "error", err)
		return err
	}
	db.log.Info("Assembly drive added", "assembly_id", element.AssemblyID, "drive_id", element.DriveID)
	return nil
}

func (db *DB) DeleteAssemblyDrive(ctx context.Context, assemblyID int, driveID int) error {
	query := `DELETE FROM assembly_drives WHERE assembly_id = $1 AND drive_id = $2`
	result, err := db.conn.Exec(query, assemblyID, driveID)
	if err != nil {
		db.log.Error("Failed to delete assembly drive", "assembly_id", assemblyID, "drive_id", driveID, "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		db.log.Error("Assembly drive not found for deletion", "assembly_id", assemblyID, "drive_id", driveID)
		return core.ErrNotFound
	}

	db.log.Info("Assembly drive deleted", "assembly_id", assemblyID, "drive_id", driveID)
	return nil
}

func (db *DB) DeleteAssemblyDrivesByAssembly(ctx context.Context, assemblyID int) error {
	query := `DELETE FROM assembly_drives WHERE assembly_id = $1`
	_, err := db.conn.Exec(query, assemblyID)
	if err != nil {
		db.log.Error("Failed to delete assembly drives by assembly", "assembly_id", assemblyID, "error", err)
		return err
	}
	db.log.Info("Assembly drives deleted by assembly", "assembly_id", assemblyID)
	return nil
}
