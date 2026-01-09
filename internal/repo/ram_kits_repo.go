package repo

import (
	"context"
	"fmt"
	"strings"

	"compConfigurator/internal/db"
)

type RamKitCatalogItem struct {
	RamKitID        int64   `json:"ram_kit_id"`
	Name            string  `json:"name"`
	RamType         *string `json:"ram_type,omitempty"`
	ModuleCount     *int    `json:"module_count,omitempty"`
	TotalCapacityGb *int    `json:"total_capacity_gb,omitempty"`
	FreqMHz         *int    `json:"freq_mhz,omitempty"`
}

type RamKitsCatalogFilter struct {
	Query         string
	RamType       string
	MinFreqMHz    *int
	MaxFreqMHz    *int
	MinCapacityGb *int
	MaxCapacityGb *int

	Sort   string // name_asc | name_desc
	Limit  int
	Offset int
}

type RamKitsRepo struct {
	pool *db.Pool
}

func NewRamKitsRepo(pool *db.Pool) *RamKitsRepo {
	return &RamKitsRepo{pool: pool}
}

func (r *RamKitsRepo) List(ctx context.Context, f RamKitsCatalogFilter) ([]RamKitCatalogItem, error) {
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
	if strings.TrimSpace(f.RamType) != "" {
		where = append(where, "ram_type = "+arg(strings.TrimSpace(f.RamType)))
	}
	if f.MinFreqMHz != nil {
		where = append(where, "freq_mhz >= "+arg(*f.MinFreqMHz))
	}
	if f.MaxFreqMHz != nil {
		where = append(where, "freq_mhz <= "+arg(*f.MaxFreqMHz))
	}
	if f.MinCapacityGb != nil {
		where = append(where, "total_capacity_gb >= "+arg(*f.MinCapacityGb))
	}
	if f.MaxCapacityGb != nil {
		where = append(where, "total_capacity_gb <= "+arg(*f.MaxCapacityGb))
	}

	sql := `
SELECT ram_kit_id, name, ram_type, module_count, total_capacity_gb, freq_mhz
FROM pc_configurator.ram_kits
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

	var out []RamKitCatalogItem
	for rows.Next() {
		var x RamKitCatalogItem
		if err := rows.Scan(&x.RamKitID, &x.Name, &x.RamType, &x.ModuleCount, &x.TotalCapacityGb, &x.FreqMHz); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
