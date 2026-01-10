package repo

import (
	"context"
	"fmt"
	"strings"

	"compConfigurator/internal/db"
	domainErr "compConfigurator/internal/domain/errors"

	"errors"

	"github.com/jackc/pgx/v5"
)

type CPU struct {
	CPUId                  int64   `json:"cpu_id"`
	Name                   string  `json:"name"`
	Brand                  *string `json:"brand,omitempty"`
	SocketCode             string  `json:"socket_code"`
	TdpWatt                *int    `json:"tdp_watt,omitempty"`
	HasIntegratedGPU       *bool   `json:"has_integrated_gpu,omitempty"`
	SupportedRamType       *string `json:"supported_ram_type,omitempty"`
	SupportedRamFreqMaxMHz *int    `json:"supported_ram_freq_max_mhz,omitempty"`
	PcieVersion            *int    `json:"pcie_version,omitempty"`
}

type CPUsCatalogFilter struct {
	Query string

	SocketCode string
	RamType    string

	HasIGPU *bool

	MinRamFreqMHz *int
	MaxTdpWatt    *int
	MinPcieVer    *int

	Sort   string // name_asc | name_desc
	Limit  int
	Offset int
}

type CPUsRepo struct {
	pool *db.Pool
}

func NewCPUsRepo(pool *db.Pool) *CPUsRepo {
	return &CPUsRepo{pool: pool}
}

func (r *CPUsRepo) List(ctx context.Context, f CPUsCatalogFilter) ([]CPU, error) {
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 50
	}
	if f.Offset < 0 {
		f.Offset = 0
	}

	orderBy := "name ASC"
	if f.Sort == "name_desc" {
		orderBy = "name DESC"
	}

	var where []string
	var args []any
	arg := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf("$%d", len(args))
	}

	if strings.TrimSpace(f.Query) != "" {
		where = append(where, "name ILIKE "+arg("%"+strings.TrimSpace(f.Query)+"%"))
	}
	if strings.TrimSpace(f.SocketCode) != "" {
		where = append(where, "socket_code = "+arg(strings.TrimSpace(f.SocketCode)))
	}
	if strings.TrimSpace(f.RamType) != "" {
		// у тебя поле называется supported_ram_type (см. struct)
		where = append(where, "supported_ram_type = "+arg(strings.TrimSpace(f.RamType)))
	}
	if f.HasIGPU != nil {
		where = append(where, "has_integrated_gpu = "+arg(*f.HasIGPU))
	}
	if f.MinRamFreqMHz != nil {
		where = append(where, "supported_ram_freq_max_mhz >= "+arg(*f.MinRamFreqMHz))
	}
	if f.MaxTdpWatt != nil {
		where = append(where, "tdp_watt <= "+arg(*f.MaxTdpWatt))
	}
	if f.MinPcieVer != nil {
		where = append(where, "pcie_version >= "+arg(*f.MinPcieVer))
	}

	sql := `
SELECT
  cpu_id, name, brand, socket_code, tdp_watt, has_integrated_gpu,
  supported_ram_type, supported_ram_freq_max_mhz, pcie_version
FROM pc_configurator.cpus
`
	if len(where) > 0 {
		sql += "WHERE " + strings.Join(where, " AND ") + "\n"
	}
	sql += "ORDER BY " + orderBy + "\n"
	sql += fmt.Sprintf("LIMIT %s OFFSET %s", arg(f.Limit), arg(f.Offset))

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []CPU
	for rows.Next() {
		var cpu CPU
		if err := rows.Scan(
			&cpu.CPUId, &cpu.Name, &cpu.Brand, &cpu.SocketCode, &cpu.TdpWatt, &cpu.HasIntegratedGPU,
			&cpu.SupportedRamType, &cpu.SupportedRamFreqMaxMHz, &cpu.PcieVersion,
		); err != nil {
			return nil, err
		}
		out = append(out, cpu)
	}
	return out, rows.Err()
}

func (r *CPUsRepo) GetByID(ctx context.Context, id int64) (CPU, error) {
	var c CPU
	err := r.pool.QueryRow(ctx, `
SELECT
  cpu_id, name, brand, socket_code, tdp_watt, has_integrated_gpu,
  supported_ram_type, supported_ram_freq_max_mhz, pcie_version
FROM pc_configurator.cpus
WHERE cpu_id = $1
`, id).Scan(
		&c.CPUId, &c.Name, &c.Brand, &c.SocketCode, &c.TdpWatt, &c.HasIntegratedGPU,
		&c.SupportedRamType, &c.SupportedRamFreqMaxMHz, &c.PcieVersion,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CPU{}, domainErr.ErrNotFound
		}
		return CPU{}, err
	}
	return c, nil
}
