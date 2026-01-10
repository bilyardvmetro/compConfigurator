package repo

import (
	"context"
	"fmt"
	"strings"

	"compConfigurator/internal/db"
)

type GPUCatalogItem struct {
	GPUID      int64    `json:"gpu_id"`
	Name       string   `json:"name"`
	Brand      *string  `json:"brand,omitempty"`
	LengthMM   *int     `json:"length_mm,omitempty"`
	WidthSlots *float64 `json:"width_slots,omitempty"` // NUMERIC(3,1) -> float64
	TDPWatt    *int     `json:"tdp_watt,omitempty"`
}

type GPUsCatalogFilter struct {
	Query string

	MinLengthMM *int
	MaxLengthMM *int

	MinWidthSlots *float64
	MaxWidthSlots *float64

	MinTDPWatt *int
	MaxTDPWatt *int

	Sort   string
	Limit  int
	Offset int
}

type GPUsRepo struct {
	pool *db.Pool
}

func NewGPUsRepo(pool *db.Pool) *GPUsRepo {
	return &GPUsRepo{pool: pool}
}

func (r *GPUsRepo) List(ctx context.Context, f GPUsCatalogFilter) ([]GPUCatalogItem, error) {
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
	if f.MinLengthMM != nil {
		where = append(where, "length_mm >= "+arg(*f.MinLengthMM))
	}
	if f.MaxLengthMM != nil {
		where = append(where, "length_mm <= "+arg(*f.MaxLengthMM))
	}
	if f.MinWidthSlots != nil {
		where = append(where, "width_slots >= "+arg(*f.MinWidthSlots))
	}
	if f.MaxWidthSlots != nil {
		where = append(where, "width_slots <= "+arg(*f.MaxWidthSlots))
	}
	if f.MinTDPWatt != nil {
		where = append(where, "tdp_watt >= "+arg(*f.MinTDPWatt))
	}
	if f.MaxTDPWatt != nil {
		where = append(where, "tdp_watt <= "+arg(*f.MaxTDPWatt))
	}

	sql := `
SELECT
  gpu_id,
  name,
  brand,
  length_mm,
  width_slots,
  tdp_watt
FROM pc_configurator.gpus
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

	var out []GPUCatalogItem
	for rows.Next() {
		var x GPUCatalogItem
		if err := rows.Scan(&x.GPUID, &x.Name, &x.Brand, &x.LengthMM, &x.WidthSlots, &x.TDPWatt); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
