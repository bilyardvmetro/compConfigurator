package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"comp_config.com/services/search/core"
	_ "github.com/jackc/pgx/v5/stdlib"
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

func (db *DB) buildListQuery(baseQuery, tableName string, opts *core.ListOptions, allowedSortFields map[string]string) (string, []interface{}) {
	var query strings.Builder
	var args []interface{}
	argCounter := 1

	query.WriteString(baseQuery)

	if opts != nil && len(opts.Filters) > 0 {
		query.WriteString(" WHERE ")
		for i, filter := range opts.Filters {
			if i > 0 {
				query.WriteString(" AND ")
			}
			query.WriteString(filter.Field)
			query.WriteString(" ")
			query.WriteString(filter.Operator)
			query.WriteString(" $")
			query.WriteString(strconv.Itoa(argCounter))
			args = append(args, filter.Value)
			argCounter++
		}
	}

	if opts != nil && opts.SortBy != "" {
		if field, ok := allowedSortFields[opts.SortBy]; ok {
			query.WriteString(" ORDER BY ")
			query.WriteString(field)
			if opts.SortOrder == "desc" {
				query.WriteString(" DESC")
			} else {
				query.WriteString(" ASC")
			}
		}
	}

	if opts != nil && opts.PageSize > 0 {
		query.WriteString(" LIMIT $")
		query.WriteString(strconv.Itoa(argCounter))
		args = append(args, opts.PageSize)
		argCounter++

		if opts.Page > 1 {
			query.WriteString(" OFFSET $")
			query.WriteString(strconv.Itoa(argCounter))
			args = append(args, (opts.Page-1)*opts.PageSize)
			argCounter++
		}
	}

	return query.String(), args
}

func (db *DB) getTotalCount(ctx context.Context, tableName string, filters []core.Filter) (int, error) {
	var query strings.Builder
	var args []interface{}
	argCounter := 1

	query.WriteString("SELECT COUNT(*) FROM ")
	query.WriteString(tableName)

	if len(filters) > 0 {
		query.WriteString(" WHERE ")
		for i, filter := range filters {
			if i > 0 {
				query.WriteString(" AND ")
			}
			query.WriteString(filter.Field)
			query.WriteString(" ")
			query.WriteString(filter.Operator)
			query.WriteString(" $")
			query.WriteString(strconv.Itoa(argCounter))
			args = append(args, filter.Value)
			argCounter++
		}
	}

	var count int
	err := db.conn.GetContext(ctx, &count, query.String(), args...)
	if err != nil {
		db.log.Error("Failed to get total count", "table", tableName, "error", err)
		return 0, err
	}

	return count, nil
}

// CheckAssemblyCompatibility проверяет совместимость сборки
func (db *DB) CheckAssemblyCompatibility(ctx context.Context, assemblyID int) (*core.ValidationResult, error) {
	result := &core.ValidationResult{
		IsValid:  true,
		Errors:   make([]string, 0),
		Warnings: make([]string, 0),
	}

	// Вызываем хранимую процедуру
	_, err := db.conn.ExecContext(ctx, "SELECT check_assembly_compatibility($1)", assemblyID)
	if err != nil {
		// Проверяем, является ли ошибка ошибкой совместимости
		if strings.Contains(err.Error(), "check_violation") {
			result.IsValid = false
			// Извлекаем понятное сообщение об ошибке
			result.Errors = append(result.Errors, extractErrorMessage(err))
			return result, nil
		}
		db.log.Error("Failed to check assembly compatibility", "assembly_id", assemblyID, "error", err)
		return nil, err
	}

	return result, nil
}

// extractErrorMessage извлекает понятное сообщение из ошибки PostgreSQL
func extractErrorMessage(err error) string {
	errStr := err.Error()
	if idx := strings.Index(errStr, "ERROR: "); idx != -1 {
		errStr = errStr[idx+7:]
		if idx = strings.LastIndex(errStr, " ("); idx != -1 {
			errStr = errStr[:idx]
		}
	}
	return errStr
}

