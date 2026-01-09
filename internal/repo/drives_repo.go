package repo

import (
	"context"
	"fmt"
	"strings"

	"compConfigurator/internal/db"
)

type DriveCatalogItem struct {
	DriveID    int64   `json:"drive_id"`
	Name       string  `json:"name"`
	Interface  *string `json:"interface,omitempty"`   // NVMe / SATA
	FormFactor *string `json:"form_factor,omitempty"` // 2.5 / 3.5 / 2280...
	CapacityGb *int    `json:"capacity_gb,omitempty"` // <-- если у тебя другое имя, поправь в SELECT
}

type DrivesCatalogFilter struct {
	Query         string
	Interface     string
	FormFactor    string
	MinCapacityGb *int
	MaxCapacityGb *int

	Sort   string
	Limit  int
	Offset int
}

type DrivesRepo struct {
	pool *db.Pool
}

func NewDrivesRepo(pool *db.Pool) *DrivesRepo {
	return &DrivesRepo{pool: pool}
}

func (r *DrivesRepo) List(ctx context.Context, f DrivesCatalogFilter) ([]DriveCatalogItem, error) {
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
	if strings.TrimSpace(f.Interface) != "" {
		where = append(where, "interface = "+arg(strings.TrimSpace(f.Interface)))
	}
	if strings.TrimSpace(f.FormFactor) != "" {
		where = append(where, "form_factor = "+arg(strings.TrimSpace(f.FormFactor)))
	}
	if f.MinCapacityGb != nil {
		where = append(where, "capacity_gb >= "+arg(*f.MinCapacityGb))
	}
	if f.MaxCapacityGb != nil {
		where = append(where, "capacity_gb <= "+arg(*f.MaxCapacityGb))
	}

	sql := `
SELECT drive_id, name, interface, form_factor, capacity_gb
FROM pc_configurator.storage_drives
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

	var out []DriveCatalogItem
	for rows.Next() {
		var x DriveCatalogItem
		if err := rows.Scan(&x.DriveID, &x.Name, &x.Interface, &x.FormFactor, &x.CapacityGb); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
