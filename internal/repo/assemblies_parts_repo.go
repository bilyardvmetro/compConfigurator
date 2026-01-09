package repo

import (
	"context"
	"errors"

	"compConfigurator/internal/db"
	domainerr "compConfigurator/internal/domain/errors"

	"github.com/jackc/pgx/v5"
)

type AssemblyPartsRepo struct {
	pool *db.Pool
}

func NewAssemblyPartsRepo(pool *db.Pool) *AssemblyPartsRepo {
	return &AssemblyPartsRepo{pool: pool}
}

// ---- RAM kits ----

func (r *AssemblyPartsRepo) AddRamKit(ctx context.Context, assemblyID, ramKitID int64) error {
	_, err := r.pool.Exec(ctx, `
INSERT INTO pc_configurator.assembly_ram_kits (assembly_id, ram_kit_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`, assemblyID, ramKitID)
	return err
}

func (r *AssemblyPartsRepo) RemoveRamKit(ctx context.Context, assemblyID, ramKitID int64) error {
	ct, err := r.pool.Exec(ctx, `
DELETE FROM pc_configurator.assembly_ram_kits
WHERE assembly_id = $1 AND ram_kit_id = $2
`, assemblyID, ramKitID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domainerr.ErrNotFound
	}
	return nil
}

// ---- Drives ----

func (r *AssemblyPartsRepo) AddDrive(ctx context.Context, assemblyID, driveID int64, mountType *string) error {
	_, err := r.pool.Exec(ctx, `
INSERT INTO pc_configurator.assembly_drives (assembly_id, drive_id, mount_type)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING
`, assemblyID, driveID, mountType)
	return err
}

func (r *AssemblyPartsRepo) RemoveDrive(ctx context.Context, assemblyID, driveID int64) error {
	ct, err := r.pool.Exec(ctx, `
DELETE FROM pc_configurator.assembly_drives
WHERE assembly_id = $1 AND drive_id = $2
`, assemblyID, driveID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return domainerr.ErrNotFound
	}
	return nil
}

// ---- helpers: existence checks (чтобы отдавать 404/400 красиво) ----

func (r *AssemblyPartsRepo) EnsureRamKitExists(ctx context.Context, ramKitID int64) error {
	var tmp int64
	err := r.pool.QueryRow(ctx, `SELECT ram_kit_id FROM pc_configurator.ram_kits WHERE ram_kit_id=$1`, ramKitID).Scan(&tmp)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerr.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *AssemblyPartsRepo) EnsureDriveExists(ctx context.Context, driveID int64) error {
	var tmp int64
	err := r.pool.QueryRow(ctx, `SELECT drive_id FROM pc_configurator.storage_drives WHERE drive_id=$1`, driveID).Scan(&tmp)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainerr.ErrNotFound
		}
		return err
	}
	return nil
}
