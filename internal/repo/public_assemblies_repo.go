package repo

import (
	"context"
	"fmt"
	"strings"

	"compConfigurator/internal/db"
	domainErr "compConfigurator/internal/domain/errors"

	"errors"

	"github.com/jackc/pgx/v5"
)

type PublicAssemblyItem struct {
	AssemblyID int64   `json:"assembly_id"`
	Name       string  `json:"name"`
	UserID     int64   `json:"user_id"`
	AuthorNick *string `json:"author_nickname,omitempty"`

	// если у тебя есть кеш цены:
	TotalPrice *int `json:"total_price_cached,omitempty"`

	// если есть created_at — будет работать sort newest/oldest
	CreatedAt *string `json:"created_at,omitempty"`
}

type PublicAssembliesFilter struct {
	Query  string
	Sort   string // newest|oldest|name_asc|name_desc
	Limit  int
	Offset int
}

type PublicAssembliesRepo struct {
	pool *db.Pool
}

func NewPublicAssembliesRepo(pool *db.Pool) *PublicAssembliesRepo {
	return &PublicAssembliesRepo{pool: pool}
}

// ВАЖНО: в SQL ниже я делаю LEFT JOIN на users для nickname.
// Если у тебя поле nickname в users называется иначе — поправь.
func (r *PublicAssembliesRepo) List(ctx context.Context, f PublicAssembliesFilter) ([]PublicAssemblyItem, error) {
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 50
	}
	if f.Offset < 0 {
		f.Offset = 0
	}

	// sorting
	orderBy := "a.assembly_id DESC"
	switch f.Sort {
	case "oldest":
		orderBy = "a.assembly_id ASC"
	case "name_asc":
		orderBy = "a.name ASC"
	case "name_desc":
		orderBy = "a.name DESC"
	case "newest":
		fallthrough
	default:
		// newest по id (если есть created_at — лучше по created_at DESC)
		orderBy = "a.assembly_id DESC"
	}

	var where []string
	var args []any
	arg := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf("$%d", len(args))
	}

	where = append(where, "a.is_public = TRUE")

	if strings.TrimSpace(f.Query) != "" {
		where = append(where, "a.name ILIKE "+arg("%"+strings.TrimSpace(f.Query)+"%"))
	}

	// Если у тебя нет total_price_cached/created_at — просто убери из SELECT/Scan.
	sql := `
SELECT
  a.assembly_id,
  a.name,
  a.user_id,
  u.nickname,
  a.total_price_cached
FROM pc_configurator.assemblies a
JOIN pc_configurator.users u ON u.user_id = a.user_id
`
	sql += "WHERE " + strings.Join(where, " AND ") + "\n"
	sql += "ORDER BY " + orderBy + "\n"
	sql += fmt.Sprintf("LIMIT %s OFFSET %s", arg(f.Limit), arg(f.Offset))

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []PublicAssemblyItem
	for rows.Next() {
		var x PublicAssemblyItem
		if err := rows.Scan(&x.AssemblyID, &x.Name, &x.UserID, &x.AuthorNick, &x.TotalPrice); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

type PublicAssembly struct {
	AssemblyID int64  `json:"assembly_id"`
	Name       string `json:"name"`
	UserID     int64  `json:"user_id"`
	IsPublic   bool   `json:"is_public"`

	CPUId         *int64 `json:"cpu_id,omitempty"`
	MotherboardId *int64 `json:"motherboard_id,omitempty"`
	PSUId         *int64 `json:"psu_id,omitempty"`
	CaseId        *int64 `json:"case_id,omitempty"`
	GPUId         *int64 `json:"gpu_id,omitempty"`
	CoolerId      *int64 `json:"cooler_id,omitempty"`

	AuthorNick *string `json:"author_nickname,omitempty"`
	TotalPrice *int    `json:"total_price_cached,omitempty"`
}

func (r *PublicAssembliesRepo) GetPublicByID(ctx context.Context, id int64) (PublicAssembly, error) {
	var a PublicAssembly
	err := r.pool.QueryRow(ctx, `
SELECT
  a.assembly_id, a.name, a.user_id, a.is_public,
  a.cpu_id, a.motherboard_id, a.psu_id, a.case_id, a.gpu_id, a.cooler_id,
  u.nickname,
  a.total_price_cached
FROM pc_configurator.assemblies a
JOIN pc_configurator.users u ON u.user_id = a.user_id
WHERE a.assembly_id = $1 AND a.is_public = TRUE
`, id).Scan(
		&a.AssemblyID, &a.Name, &a.UserID, &a.IsPublic,
		&a.CPUId, &a.MotherboardId, &a.PSUId, &a.CaseId, &a.GPUId, &a.CoolerId,
		&a.AuthorNick,
		&a.TotalPrice,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return PublicAssembly{}, domainErr.ErrNotFound
		}
		return PublicAssembly{}, err
	}
	return a, nil
}
