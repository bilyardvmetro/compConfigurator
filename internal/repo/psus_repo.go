package repo

import (
	"context"
	"fmt"
	"strings"

	"compConfigurator/internal/db"
)

type PSUCatalogItem struct {
	PSUID      int64   `json:"psu_id"`
	Name       string  `json:"name"`
	Brand      *string `json:"brand,omitempty"`
	PowerWatt  int     `json:"power_watt"`
	FormFactor *string `json:"form_factor,omitempty"`
	LengthMM   *int    `json:"length_mm,omitempty"`
}

type PSUsCatalogFilter struct {
	Query        string
	FormFactor   string
	MinPowerWatt *int
	MaxLengthMM  *int

	Sort   string
	Limit  int
	Offset int
}

type PSUsRepo struct {
	pool *db.Pool
}

func NewPSUsRepo(pool *db.Pool) *PSUsRepo {
	return &PSUsRepo{pool: pool}
}

func (r *PSUsRepo) List(ctx context.Context, f PSUsCatalogFilter) ([]PSUCatalogItem, error) {
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
	if strings.TrimSpace(f.FormFactor) != "" {
		where = append(where, "form_factor = "+arg(f.FormFactor))
	}
	if f.MinPowerWatt != nil {
		where = append(where, "power_watt >= "+arg(*f.MinPowerWatt))
	}
	if f.MaxLengthMM != nil {
		where = append(where, "length_mm <= "+arg(*f.MaxLengthMM))
	}

	sql := `
SELECT
  psu_id,
  name,
  brand,
  power_watt,
  form_factor,
  length_mm
FROM pc_configurator.psus
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

	var out []PSUCatalogItem
	for rows.Next() {
		var x PSUCatalogItem
		if err := rows.Scan(
			&x.PSUID,
			&x.Name,
			&x.Brand,
			&x.PowerWatt,
			&x.FormFactor,
			&x.LengthMM,
		); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
