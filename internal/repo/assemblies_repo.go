package repo

import (
	"compConfigurator/internal/db"
	domainerr "compConfigurator/internal/domain/errors"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Assembly struct {
	AssemblyID       int64   `json:"assembly_id"`
	UserID           int64   `json:"user_id"`
	Name             string  `json:"name"`
	IsPublic         bool    `json:"is_public"`
	TotalPriceCached *string `json:"total_price_cached,omitempty"` // NUMERIC -> text

	CPUId         int64  `json:"cpu_id"`
	MotherboardId int64  `json:"motherboard_id"`
	PSUId         int64  `json:"psu_id"`
	CaseId        int64  `json:"case_id"`
	GPUId         *int64 `json:"gpu_id,omitempty"`
	CoolerId      *int64 `json:"cooler_id,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CompatibilityResult struct {
	IsOK         bool    `json:"is_ok"`
	ErrorMessage *string `json:"error_message,omitempty"`
}

type AssembliesRepo struct {
	pool *db.Pool
}

func NewAssembliesRepo(pool *db.Pool) *AssembliesRepo {
	return &AssembliesRepo{pool: pool}
}

func (r *AssembliesRepo) Create(
	ctx context.Context,
	userID int64,
	name string,
	cpuID, motherboardID, psuID, caseID int64,
	gpuID, coolerID *int64,
) (Assembly, error) {
	var id int64
	err := r.pool.QueryRow(ctx, `
	INSERT INTO pc_configurator.assemblies
	  (user_id, name, is_public, cpu_id, motherboard_id, psu_id, case_id, gpu_id, cooler_id)
	VALUES
	  ($1, $2, false, $3, $4, $5, $6, $7, $8)
	RETURNING assembly_id
`, userID, name, cpuID, motherboardID, psuID, caseID, gpuID, coolerID).Scan(&id)

	if err != nil {
		return Assembly{}, err
	}

	return r.GetByID(ctx, id, userID)
}

func (r *AssembliesRepo) ListByUser(ctx context.Context, userID int64, limit, offset int) ([]Assembly, error) {
	rows, err := r.pool.Query(ctx, `
	SELECT
	  assembly_id, user_id, name, is_public, total_price_cached::text,
	  cpu_id, gpu_id, motherboard_id, psu_id, case_id, cooler_id,
	  created_at, updated_at
	FROM pc_configurator.assemblies
	WHERE user_id = $1
	ORDER BY updated_at DESC
	LIMIT $2 OFFSET $3
`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Assembly
	for rows.Next() {
		var a Assembly
		if err := rows.Scan(
			&a.AssemblyID, &a.UserID, &a.Name, &a.IsPublic, &a.TotalPriceCached,
			&a.CPUId, &a.GPUId, &a.MotherboardId, &a.PSUId, &a.CaseId, &a.CoolerId,
			&a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *AssembliesRepo) GetByID(ctx context.Context, assemblyID, userID int64) (Assembly, error) {
	var a Assembly
	err := r.pool.QueryRow(ctx, `
	SELECT
	  assembly_id, user_id, name, is_public, total_price_cached::text,
	  cpu_id, gpu_id, motherboard_id, psu_id, case_id, cooler_id,
	  created_at, updated_at
	FROM pc_configurator.assemblies
	WHERE assembly_id = $1 AND user_id = $2
`, assemblyID, userID).Scan(
		&a.AssemblyID, &a.UserID, &a.Name, &a.IsPublic, &a.TotalPriceCached,
		&a.CPUId, &a.GPUId, &a.MotherboardId, &a.PSUId, &a.CaseId, &a.CoolerId,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Assembly{}, domainerr.ErrNotFound
		}
		return Assembly{}, err
	}
	return a, nil
}

func (r *AssembliesRepo) GetCompatibility(ctx context.Context, assemblyID int64) (CompatibilityResult, error) {
	var res CompatibilityResult
	err := r.pool.QueryRow(ctx, `
	SELECT is_ok, error_message
	FROM pc_configurator.is_assembly_compatible($1)
`, assemblyID).Scan(&res.IsOK, &res.ErrorMessage)
	if err != nil {
		return CompatibilityResult{}, err
	}
	return res, nil
}

func (r *AssembliesRepo) RecalcTotalPrice(ctx context.Context, assemblyID int64) error {
	_, err := r.pool.Exec(ctx, `SELECT pc_configurator.recalc_assembly_total_price($1)`, assemblyID)
	return err
}
