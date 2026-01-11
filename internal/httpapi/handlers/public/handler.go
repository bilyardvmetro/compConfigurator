package public

import (
	"compConfigurator/internal/service"
	"errors"
	"net/http"
	"strconv"

	domainErr "compConfigurator/internal/domain/errors"
	"compConfigurator/internal/httpapi/response"
	"compConfigurator/internal/repo"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	pub        *repo.PublicAssembliesRepo
	clone      *repo.CloneRepo
	detailsSvc *service.AssemblyDetailsService
	// + тебе нужно получить текущего user_id из контекста (как у тебя сделано в auth middleware)
}

func New(pub *repo.PublicAssembliesRepo, clone *repo.CloneRepo, details *service.AssemblyDetailsService) *Handler {
	return &Handler{pub: pub, clone: clone, detailsSvc: details}
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

func (h *Handler) ListPublicAssemblies(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if offset < 0 {
		offset = 0
	}
	limit = clampLimit(limit)

	filter := repo.PublicAssembliesFilter{
		Query:  q.Get("q"),
		Sort:   q.Get("sort"),
		Limit:  limit,
		Offset: offset,
	}

	items, err := h.pub.List(r.Context(), filter)
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

func (h *Handler) GetPublicAssembly(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.Fail(w, http.StatusBadRequest, domainErr.ErrInvalidInput.Error())
		return
	}

	item, err := h.pub.GetPublicByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrNotFound):
			response.Fail(w, http.StatusNotFound, "assembly not found")
		default:
			response.Fail(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.JSON(w, http.StatusOK, item)
}

func (h *Handler) PublicDetails(w http.ResponseWriter, r *http.Request) {
	assemblyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || assemblyID <= 0 {
		response.Fail(w, http.StatusBadRequest, domainErr.ErrInvalidInput.Error())
		return
	}

	out, err := h.detailsSvc.GetPublic(r.Context(), assemblyID)
	if err != nil {
		if errors.Is(err, domainErr.ErrNotFound) {
			response.Fail(w, http.StatusNotFound, "assembly not found")
			return
		}
		response.Fail(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusOK, out)
}
