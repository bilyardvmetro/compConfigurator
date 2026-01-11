package repo

import (
	"context"
	"fmt"
	"strings"

	"compConfigurator/internal/db"
)

type Offer struct {
	OfferID  int64   `json:"offer_id"`
	ShopID   int64   `json:"shop_id"`
	ShopName string  `json:"shop_name"`
	ShopURL  *string `json:"shop_url,omitempty"`

	ComponentType string `json:"component_type"`
	ComponentID   int64  `json:"component_id"`

	// NUMERIC(12,2) будем отдавать как:
	PriceCents int64 `json:"price_cents"` // цена * 100
	Available  bool  `json:"available"`
}

type OffersFilter struct {
	ComponentType string
	ComponentID   int64

	OnlyAvailable bool
	Sort          string // price_asc | price_desc
	Limit         int
	Offset        int
}

type OffersRepo struct {
	pool *db.Pool
}

func NewOffersRepo(pool *db.Pool) *OffersRepo {
	return &OffersRepo{pool: pool}
}

func (r *OffersRepo) List(ctx context.Context, f OffersFilter) ([]Offer, error) {
	if f.Limit <= 0 || f.Limit > 200 {
		f.Limit = 50
	}
	if f.Offset < 0 {
		f.Offset = 0
	}

	orderBy := "po.price ASC"
	if f.Sort == "price_desc" {
		orderBy = "po.price DESC"
	}

	var where []string
	var args []any
	arg := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf("$%d", len(args))
	}

	where = append(where, "po.component_type = "+arg(f.ComponentType))
	where = append(where, "po.component_id = "+arg(f.ComponentID))

	if f.OnlyAvailable {
		where = append(where, "po.available = TRUE")
	}

	// ВАЖНО: price NUMERIC(12,2) -> копейки:
	// CAST(price * 100 AS BIGINT)
	sql := `
SELECT
  po.offer_id,
  po.shop_id,
  s.name,
  s.url,
  po.component_type,
  po.component_id,
  CAST(po.price * 100 AS BIGINT) AS price_cents,
  po.available
FROM pc_configurator.product_offers po
JOIN pc_configurator.shops s ON s.shop_id = po.shop_id
`
	sql += "WHERE " + strings.Join(where, " AND ") + "\n"
	sql += "ORDER BY " + orderBy + "\n"
	sql += fmt.Sprintf("LIMIT %s OFFSET %s", arg(f.Limit), arg(f.Offset))

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Offer
	for rows.Next() {
		var x Offer
		if err := rows.Scan(
			&x.OfferID,
			&x.ShopID,
			&x.ShopName,
			&x.ShopURL,
			&x.ComponentType,
			&x.ComponentID,
			&x.PriceCents,
			&x.Available,
		); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *OffersRepo) Best(ctx context.Context, componentType string, componentID int64, onlyAvailable bool) (*Offer, error) {
	where := `
WHERE po.component_type = $1 AND po.component_id = $2
`
	if onlyAvailable {
		where += " AND po.available = TRUE\n"
	}

	sql := `
SELECT
  po.offer_id,
  po.shop_id,
  s.name,
  s.url,
  po.component_type,
  po.component_id,
  CAST(po.price * 100 AS BIGINT) AS price_cents,
  po.available
FROM pc_configurator.product_offers po
JOIN pc_configurator.shops s ON s.shop_id = po.shop_id
` + where + `
ORDER BY po.price ASC
LIMIT 1
`

	var x Offer
	err := r.pool.QueryRow(ctx, sql, componentType, componentID).Scan(
		&x.OfferID,
		&x.ShopID,
		&x.ShopName,
		&x.ShopURL,
		&x.ComponentType,
		&x.ComponentID,
		&x.PriceCents,
		&x.Available,
	)
	if err != nil {
		// нет офферов → nil, nil (удобно для pricing)
		return nil, nil
	}
	return &x, nil
}
