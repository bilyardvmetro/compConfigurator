package repo

import (
	"context"
	"fmt"
	"strings"

	"compConfigurator/internal/db"
)

type CPUCoolerCatalogItem struct {
	CoolerID     int64   `json:"cooler_id"`
	Name         string  `json:"name"`
	Brand        *string `json:"brand,omitempty"`
	HeightMM     *int    `json:"height_mm,omitempty"`
	TDPLimitWatt *int    `json:"tdp_limit_watt,omitempty"`
}

type CPUCoolersCatalogFilter struct {
	Query           string
	SocketCode      string
	MinTDPLimitWatt *int
	MaxHeightMM     *int

	Sort   string
	Limit  int
	Offset int
}

type CPUCoolersRepo struct {
	pool *db.Pool
}

func NewCPUCoolersRepo(pool *db.Pool) *CPUCoolersRepo {
	return &CPUCoolersRepo{pool: pool}
}

func (r *CPUCoolersRepo) List(ctx context.Context, f CPUCoolersCatalogFilter) ([]CPUCoolerCatalogItem, error) {
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

	var joins []string
	var where []string
	var args []any
	arg := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf("$%d", len(args))
	}

	// Фильтр по сокету через M:N (cooler_sockets)
	if strings.TrimSpace(f.SocketCode) != "" {
		joins = append(joins, `
JOIN pc_configurator.cooler_sockets cs
  ON cs.cooler_id = c.cooler_id
`)
		where = append(where, "cs.socket_code = "+arg(strings.TrimSpace(f.SocketCode)))
	}

	if strings.TrimSpace(f.Query) != "" {
		where = append(where, "c.name ILIKE "+arg("%"+strings.TrimSpace(f.Query)+"%"))
	}
	if f.MinTDPLimitWatt != nil {
		where = append(where, "c.tdp_limit_watt >= "+arg(*f.MinTDPLimitWatt))
	}
	if f.MaxHeightMM != nil {
		where = append(where, "c.height_mm <= "+arg(*f.MaxHeightMM))
	}

	sql := `
SELECT DISTINCT
  c.cooler_id,
  c.name,
  c.brand,
  c.height_mm,
  c.tdp_limit_watt
FROM pc_configurator.cpu_coolers c
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

	var out []CPUCoolerCatalogItem
	for rows.Next() {
		var x CPUCoolerCatalogItem
		if err := rows.Scan(&x.CoolerID, &x.Name, &x.Brand, &x.HeightMM, &x.TDPLimitWatt); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