// ValidateAssemblyComponents проверяет совместимость конкретных компонентов
func (db *DB) ValidateAssemblyComponents(ctx context.Context, components core.AssemblyComponents) (*core.ValidationResult, error) {
	tx, err := db.conn.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Проверяем обязательные поля
	if components.CpuID == nil || components.MotherboardID == nil || components.PsuID == nil || components.CaseID == nil {
		return &core.ValidationResult{
			IsValid: false,
			Errors:  []string{"Не указаны обязательные компоненты (CPU, материнская плата, БП, корпус)"},
		}, nil
	}

	tempAssembly := core.Assembly{
		UserID:        -1,
		Name:          "temp_validation",
		IsPublic:      false,
		CpuID:         *components.CpuID,
		GpuID:         components.GpuID,
		MotherboardID: *components.MotherboardID,
		PsuID:         *components.PsuID,
		CaseID:        *components.CaseID,
		CoolerID:      components.CoolerID,
	}

	var assemblyID int
	err = tx.QueryRowContext(ctx, `
        INSERT INTO assemblies (user_id, name, is_public, cpu_id, gpu_id, motherboard_id, psu_id, case_id, cooler_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING assembly_id
    `,
		tempAssembly.UserID, tempAssembly.Name, tempAssembly.IsPublic,
		tempAssembly.CpuID, tempAssembly.GpuID, tempAssembly.MotherboardID,
		tempAssembly.PsuID, tempAssembly.CaseID, tempAssembly.CoolerID,
	).Scan(&assemblyID)

	if err != nil {
		return nil, err
	}

	for _, ramKitID := range components.RAMKitIDs {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO assembly_ram_kits (assembly_id, ram_kit_id)
            VALUES ($1, $2)
            ON CONFLICT DO NOTHING
        `, assemblyID, ramKitID)
		if err != nil {
			return nil, err
		}
	}

	for _, drive := range components.Drives {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO assembly_drives (assembly_id, drive_id, mount_type)
            VALUES ($1, $2, $3)
            ON CONFLICT DO NOTHING
        `, assemblyID, drive.DriveID, drive.MountType)
		if err != nil {
			return nil, err
		}
	}

	_, err = tx.ExecContext(ctx, "SELECT check_assembly_compatibility($1)", assemblyID)
	if err != nil {
		result := &core.ValidationResult{
			IsValid: false,
			Errors:  []string{extractErrorMessage(err)},
		}
		return result, nil
	}

	return &core.ValidationResult{IsValid: true}, nil
}

