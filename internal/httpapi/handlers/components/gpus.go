package components

import (
	"net/http"
	"strconv"

	"compConfigurator/internal/httpapi/response"
	"compConfigurator/internal/repo"
)

func (h *Handler) ListGPUs(w http.ResponseWriter, r *http.Request) {
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

	filter := repo.GPUsCatalogFilter{
		Query: q.Get("q"),

		MinLengthMM: parseIntPtr("min_length_mm"),
		MaxLengthMM: parseIntPtr("max_length_mm"),

		MinWidthSlots: parseFloatPtr(q, "min_width_slots"),
		MaxWidthSlots: parseFloatPtr(q, "max_width_slots"),

		MinTDPWatt: parseIntPtr("min_tdp_watt"),
		MaxTDPWatt: parseIntPtr("max_tdp_watt"),

		Sort:   q.Get("sort"),
		Limit:  limit,
		Offset: offset,
	}

	items, err := h.gpus.List(r.Context(), filter)
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
