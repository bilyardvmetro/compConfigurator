package components

import (
	"net/http"
	"net/url"
	"strconv"

	"compConfigurator/internal/httpapi/response"
	"compConfigurator/internal/repo"
)

type Handler struct {
	cpus         *repo.CPUsRepo
	ram          *repo.RamKitsRepo
	drv          *repo.DrivesRepo
	motherboards *repo.MotherboardsRepo
	cases        *repo.CasesRepo
	psus         *repo.PSUsRepo
	cpuCoolers   *repo.CPUCoolersRepo
	gpus         *repo.GPUsRepo
}

func New(
	cpus *repo.CPUsRepo,
	ram *repo.RamKitsRepo,
	drv *repo.DrivesRepo,
	mb *repo.MotherboardsRepo,
	cs *repo.CasesRepo,
	psu *repo.PSUsRepo,
	coolers *repo.CPUCoolersRepo,
	gpu *repo.GPUsRepo,
) *Handler {
	return &Handler{
		cpus:         cpus,
		ram:          ram,
		drv:          drv,
		motherboards: mb,
		cases:        cs,
		psus:         psu,
		cpuCoolers:   coolers,
		gpus:         gpu,
	}
}

func clampLimit(v int) int {
	if v <= 0 {
		return 50
	}
	if v > 200 {
		return 200
	}
	return v
}

func parseIntPtr(q map[string][]string, key string) *int {
	raw := ""
	if xs := q[key]; len(xs) > 0 {
		raw = xs[0]
	}
	if raw == "" {
		return nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return nil
	}
	return &n
}

func parseFloatPtr(q url.Values, key string) *float64 {
	v := q.Get(key)
	if v == "" {
		return nil
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil
	}
	return &f
}

func (h *Handler) ListRamKits(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if offset < 0 {
		offset = 0
	}
	limit = clampLimit(limit)

	filter := repo.RamKitsCatalogFilter{
		Query:         q.Get("q"),
		RamType:       q.Get("ram_type"),
		MinFreqMHz:    parseIntPtr(q, "min_freq"),
		MaxFreqMHz:    parseIntPtr(q, "max_freq"),
		MinCapacityGb: parseIntPtr(q, "min_capacity_gb"),
		MaxCapacityGb: parseIntPtr(q, "max_capacity_gb"),
		Sort:          q.Get("sort"),
		Limit:         limit,
		Offset:        offset,
	}

	items, err := h.ram.List(r.Context(), filter)
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

func (h *Handler) ListDrives(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if offset < 0 {
		offset = 0
	}
	limit = clampLimit(limit)

	filter := repo.DrivesCatalogFilter{
		Query:         q.Get("q"),
		Interface:     q.Get("interface"),
		FormFactor:    q.Get("form_factor"),
		MinCapacityGb: parseIntPtr(q, "min_capacity_gb"),
		MaxCapacityGb: parseIntPtr(q, "max_capacity_gb"),
		Sort:          q.Get("sort"),
		Limit:         limit,
		Offset:        offset,
	}

	items, err := h.drv.List(r.Context(), filter)
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
