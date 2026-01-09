package repo

import (
	"context"

	"compConfigurator/internal/db"
)

type RamKit struct {
	RamKitID        int64   `json:"ram_kit_id"`
	Name            string  `json:"name"`
	RamType         *string `json:"ram_type,omitempty"`
	ModuleCount     *int    `json:"module_count,omitempty"`
	TotalCapacityGb *int    `json:"total_capacity_gb,omitempty"`
	FreqMHz         *int    `json:"freq_mhz,omitempty"`
	FormFactor      *string `json:"form_factor,omitempty"`
}

type Drive struct {
	DriveID    int64   `json:"drive_id"`
	Name       string  `json:"name"`
	Interface  *string `json:"interface,omitempty"`   // NVMe/SATA
	FormFactor *string `json:"form_factor,omitempty"` // 2.5/3.5/2280...
	CapacityGb *int    `json:"capacity_gb,omitempty"`
	MountType  *string `json:"mount_type,omitempty"` // если есть в assembly_drives
}

type AssemblyDetailsRepo struct {
	pool *db.Pool
}

func NewAssemblyDetailsRepo(pool *db.Pool) *AssemblyDetailsRepo {
	return &AssemblyDetailsRepo{pool: pool}
}

func (r *AssemblyDetailsRepo) ListRamKits(ctx context.Context, assemblyID int64) ([]RamKit, error) {
	rows, err := r.pool.Query(ctx, `
SELECT k.ram_kit_id, k.name, k.ram_type, k.module_count, k.total_capacity_gb, k.freq_mhz, k.form_factor
FROM pc_configurator.assembly_ram_kits ark
JOIN pc_configurator.ram_kits k ON k.ram_kit_id = ark.ram_kit_id
WHERE ark.assembly_id = $1
ORDER BY k.name
`, assemblyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []RamKit
	for rows.Next() {
		var x RamKit
		if err := rows.Scan(&x.RamKitID, &x.Name, &x.RamType, &x.ModuleCount, &x.TotalCapacityGb, &x.FreqMHz, &x.FormFactor); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *AssemblyDetailsRepo) ListDrives(ctx context.Context, assemblyID int64) ([]Drive, error) {
	rows, err := r.pool.Query(ctx, `
SELECT d.drive_id, d.name, d.interface, d.form_factor, d.capacity_gb, ad.mount_type
FROM pc_configurator.assembly_drives ad
JOIN pc_configurator.storage_drives d ON d.drive_id = ad.drive_id
WHERE ad.assembly_id = $1
ORDER BY d.name
`, assemblyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Drive
	for rows.Next() {
		var x Drive
		if err := rows.Scan(&x.DriveID, &x.Name, &x.Interface, &x.FormFactor, &x.CapacityGb, &x.MountType); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
