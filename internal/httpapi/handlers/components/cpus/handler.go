package cpus

import (
	domainErr "compConfigurator/internal/domain/errors"
	"compConfigurator/internal/httpapi/response"
	"compConfigurator/internal/repo"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	repo *repo.CPUsRepo
}

func New(r *repo.CPUsRepo) *Handler {
	return &Handler{repo: r}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0

	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 && n < 200 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	items, err := h.repo.List(r.Context(), limit, offset)
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

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.Fail(w, http.StatusBadRequest, domainErr.ErrInvalidInput.Error())
		return
	}

	item, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrNotFound):
			response.Fail(w, http.StatusNotFound, "cpu not found")
		default:
			response.Fail(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.JSON(w, http.StatusOK, item)
}
