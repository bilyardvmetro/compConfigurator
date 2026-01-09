package assemblies

import (
	"encoding/json"
	"net/http"
	"strconv"

	domainerr "compConfigurator/internal/domain/errors"
	"compConfigurator/internal/httpapi/middleware"
	"compConfigurator/internal/httpapi/response"
	"compConfigurator/internal/service"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *service.AssemblyService
}

func New(s *service.AssemblyService) *Handler {
	return &Handler{svc: s}
}

type createReq struct {
	Name          string `json:"name"`
	CPUId         int64  `json:"cpu_id"`
	MotherboardId int64  `json:"motherboard_id"`
	PSUId         int64  `json:"psu_id"`
	CaseId        int64  `json:"case_id"`
	GPUId         *int64 `json:"gpu_id"`
	CoolerId      *int64 `json:"cooler_id"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, domainerr.ErrUnauthorized.Error())
		return
	}

	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.Name == "" {
		req.Name = "My build"
	}
	// обязательные поля по схеме assemblies NOT NULL
	if req.CPUId <= 0 || req.MotherboardId <= 0 || req.PSUId <= 0 || req.CaseId <= 0 {
		response.Fail(w, http.StatusBadRequest, "cpu_id, motherboard_id, psu_id, case_id are required")
		return
	}

	view, err := h.svc.Create(
		r.Context(),
		userID,
		req.Name,
		req.CPUId,
		req.MotherboardId,
		req.PSUId,
		req.CaseId,
		req.GPUId,
		req.CoolerId,
	)
	if err != nil {
		// если check_assembly_compatibility кидает EXCEPTION, ты можешь вернуть 409/422 —
		// но пока на этапе 3 можно вернуть 400/500. Позже сделаем маппинг pgerr.Code -> 409.
		response.Fail(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusCreated, view)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, domainerr.ErrUnauthorized.Error())
		return
	}

	limit := 50
	offset := 0

	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	items, err := h.svc.List(r.Context(), userID, limit, offset)
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

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, domainerr.ErrUnauthorized.Error())
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.Fail(w, http.StatusBadRequest, domainerr.ErrInvalidInput.Error())
		return
	}

	view, err := h.svc.Get(r.Context(), userID, id)
	if err != nil {
		switch err {
		case domainerr.ErrNotFound:
			response.Fail(w, http.StatusNotFound, "assembly not found")
		default:
			response.Fail(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	response.JSON(w, http.StatusOK, view)
}
