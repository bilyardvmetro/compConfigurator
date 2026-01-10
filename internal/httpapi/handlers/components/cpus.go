package components

import (
	"net/http"
	"strconv"

	"compConfigurator/internal/httpapi/response"
	"compConfigurator/internal/repo"
)

func (h *Handler) ListCPUs(w http.ResponseWriter, r *http.Request) {
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

	parseBoolPtr := func(key string) *bool {
		v := q.Get(key)
		if v == "" {
			return nil
		}
		b, err := strconv.ParseBool(v) // true/false/1/0
		if err != nil {
			return nil
		}
		return &b
	}

	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if offset < 0 {
		offset = 0
	}
	limit = clampLimit(limit)

	filter := repo.CPUsCatalogFilter{
		Query:      q.Get("q"),
		SocketCode: q.Get("socket_code"),
		RamType:    q.Get("ram_type"),
		HasIGPU:    parseBoolPtr("has_igpu"),

		MinRamFreqMHz: parseIntPtr("min_ram_freq_mhz"),
		MaxTdpWatt:    parseIntPtr("max_tdp_watt"),
		MinPcieVer:    parseIntPtr("min_pcie_version"),

		Sort:   q.Get("sort"),
		Limit:  limit,
		Offset: offset,
	}

	items, err := h.cpus.List(r.Context(), filter)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"items":  items,
		"limit":  limit,
		"offset": offset,
	})
}
