package offers

import (
	"net/http"
	"strconv"

	"compConfigurator/internal/httpapi/response"
	"compConfigurator/internal/repo"
)

type Handler struct {
	offers *repo.OffersRepo
}

func New(offers *repo.OffersRepo) *Handler {
	return &Handler{offers: offers}
}

func clampLimit(n int) int {
	if n <= 0 {
		return 50
	}
	if n > 200 {
		return 200
	}
	return n
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	componentType := q.Get("component_type")
	componentID, _ := strconv.ParseInt(q.Get("component_id"), 10, 64)
	if componentType == "" || componentID <= 0 {
		response.Fail(w, http.StatusBadRequest, "component_type and component_id are required")
		return
	}

	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	limit = clampLimit(limit)
	if offset < 0 {
		offset = 0
	}

	onlyAvailable := q.Get("only_available") == "true"

	items, err := h.offers.List(r.Context(), repo.OffersFilter{
		ComponentType: componentType,
		ComponentID:   componentID,
		OnlyAvailable: onlyAvailable,
		Sort:          q.Get("sort"),
		Limit:         limit,
		Offset:        offset,
	})
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

func (h *Handler) Best(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	componentType := q.Get("component_type")
	componentID, _ := strconv.ParseInt(q.Get("component_id"), 10, 64)
	if componentType == "" || componentID <= 0 {
		response.Fail(w, http.StatusBadRequest, "component_type and component_id are required")
		return
	}

	onlyAvailable := q.Get("only_available") == "true"

	item, err := h.offers.Best(r.Context(), componentType, componentID, onlyAvailable)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"offer": item, // может быть null
	})
}
