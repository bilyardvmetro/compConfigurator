package repo

import (
	"context"
	"fmt"
	"strings"

	"compConfigurator/internal/db"
)

type MotherboardCatalogItem struct {
	MotherboardID  int64  `json:"motherboard_id"`
	Name           string `json:"name"`
	SocketCode     string `json:"socket_code"`
	RamType        string `json:"ram_type"`
	FormFactorCode string `json:"form_factor_code"`
	RamSlots       int    `json:"ram_slots"`
	M2SlotsCount   int    `json:"m2_slots_count"`
	SataPortsCount int    `json:"sata_ports_count"`
}

type MotherboardsCatalogFilter struct {
	Query          string
	SocketCode     string
	RamType        string
	FormFactorCode string
	MinRamSlots    *int
	MinM2Slots     *int
	MinSataPorts   *int

	Sort   string
	Limit  int
	Offset int
}

type MotherboardsRepo struct {
	pool *db.Pool
}

func NewMotherboardsRepo(pool *db.Pool) *MotherboardsRepo {
	return &MotherboardsRepo{pool: pool}
}

func (r *MotherboardsRepo) List(ctx context.Context, f MotherboardsCatalogFilter) ([]MotherboardCatalogItem, error) {
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
	if strings.TrimSpace(f.SocketCode) != "" {
		where = append(where, "socket_code = "+arg(f.SocketCode))
	}
	if strings.TrimSpace(f.RamType) != "" {
		where = append(where, "ram_type = "+arg(f.RamType))
	}
	if strings.TrimSpace(f.FormFactorCode) != "" {
		where = append(where, "form_factor_code = "+arg(f.FormFactorCode))
	}
	if f.MinRamSlots != nil {
		where = append(where, "ram_slots >= "+arg(*f.MinRamSlots))
	}
	if f.MinM2Slots != nil {
		where = append(where, "m2_slots_count >= "+arg(*f.MinM2Slots))
	}
	if f.MinSataPorts != nil {
		where = append(where, "sata_ports_count >= "+arg(*f.MinSataPorts))
	}

	sql := `
SELECT
  motherboard_id,
  name,
  socket_code,
  ram_type,
  form_factor_code,
  ram_slots,
  m2_slots_count,
  sata_ports_count
FROM pc_configurator.motherboards
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

	var out []MotherboardCatalogItem
	for rows.Next() {
		var x MotherboardCatalogItem
		if err := rows.Scan(
			&x.MotherboardID,
			&x.Name,
			&x.SocketCode,
			&x.RamType,
			&x.FormFactorCode,
			&x.RamSlots,
			&x.M2SlotsCount,
			&x.SataPortsCount,
		); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
