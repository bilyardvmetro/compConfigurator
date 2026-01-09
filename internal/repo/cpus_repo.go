package repo

import (
	"compConfigurator/internal/db"
	domainErr "compConfigurator/internal/domain/errors"
	"context"
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

type CPUsRepo struct {
	pool *db.Pool
}

func NewCPUsRepo(pool *db.Pool) *CPUsRepo {
	return &CPUsRepo{pool: pool}
}

func (r *CPUsRepo) List(ctx context.Context, limit, offset int) ([]CPU, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT 
		cpu_id, name, brand, socket_code, tdp_watt, has_integrated_gpu,
		supported_ram_type, supported_ram_freq_max_mhz, pcie_version
		FROM pc_configurator.cpus
		ORDER BY name
		LIMIT $1 OFFSET $2
`, limit, offset)

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
