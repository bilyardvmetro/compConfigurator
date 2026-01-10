package repo

import (
	"context"
	"fmt"
	"strings"

	"compConfigurator/internal/db"
)

type CaseCatalogItem struct {
	CaseID        int64   `json:"case_id"`
	Name          string  `json:"name"`
	Brand         *string `json:"brand,omitempty"`
	PSUFormFactor *string `json:"psu_form_factor,omitempty"`

	MaxGPULengthMM    *int     `json:"max_gpu_length_mm,omitempty"`
	MaxGPUWidthSlots  *float64 `json:"max_gpu_width_slots,omitempty"`
	MaxCoolerHeightMM *int     `json:"max_cooler_height_mm,omitempty"`
	MaxPSULengthMM    *int     `json:"max_psu_length_mm,omitempty"`

	DriveBays35Count *int `json:"drive_bays_3_5_count,omitempty"`
	DriveBays25Count *int `json:"drive_bays_2_5_count,omitempty"`
}

type CasesCatalogFilter struct {
	Query          string
	FormFactorCode string // ← фильтр через M:N
	PSUFormFactor  string

	MinMaxGPULengthMM    *int
	MinMaxGPUWidthSlots  *float64
	MinMaxCoolerHeightMM *int
	MinMaxPSULengthMM    *int
	MinDriveBays35       *int
	MinDriveBays25       *int

	Sort   string
	Limit  int
	Offset int
}

type CasesRepo struct {
	pool *db.Pool
}

func NewCasesRepo(pool *db.Pool) *CasesRepo {
	return &CasesRepo{pool: pool}
}

func (r *CasesRepo) List(ctx context.Context, f CasesCatalogFilter) ([]CaseCatalogItem, error) {
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 50
	}
	if f.Offset < 0 {
		f.Offset = 0
	}

	orderBy := "c.name ASC"
	if f.Sort == "name_desc" {
		orderBy = "c.name DESC"
	}

	var where []string
	var joins []string
	var args []any

	arg := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf("$%d", len(args))
	}

	// ---- JOIN под форм-фактор платы
	if strings.TrimSpace(f.FormFactorCode) != "" {
		joins = append(joins, `
JOIN pc_configurator.case_form_factor_support cffs
  ON cffs.case_id = c.case_id
`)
		where = append(where, "cffs.form_factor_code = "+arg(f.FormFactorCode))
	}

	// ---- фильтры
	if strings.TrimSpace(f.Query) != "" {
		where = append(where, "c.name ILIKE "+arg("%"+strings.TrimSpace(f.Query)+"%"))
	}
	if strings.TrimSpace(f.PSUFormFactor) != "" {
		where = append(where, "c.psu_form_factor = "+arg(f.PSUFormFactor))
	}
	if f.MinMaxGPULengthMM != nil {
		where = append(where, "c.max_gpu_length_mm >= "+arg(*f.MinMaxGPULengthMM))
	}
	if f.MinMaxGPUWidthSlots != nil {
		where = append(where, "c.max_gpu_width_slots >= "+arg(*f.MinMaxGPUWidthSlots))
	}
	if f.MinMaxCoolerHeightMM != nil {
		where = append(where, "c.max_cooler_height_mm >= "+arg(*f.MinMaxCoolerHeightMM))
	}
	if f.MinMaxPSULengthMM != nil {
		where = append(where, "c.max_psu_length_mm >= "+arg(*f.MinMaxPSULengthMM))
	}
	if f.MinDriveBays35 != nil {
		where = append(where, "c.drive_bays_3_5_count >= "+arg(*f.MinDriveBays35))
	}
	if f.MinDriveBays25 != nil {
		where = append(where, "c.drive_bays_2_5_count >= "+arg(*f.MinDriveBays25))
	}

	sql := `
SELECT DISTINCT
  c.case_id,
  c.name,
  c.brand,
  c.psu_form_factor,
  c.max_gpu_length_mm,
  c.max_gpu_width_slots,
  c.max_cooler_height_mm,
  c.max_psu_length_mm,
  c.drive_bays_3_5_count,
  c.drive_bays_2_5_count
FROM pc_configurator.cases c
`

	if len(joins) > 0 {
		sql += strings.Join(joins, "\n") + "\n"
	}
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

	var out []CaseCatalogItem
	for rows.Next() {
		var x CaseCatalogItem
		if err := rows.Scan(
			&x.CaseID,
			&x.Name,
			&x.Brand,
			&x.PSUFormFactor,
			&x.MaxGPULengthMM,
			&x.MaxGPUWidthSlots,
			&x.MaxCoolerHeightMM,
			&x.MaxPSULengthMM,
			&x.DriveBays35Count,
			&x.DriveBays25Count,
		); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
