package repo

import (
	"context"
	"errors"

	"compConfigurator/internal/db"
	domainerr "compConfigurator/internal/domain/errors"

	"github.com/jackc/pgx/v5"
)

type AssemblyComponents struct {
	AssemblyID int64 `json:"assembly_id"`

	CPUId   *int64  `json:"cpu_id,omitempty"`
	CPUName *string `json:"cpu_name,omitempty"`

	GPUId   *int64  `json:"gpu_id,omitempty"`
	GPUName *string `json:"gpu_name,omitempty"`

	MotherboardId   *int64  `json:"motherboard_id,omitempty"`
	MotherboardName *string `json:"motherboard_name,omitempty"`

	PSUId   *int64  `json:"psu_id,omitempty"`
	PSUName *string `json:"psu_name,omitempty"`

	CaseId   *int64  `json:"case_id,omitempty"`
	CaseName *string `json:"case_name,omitempty"`

	CoolerId   *int64  `json:"cooler_id,omitempty"`
	CoolerName *string `json:"cooler_name,omitempty"`
}

type AssemblyComponentsRepo struct {
	pool *db.Pool
}

func NewAssemblyComponentsRepo(pool *db.Pool) *AssemblyComponentsRepo {
	return &AssemblyComponentsRepo{pool: pool}
}

// Проверяем владельца (как и в остальных owner-endpoint’ах)
func (r *AssemblyComponentsRepo) GetForOwner(ctx context.Context, assemblyID, userID int64) (AssemblyComponents, error) {
	var x AssemblyComponents

	// ⚠️ Если у тебя названия таблиц отличаются — поправь здесь.
	err := r.pool.QueryRow(ctx, `
SELECT
  a.assembly_id,

  a.cpu_id, c.name,
  a.gpu_id, g.name,
  a.motherboard_id, m.name,
  a.psu_id, p.name,
  a.case_id, ca.name,
  a.cooler_id, co.name

FROM pc_configurator.assemblies a
LEFT JOIN pc_configurator.cpus c ON c.cpu_id = a.cpu_id
LEFT JOIN pc_configurator.gpus g ON g.gpu_id = a.gpu_id
LEFT JOIN pc_configurator.motherboards m ON m.motherboard_id = a.motherboard_id
LEFT JOIN pc_configurator.psus p ON p.psu_id = a.psu_id
LEFT JOIN pc_configurator.cases ca ON ca.case_id = a.case_id
LEFT JOIN pc_configurator.cpu_coolers co ON co.cooler_id = a.cooler_id

WHERE a.assembly_id = $1 AND a.user_id = $2
`, assemblyID, userID).Scan(
		&x.AssemblyID,

		&x.CPUId, &x.CPUName,
		&x.GPUId, &x.GPUName,
		&x.MotherboardId, &x.MotherboardName,
		&x.PSUId, &x.PSUName,
		&x.CaseId, &x.CaseName,
		&x.CoolerId, &x.CoolerName,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AssemblyComponents{}, domainerr.ErrNotFound
		}
		return AssemblyComponents{}, err
	}
	return x, nil
}
