package repo

import (
	"context"
	"fmt"

	"compConfigurator/internal/db"
	domainErr "compConfigurator/internal/domain/errors"

	"errors"

	"github.com/jackc/pgx/v5"
)

type CloneRepo struct {
	pool *db.Pool
}

func NewCloneRepo(pool *db.Pool) *CloneRepo {
	return &CloneRepo{pool: pool}
}

// Создаёт копию публичной сборки srcAssemblyID для пользователя dstUserID.
// Возвращает новый assembly_id.
func (r *CloneRepo) ClonePublicAssembly(ctx context.Context, srcAssemblyID, dstUserID int64, newName string) (int64, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// 1) прочитать исходную публичную сборку
	var (
		name     string
		cpuID    *int64
		mbID     *int64
		psuID    *int64
		caseID   *int64
		gpuID    *int64
		coolerID *int64
	)
	err = tx.QueryRow(ctx, `
SELECT name, cpu_id, motherboard_id, psu_id, case_id, gpu_id, cooler_id
FROM pc_configurator.assemblies
WHERE assembly_id = $1 AND is_public = TRUE
`, srcAssemblyID).Scan(&name, &cpuID, &mbID, &psuID, &caseID, &gpuID, &coolerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, domainErr.ErrNotFound
		}
		return 0, err
	}

	if newName == "" {
		newName = fmt.Sprintf("Copy of %s", name)
	}

	// 2) создать новую сборку (is_public=false)
	var newAssemblyID int64
	err = tx.QueryRow(ctx, `
INSERT INTO pc_configurator.assemblies
  (user_id, name, is_public, cpu_id, motherboard_id, psu_id, case_id, gpu_id, cooler_id)
VALUES
  ($1, $2, FALSE, $3, $4, $5, $6, $7, $8)
RETURNING assembly_id
`, dstUserID, newName, cpuID, mbID, psuID, caseID, gpuID, coolerID).Scan(&newAssemblyID)
	if err != nil {
		return 0, err
	}

	// 3) копируем RAM kits
	_, err = tx.Exec(ctx, `
INSERT INTO pc_configurator.assembly_ram_kits (assembly_id, ram_kit_id)
SELECT $1, ark.ram_kit_id
FROM pc_configurator.assembly_ram_kits ark
WHERE ark.assembly_id = $2
ON CONFLICT DO NOTHING
`, newAssemblyID, srcAssemblyID)
	if err != nil {
		return 0, err
	}

	// 4) копируем drives
	_, err = tx.Exec(ctx, `
INSERT INTO pc_configurator.assembly_drives (assembly_id, drive_id, mount_type)
SELECT $1, ad.drive_id, ad.mount_type
FROM pc_configurator.assembly_drives ad
WHERE ad.assembly_id = $2
ON CONFLICT DO NOTHING
`, newAssemblyID, srcAssemblyID)
	if err != nil {
		return 0, err
	}

	// 5) commit
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return newAssemblyID, nil
}
