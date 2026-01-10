package components

import (
	"net/http"
	"strconv"

	"compConfigurator/internal/httpapi/response"
	"compConfigurator/internal/repo"
)

func (h *Handler) ListCases(w http.ResponseWriter, r *http.Request) {
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

	filter := repo.CasesCatalogFilter{
		Query:                q.Get("q"),
		FormFactorCode:       q.Get("form_factor_code"),
		PSUFormFactor:        q.Get("psu_form_factor"),
		MinMaxGPULengthMM:    parseIntPtr("min_max_gpu_length_mm"),
		MinMaxGPUWidthSlots:  parseFloatPtr(q, "min_max_gpu_width_slots"),
		MinMaxCoolerHeightMM: parseIntPtr("min_max_cooler_height_mm"),
		MinDriveBays35:       parseIntPtr("min_drive_bays_3_5"),
		MinDriveBays25:       parseIntPtr("min_drive_bays_2_5"),
		Sort:                 q.Get("sort"),
		Limit:                limit,
		Offset:               offset,
	}

	items, err := h.cases.List(r.Context(), filter)
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
