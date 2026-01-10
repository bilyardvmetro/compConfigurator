package components

import (
	"net/http"
	"strconv"

	"compConfigurator/internal/httpapi/response"
	"compConfigurator/internal/repo"
)

func (h *Handler) ListPSUs(w http.ResponseWriter, r *http.Request) {
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

	filter := repo.PSUsCatalogFilter{
		Query:        q.Get("q"),
		FormFactor:   q.Get("form_factor"),
		MinPowerWatt: parseIntPtr("min_power_watt"),
		MaxLengthMM:  parseIntPtr("max_length_mm"),
		Sort:         q.Get("sort"),
		Limit:        limit,
		Offset:       offset,
	}

	items, err := h.psus.List(r.Context(), filter)
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
