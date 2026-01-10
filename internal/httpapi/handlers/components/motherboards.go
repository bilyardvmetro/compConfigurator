package components

import (
	"net/http"
	"strconv"

	"compConfigurator/internal/httpapi/response"
	"compConfigurator/internal/repo"
)

func (h *Handler) ListMotherboards(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	parseIntPtr := func(key string) *int {
		v := q.Get(key)
		if v == "" {
			return nil
		}
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil
		}
		return &n
	}

	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if offset < 0 {
		offset = 0
	}
	limit = clampLimit(limit)

	filter := repo.MotherboardsCatalogFilter{
		Query:          q.Get("q"),
		SocketCode:     q.Get("socket_code"),
		RamType:        q.Get("ram_type"),
		FormFactorCode: q.Get("form_factor_code"),
		MinRamSlots:    parseIntPtr("min_ram_slots"),
		MinM2Slots:     parseIntPtr("min_m2_slots"),
		MinSataPorts:   parseIntPtr("min_sata_ports"),
		Sort:           q.Get("sort"),
		Limit:          limit,
		Offset:         offset,
	}

	items, err := h.motherboards.List(r.Context(), filter)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"items":  items,
		"limit":  limit,
		"offset": offset,
	})
}