// Вспомогательная функция для создания ListResult
func createListResult[T any](items []T, totalCount int, opts *core.ListOptions) *core.ListResult[T] {
	totalPages := 0
	page := 1
	pageSize := 0

	if opts != nil {
		page = opts.Page
		pageSize = opts.PageSize
		if pageSize > 0 {
			totalPages = (totalCount + pageSize - 1) / pageSize
		}
	}

	return &core.ListResult[T]{
		Items:      items,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// GPU listing
func (db *DB) ListGPUs(ctx context.Context, opts *core.ListOptions) (*core.ListResult[core.GPU], error) {
	allowedSortFields := map[string]string{
		"id":         "gpu_id",
		"name":       "name",
		"brand":      "brand",
		"tdp_watt":   "tdp_watt",
		"length_mm":  "length_mm",
		"created_at": "created_at",
	}

	baseQuery := "SELECT gpu_id, name, brand, interface, pcie_version, pcie_lanes_required, tdp_watt, power_connectors, length_mm, width_slots, height_mm FROM gpus"

	query, args := db.buildListQuery(baseQuery, "gpus", opts, allowedSortFields)

	var items []core.GPU
	err := db.conn.SelectContext(ctx, &items, query, args...)
	if err != nil {
		db.log.Error("Failed to list GPUs", "error", err)
		return nil, err
	}

	totalCount, err := db.getTotalCount(ctx, "gpus", opts.Filters)
	if err != nil {
		return nil, err
	}

	return createListResult(items, totalCount, opts), nil
}

// CPU Cooler listing
func (db *DB) ListCpuCoolers(ctx context.Context, opts *core.ListOptions) (*core.ListResult[core.CpuCooler], error) {
	allowedSortFields := map[string]string{
		"id":             "cooler_id",
		"name":           "name",
		"brand":          "brand",
		"cooling_type":   "cooling_type",
		"tdp_limit_watt": "tdp_limit_watt",
		"height_mm":      "height_mm",
		"created_at":     "created_at",
	}

	baseQuery := "SELECT cooler_id, name, brand, cooling_type, tdp_limit_watt, height_mm, fan_count, fan_control_type FROM cpu_coolers"

	query, args := db.buildListQuery(baseQuery, "cpu_coolers", opts, allowedSortFields)

	var items []core.CpuCooler
	err := db.conn.SelectContext(ctx, &items, query, args...)
	if err != nil {
		db.log.Error("Failed to list CPU coolers", "error", err)
		return nil, err
	}

	totalCount, err := db.getTotalCount(ctx, "cpu_coolers", opts.Filters)
	if err != nil {
		return nil, err
	}

	return createListResult(items, totalCount, opts), nil
}

// Case listing
func (db *DB) ListCases(ctx context.Context, opts *core.ListOptions) (*core.ListResult[core.Case], error) {
	allowedSortFields := map[string]string{
		"id":                   "case_id",
		"name":                 "name",
		"brand":                "brand",
		"psu_form_factor":      "psu_form_factor",
		"max_gpu_length_mm":    "max_gpu_length_mm",
		"max_cooler_height_mm": "max_cooler_height_mm",
		"created_at":           "created_at",
	}

	baseQuery := "SELECT case_id, name, brand, psu_form_factor, max_gpu_length_mm, max_gpu_width_slots, max_cooler_height_mm, max_psu_length_mm, drive_bays_3_5_count, drive_bays_2_5_count FROM cases"

	query, args := db.buildListQuery(baseQuery, "cases", opts, allowedSortFields)

	var items []core.Case
	err := db.conn.SelectContext(ctx, &items, query, args...)
	if err != nil {
		db.log.Error("Failed to list cases", "error", err)
		return nil, err
	}

	totalCount, err := db.getTotalCount(ctx, "cases", opts.Filters)
	if err != nil {
		return nil, err
	}

	return createListResult(items, totalCount, opts), nil
}

// PSU listing
func (db *DB) ListPSUs(ctx context.Context, opts *core.ListOptions) (*core.ListResult[core.PSU], error) {
	allowedSortFields := map[string]string{
		"id":                "psu_id",
		"name":              "name",
		"brand":             "brand",
		"power_watt":        "power_watt",
		"form_factor":       "form_factor",
		"efficiency_rating": "efficiency_rating",
		"created_at":        "created_at",
	}

	baseQuery := "SELECT psu_id, name, brand, power_watt, form_factor, efficiency_rating, pcie_connectors_6_8pin_count, cpu_8pin_connectors_count, sata_connectors_count, molex_connectors_count, length_mm FROM psus"

	query, args := db.buildListQuery(baseQuery, "psus", opts, allowedSortFields)

	var items []core.PSU
	err := db.conn.SelectContext(ctx, &items, query, args...)
	if err != nil {
		db.log.Error("Failed to list PSUs", "error", err)
		return nil, err
	}

	totalCount, err := db.getTotalCount(ctx, "psus", opts.Filters)
	if err != nil {
		return nil, err
	}

	return createListResult(items, totalCount, opts), nil
}

// CPU Socket listing
func (db *DB) ListCpuSockets(ctx context.Context, opts *core.ListOptions) (*core.ListResult[core.CpuSocket], error) {
	allowedSortFields := map[string]string{
		"code":        "socket_code",
		"description": "description",
	}

	baseQuery := "SELECT socket_code, description FROM cpu_sockets"

	query, args := db.buildListQuery(baseQuery, "cpu_sockets", opts, allowedSortFields)

	var items []core.CpuSocket
	err := db.conn.SelectContext(ctx, &items, query, args...)
	if err != nil {
		db.log.Error("Failed to list CPU sockets", "error", err)
		return nil, err
	}

	totalCount, err := db.getTotalCount(ctx, "cpu_sockets", opts.Filters)
	if err != nil {
		return nil, err
	}

	return createListResult(items, totalCount, opts), nil
}

// Motherboard Form Factor listing
func (db *DB) ListMotherboardFormFactors(ctx context.Context, opts *core.ListOptions) (*core.ListResult[core.MotherboardFormFactor], error) {
	allowedSortFields := map[string]string{
		"code":        "form_factor_code",
		"description": "description",
		"width_mm":    "width_mm",
		"height_mm":   "height_mm",
	}

	baseQuery := "SELECT form_factor_code, description, width_mm, height_mm FROM motherboard_form_factors"

	query, args := db.buildListQuery(baseQuery, "motherboard_form_factors", opts, allowedSortFields)

	var items []core.MotherboardFormFactor
	err := db.conn.SelectContext(ctx, &items, query, args...)
	if err != nil {
		db.log.Error("Failed to list motherboard form factors", "error", err)
		return nil, err
	}

	totalCount, err := db.getTotalCount(ctx, "motherboard_form_factors", opts.Filters)
	if err != nil {
		return nil, err
	}

	return createListResult(items, totalCount, opts), nil
}

// CPU listing
func (db *DB) ListCPUs(ctx context.Context, opts *core.ListOptions) (*core.ListResult[core.CPU], error) {
	allowedSortFields := map[string]string{
		"id":             "cpu_id",
		"name":           "name",
		"brand":          "brand",
		"socket_code":    "socket_code",
		"core_count":     "core_count",
		"base_clock_mhz": "base_clock_mhz",
		"tdp_watt":       "tdp_watt",
		"created_at":     "created_at",
	}

	baseQuery := "SELECT cpu_id, name, brand, socket_code, architecture, core_count, thread_count, base_clock_mhz, boost_clock_mhz, tdp_watt, has_integrated_gpu, supported_ram_type, supported_ram_freq_max_mhz, memory_channels, pcie_version, pcie_lanes_total FROM cpus"

	query, args := db.buildListQuery(baseQuery, "cpus", opts, allowedSortFields)

	var items []core.CPU
	err := db.conn.SelectContext(ctx, &items, query, args...)
	if err != nil {
		db.log.Error("Failed to list CPUs", "error", err)
		return nil, err
	}

	totalCount, err := db.getTotalCount(ctx, "cpus", opts.Filters)
	if err != nil {
		return nil, err
	}

	return createListResult(items, totalCount, opts), nil
}

// Motherboard listing
func (db *DB) ListMotherboards(ctx context.Context, opts *core.ListOptions) (*core.ListResult[core.Motherboard], error) {
	allowedSortFields := map[string]string{
		"id":                  "motherboard_id",
		"name":                "name",
		"brand":               "brand",
		"socket_code":         "socket_code",
		"ram_capacity_max_gb": "ram_capacity_max_gb",
		"ram_freq_max_mhz":    "ram_freq_max_mhz",
		"created_at":          "created_at",
	}

	baseQuery := "SELECT motherboard_id, name, brand, socket_code, form_factor_code, chipset, ram_type, ram_slots, ram_capacity_max_gb, ram_freq_max_mhz, pcie_x16_slots_count, pcie_version_max, m2_slots_count, sata_ports_count, psu_main_connector_type, cpu_power_connector_type FROM motherboards"

	query, args := db.buildListQuery(baseQuery, "motherboards", opts, allowedSortFields)

	var items []core.Motherboard
	err := db.conn.SelectContext(ctx, &items, query, args...)
	if err != nil {
		db.log.Error("Failed to list motherboards", "error", err)
		return nil, err
	}

	totalCount, err := db.getTotalCount(ctx, "motherboards", opts.Filters)
	if err != nil {
		return nil, err
	}

	return createListResult(items, totalCount, opts), nil
}

// RAM Kit listing
func (db *DB) ListRAMKits(ctx context.Context, opts *core.ListOptions) (*core.ListResult[core.RAMKit], error) {
	allowedSortFields := map[string]string{
		"id":                "ram_kit_id",
		"name":              "name",
		"brand":             "brand",
		"ram_type":          "ram_type",
		"total_capacity_gb": "total_capacity_gb",
		"freq_mhz":          "freq_mhz",
		"created_at":        "created_at",
	}

	baseQuery := "SELECT ram_kit_id, name, brand, ram_type, module_capacity_gb, module_count, total_capacity_gb, freq_mhz, timings, voltage_v, form_factor FROM ram_kits"

	query, args := db.buildListQuery(baseQuery, "ram_kits", opts, allowedSortFields)

	var items []core.RAMKit
	err := db.conn.SelectContext(ctx, &items, query, args...)
	if err != nil {
		db.log.Error("Failed to list RAM kits", "error", err)
		return nil, err
	}

	totalCount, err := db.getTotalCount(ctx, "ram_kits", opts.Filters)
	if err != nil {
		return nil, err
	}

	return createListResult(items, totalCount, opts), nil
}

// Storage Drive listing
func (db *DB) ListStorageDrives(ctx context.Context, opts *core.ListOptions) (*core.ListResult[core.StorageDrive], error) {
	allowedSortFields := map[string]string{
		"id":          "drive_id",
		"name":        "name",
		"brand":       "brand",
		"drive_type":  "drive_type",
		"capacity_gb": "capacity_gb",
		"created_at":  "created_at",
	}

	baseQuery := "SELECT drive_id, name, brand, drive_type, form_factor, interface, capacity_gb FROM storage_drives"

	query, args := db.buildListQuery(baseQuery, "storage_drives", opts, allowedSortFields)

	var items []core.StorageDrive
	err := db.conn.SelectContext(ctx, &items, query, args...)
	if err != nil {
		db.log.Error("Failed to list storage drives", "error", err)
		return nil, err
	}

	totalCount, err := db.getTotalCount(ctx, "storage_drives", opts.Filters)
	if err != nil {
		return nil, err
	}

	return createListResult(items, totalCount, opts), nil
}

// User listing
func (db *DB) ListUsers(ctx context.Context, opts *core.ListOptions) (*core.ListResult[core.User], error) {
	allowedSortFields := map[string]string{
		"id":         "user_id",
		"email":      "email",
		"nickname":   "nickname",
		"created_at": "created_at",
	}

	baseQuery := "SELECT user_id, email, password_hash, nickname, avatar_url, is_admin, created_at FROM users"

	query, args := db.buildListQuery(baseQuery, "users", opts, allowedSortFields)

	var items []core.User
	err := db.conn.SelectContext(ctx, &items, query, args...)
	if err != nil {
		db.log.Error("Failed to list users", "error", err)
		return nil, err
	}

	totalCount, err := db.getTotalCount(ctx, "users", opts.Filters)
	if err != nil {
		return nil, err
	}

	return createListResult(items, totalCount, opts), nil
}

// Shop listing
func (db *DB) ListShops(ctx context.Context, opts *core.ListOptions) (*core.ListResult[core.Shop], error) {
	allowedSortFields := map[string]string{
		"id":   "shop_id",
		"name": "name",
		"url":  "url",
	}

	baseQuery := "SELECT shop_id, name, url FROM shops"

	query, args := db.buildListQuery(baseQuery, "shops", opts, allowedSortFields)

	var items []core.Shop
	err := db.conn.SelectContext(ctx, &items, query, args...)
	if err != nil {
		db.log.Error("Failed to list shops", "error", err)
		return nil, err
	}

	totalCount, err := db.getTotalCount(ctx, "shops", opts.Filters)
	if err != nil {
		return nil, err
	}

	return createListResult(items, totalCount, opts), nil
}

// Product Offer listing
func (db *DB) ListProductOffers(ctx context.Context, opts *core.ListOptions) (*core.ListResult[core.ProductOffer], error) {
	allowedSortFields := map[string]string{
		"id":             "offer_id",
		"shop_id":        "shop_id",
		"component_type": "component_type",
		"price":          "price",
		"created_at":     "created_at",
	}

	baseQuery := "SELECT offer_id, shop_id, component_type, component_id, price, available FROM product_offers"

	query, args := db.buildListQuery(baseQuery, "product_offers", opts, allowedSortFields)

	var items []core.ProductOffer
	err := db.conn.SelectContext(ctx, &items, query, args...)
	if err != nil {
		db.log.Error("Failed to list product offers", "error", err)
		return nil, err
	}

	totalCount, err := db.getTotalCount(ctx, "product_offers", opts.Filters)
	if err != nil {
		return nil, err
	}

	return createListResult(items, totalCount, opts), nil
}

// Assembly listing
func (db *DB) ListAssemblies(ctx context.Context, opts *core.ListOptions) (*core.ListResult[core.Assembly], error) {
	allowedSortFields := map[string]string{
		"id":          "assembly_id",
		"name":        "name",
		"user_id":     "user_id",
		"total_price": "total_price_cached",
		"created_at":  "created_at",
		"updated_at":  "updated_at",
	}

	baseQuery := "SELECT assembly_id, user_id, name, is_public, total_price_cached, cpu_id, gpu_id, motherboard_id, psu_id, case_id, cooler_id, created_at, updated_at FROM assemblies"

	query, args := db.buildListQuery(baseQuery, "assemblies", opts, allowedSortFields)

	var items []core.Assembly
	err := db.conn.SelectContext(ctx, &items, query, args...)
	if err != nil {
		db.log.Error("Failed to list assemblies", "error", err)
		return nil, err
	}

	totalCount, err := db.getTotalCount(ctx, "assemblies", opts.Filters)
	if err != nil {
		return nil, err
	}

	return createListResult(items, totalCount, opts), nil
}

// Assembly listing by user
func (db *DB) ListAssembliesByUser(ctx context.Context, userID int, opts *core.ListOptions) (*core.ListResult[core.Assembly], error) {
	if opts == nil {
		opts = &core.ListOptions{}
	}

	// Добавляем фильтр по user_id
	opts.Filters = append(opts.Filters, core.Filter{
		Field:    "user_id",
		Operator: "=",
		Value:    userID,
	})

	return db.ListAssemblies(ctx, opts)
}

// Get методы
func (db *DB) GetGPU(ctx context.Context, id int) (*core.GPU, error) {
	var gpu core.GPU
	query := `SELECT gpu_id, name, brand, interface, pcie_version, pcie_lanes_required, 
                     tdp_watt, power_connectors, length_mm, width_slots, height_mm 
              FROM gpus WHERE gpu_id = $1`

	err := db.conn.GetContext(ctx, &gpu, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		db.log.Error("Failed to get GPU", "id", id, "error", err)
		return nil, err
	}
	return &gpu, nil
}

func (db *DB) GetCPU(ctx context.Context, id int) (*core.CPU, error) {
	var cpu core.CPU
	query := `SELECT cpu_id, name, brand, socket_code, architecture, core_count,
                     thread_count, base_clock_mhz, boost_clock_mhz, tdp_watt,
                     has_integrated_gpu, supported_ram_type, supported_ram_freq_max_mhz,
                     memory_channels, pcie_version, pcie_lanes_total 
              FROM cpus WHERE cpu_id = $1`

	err := db.conn.GetContext(ctx, &cpu, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		db.log.Error("Failed to get CPU", "id", id, "error", err)
		return nil, err
	}
	return &cpu, nil
}

func (db *DB) GetMotherboard(ctx context.Context, id int) (*core.Motherboard, error) {
	var mb core.Motherboard
	query := `SELECT motherboard_id, name, brand, socket_code, form_factor_code,
                     chipset, ram_type, ram_slots, ram_capacity_max_gb, ram_freq_max_mhz,
                     pcie_x16_slots_count, pcie_version_max, m2_slots_count, 
                     sata_ports_count, psu_main_connector_type, cpu_power_connector_type 
              FROM motherboards WHERE motherboard_id = $1`

	err := db.conn.GetContext(ctx, &mb, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		db.log.Error("Failed to get motherboard", "id", id, "error", err)
		return nil, err
	}
	return &mb, nil
}

func (db *DB) GetRAMKit(ctx context.Context, id int) (*core.RAMKit, error) {
	var ram core.RAMKit
	query := `SELECT ram_kit_id, name, brand, ram_type, module_capacity_gb, module_count,
                     total_capacity_gb, freq_mhz, timings, voltage_v, form_factor 
              FROM ram_kits WHERE ram_kit_id = $1`

	err := db.conn.GetContext(ctx, &ram, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		db.log.Error("Failed to get RAM kit", "id", id, "error", err)
		return nil, err
	}
	return &ram, nil
}

func (db *DB) GetPSU(ctx context.Context, id int) (*core.PSU, error) {
	var psu core.PSU
	query := `SELECT psu_id, name, brand, power_watt, form_factor, efficiency_rating,
                     pcie_connectors_6_8pin_count, cpu_8pin_connectors_count,
                     sata_connectors_count, molex_connectors_count, length_mm 
              FROM psus WHERE psu_id = $1`

	err := db.conn.GetContext(ctx, &psu, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		db.log.Error("Failed to get PSU", "id", id, "error", err)
		return nil, err
	}
	return &psu, nil
}

func (db *DB) GetCase(ctx context.Context, id int) (*core.Case, error) {
	var c core.Case
	query := `SELECT case_id, name, brand, psu_form_factor, max_gpu_length_mm, 
                     max_gpu_width_slots, max_cooler_height_mm, max_psu_length_mm,
                     drive_bays_3_5_count, drive_bays_2_5_count 
              FROM cases WHERE case_id = $1`

	err := db.conn.GetContext(ctx, &c, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		db.log.Error("Failed to get case", "id", id, "error", err)
		return nil, err
	}
	return &c, nil
}

func (db *DB) GetCpuCooler(ctx context.Context, id int) (*core.CpuCooler, error) {
	var cooler core.CpuCooler
	query := `SELECT cooler_id, name, brand, cooling_type, tdp_limit_watt, 
                     height_mm, fan_count, fan_control_type 
              FROM cpu_coolers WHERE cooler_id = $1`

	err := db.conn.GetContext(ctx, &cooler, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		db.log.Error("Failed to get CPU cooler", "id", id, "error", err)
		return nil, err
	}
	return &cooler, nil
}

func (db *DB) GetStorageDrive(ctx context.Context, id int) (*core.StorageDrive, error) {
	var drive core.StorageDrive
	query := `SELECT drive_id, name, brand, drive_type, form_factor, interface, capacity_gb 
              FROM storage_drives WHERE drive_id = $1`

	err := db.conn.GetContext(ctx, &drive, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		db.log.Error("Failed to get storage drive", "id", id, "error", err)
		return nil, err
	}
	return &drive, nil
}

func (db *DB) GetCpuSocket(ctx context.Context, code string) (*core.CpuSocket, error) {
	var socket core.CpuSocket
	query := `SELECT socket_code, description FROM cpu_sockets WHERE socket_code = $1`

	err := db.conn.GetContext(ctx, &socket, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		db.log.Error("Failed to get CPU socket", "code", code, "error", err)
		return nil, err
	}
	return &socket, nil
}

func (db *DB) GetMotherboardFormFactor(ctx context.Context, code string) (*core.MotherboardFormFactor, error) {
	var ff core.MotherboardFormFactor
	query := `SELECT form_factor_code, description, width_mm, height_mm 
              FROM motherboard_form_factors WHERE form_factor_code = $1`

	err := db.conn.GetContext(ctx, &ff, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		db.log.Error("Failed to get motherboard form factor", "code", code, "error", err)
		return nil, err
	}
	return &ff, nil
}

func (db *DB) GetUser(ctx context.Context, id int) (*core.User, error) {
	var user core.User
	query := `SELECT user_id, email, password_hash, nickname, avatar_url, is_admin, created_at
              FROM users WHERE user_id = $1`

	err := db.conn.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		db.log.Error("Failed to get user", "id", id, "error", err)
		return nil, err
	}
	return &user, nil
}

func (db *DB) GetShop(ctx context.Context, id int) (*core.Shop, error) {
	var shop core.Shop
	query := `SELECT shop_id, name, url FROM shops WHERE shop_id = $1`

	err := db.conn.GetContext(ctx, &shop, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		db.log.Error("Failed to get shop", "id", id, "error", err)
		return nil, err
	}
	return &shop, nil
}

func (db *DB) GetProductOffer(ctx context.Context, id int) (*core.ProductOffer, error) {
	var offer core.ProductOffer
	query := `SELECT offer_id, shop_id, component_type, component_id, price, available 
              FROM product_offers WHERE offer_id = $1`

	err := db.conn.GetContext(ctx, &offer, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		db.log.Error("Failed to get product offer", "id", id, "error", err)
		return nil, err
	}
	return &offer, nil
}

func (db *DB) GetAssembly(ctx context.Context, id int) (*core.Assembly, error) {
	var assembly core.Assembly
	query := `SELECT assembly_id, user_id, name, is_public, total_price_cached,
                     cpu_id, gpu_id, motherboard_id, psu_id, case_id, cooler_id,
                     created_at, updated_at 
              FROM assemblies WHERE assembly_id = $1`

	err := db.conn.GetContext(ctx, &assembly, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		db.log.Error("Failed to get assembly", "id", id, "error", err)
		return nil, err
	}
	return &assembly, nil
}

func (db *DB) GetComponentByID(ctx context.Context, componentType string, id int) (interface{}, error) {
	switch componentType {
	case "gpu":
		return db.GetGPU(ctx, id)
	case "cpu":
		return db.GetCPU(ctx, id)
	case "motherboard":
		return db.GetMotherboard(ctx, id)
	case "ram_kit":
		return db.GetRAMKit(ctx, id)
	case "psu":
		return db.GetPSU(ctx, id)
	case "case":
		return db.GetCase(ctx, id)
	case "cooler":
		return db.GetCpuCooler(ctx, id)
	case "drive":
		return db.GetStorageDrive(ctx, id)
	case "user":
		return db.GetUser(ctx, id)
	case "shop":
		return db.GetShop(ctx, id)
	case "product_offer":
		return db.GetProductOffer(ctx, id)
	case "assembly":
		return db.GetAssembly(ctx, id)
	default:
		return nil, fmt.Errorf("unknown component type: %s", componentType)
	}
}

func (db *DB) SearchComponents(ctx context.Context, query string, componentTypes []string, limit, offset int, filters map[string]string) ([]core.SearchResult, int, error) {
	var results []core.SearchResult
	var totalCount int
	searchQuery := `
        WITH search_results AS (
            SELECT 'gpu' as component_type, gpu_id as id, name, brand, 
                   ts_rank_cd(to_tsvector('russian', name || ' ' || brand), plainto_tsquery('russian', $1)) as rank
            FROM gpus WHERE to_tsvector('russian', name || ' ' || brand) @@ plainto_tsquery('russian', $1)
            UNION ALL
            SELECT 'cpu' as component_type, cpu_id as id, name, brand,
                   ts_rank_cd(to_tsvector('russian', name || ' ' || brand), plainto_tsquery('russian', $1)) as rank
            FROM cpus WHERE to_tsvector('russian', name || ' ' || brand) @@ plainto_tsquery('russian', $1)
            UNION ALL
            SELECT 'motherboard' as component_type, motherboard_id as id, name, brand,
                   ts_rank_cd(to_tsvector('russian', name || ' ' || brand), plainto_tsquery('russian', $1)) as rank
            FROM motherboards WHERE to_tsvector('russian', name || ' ' || brand) @@ plainto_tsquery('russian', $1)
            -- Добавьте остальные таблицы по аналогии
        )
        SELECT component_type, id, name, brand, rank 
        FROM search_results 
        ORDER BY rank DESC
        LIMIT $2 OFFSET $3
    `

	err := db.conn.SelectContext(ctx, &results, searchQuery, query, limit, offset)
	if err != nil {
		db.log.Error("Failed to search components", "query", query, "error", err)
		return nil, 0, err
	}

	countQuery := `
        SELECT COUNT(*) FROM (
            SELECT gpu_id FROM gpus WHERE to_tsvector('russian', name || ' ' || brand) @@ plainto_tsquery('russian', $1)
            UNION ALL
            SELECT cpu_id FROM cpus WHERE to_tsvector('russian', name || ' ' || brand) @@ plainto_tsquery('russian', $1)
            UNION ALL
            SELECT motherboard_id FROM motherboards WHERE to_tsvector('russian', name || ' ' || brand) @@ plainto_tsquery('russian', $1)
        ) as total
    `

	err = db.conn.GetContext(ctx, &totalCount, countQuery, query)
	if err != nil {
		return nil, 0, err
	}

	return results, totalCount, nil
}

func (db *DB) GetComponentPrices(ctx context.Context, componentType string, componentID int) ([]core.PriceInfo, error) {
	var prices []core.PriceInfo

	query := `
        SELECT po.offer_id, po.shop_id, s.name as shop_name, po.price, po.available, 
               s.url, po.updated_at as last_updated
        FROM product_offers po
        JOIN shops s ON po.shop_id = s.shop_id
        WHERE po.component_type = $1 AND po.component_id = $2
        ORDER BY po.price ASC, po.available DESC
    `

	err := db.conn.SelectContext(ctx, &prices, query, componentType, componentID)
	if err != nil {
		db.log.Error("Failed to get component prices", "type", componentType, "id", componentID, "error", err)
		return nil, err
	}

	return prices, nil
}

func (db *DB) GetBestOffers(ctx context.Context, componentType string, componentID int, limit int, cheapestFirst bool) ([]core.Offer, error) {
	var offers []core.Offer

	orderBy := "po.price DESC"
	if cheapestFirst {
		orderBy = "po.price ASC"
	}

	query := fmt.Sprintf(`
        SELECT po.offer_id, po.shop_id, s.name as shop_name, po.price, 
               po.available, s.url, true as in_stock, '1-3 дня' as delivery_time
        FROM product_offers po
        JOIN shops s ON po.shop_id = s.shop_id
        WHERE po.component_type = $1 AND po.component_id = $2 AND po.available = true
        ORDER BY %s
        LIMIT $3
    `, orderBy)

	err := db.conn.SelectContext(ctx, &offers, query, componentType, componentID, limit)
	if err != nil {
		db.log.Error("Failed to get best offers", "type", componentType, "id", componentID, "error", err)
		return nil, err
	}

	return offers, nil
}
