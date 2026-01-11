package assemblies

import (
	"compConfigurator/internal/service/pricing"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	domainerr "compConfigurator/internal/domain/errors"
	"compConfigurator/internal/httpapi/middleware"
	"compConfigurator/internal/httpapi/response"
	"compConfigurator/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Handler struct {
	svc        *service.AssemblyService
	partsSvc   *service.AssemblyPartsService
	detailsSvc *service.AssemblyDetailsService
	pricingSvc *pricing.Service
}

func New(s *service.AssemblyService, p *service.AssemblyPartsService, d *service.AssemblyDetailsService, pr *pricing.Service) *Handler {
	return &Handler{svc: s, partsSvc: p, detailsSvc: d, pricingSvc: pr}
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

type updateReq struct {
	Name     string `json:"name"`
	IsPublic bool   `json:"is_public"`

	CPUId         int64 `json:"cpu_id"`
	MotherboardId int64 `json:"motherboard_id"`
	PSUId         int64 `json:"psu_id"`
	CaseId        int64 `json:"case_id"`

	GPUId    *int64 `json:"gpu_id"`    // null допустим
	CoolerId *int64 `json:"cooler_id"` // null допустим
}

type addDriveReq struct {
	MountType *string `json:"mount_type"` // можно null/omitted
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

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, domainerr.ErrUnauthorized.Error())
		return
	}

	idStr := chi.URLParam(r, "id")
	assemblyID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || assemblyID <= 0 {
		response.Fail(w, http.StatusBadRequest, domainerr.ErrInvalidInput.Error())
		return
	}

	var req updateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.Name == "" {
		req.Name = "My build"
	}
	if req.CPUId <= 0 || req.MotherboardId <= 0 || req.PSUId <= 0 || req.CaseId <= 0 {
		response.Fail(w, http.StatusBadRequest, "cpu_id, motherboard_id, psu_id, case_id are required")
		return
	}

	view, err := h.svc.Update(
		r.Context(),
		userID,
		assemblyID,
		req.Name,
		req.IsPublic,
		req.CPUId,
		req.MotherboardId,
		req.PSUId,
		req.CaseId,
		req.GPUId,
		req.CoolerId,
	)
	if err != nil {
		// Маппинг ошибок Postgres
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23514":
				// check_violation — это наша несовместимость
				response.Fail(w, http.StatusConflict, pgErr.Message)
				return
			case "23503":
				// foreign_key_violation — компонент не найден
				response.Fail(w, http.StatusBadRequest, "component reference not found")
				return
			}
		}

		if errors.Is(err, domainerr.ErrNotFound) {
			response.Fail(w, http.StatusNotFound, "assembly not found")
			return
		}

		response.Fail(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusOK, view)
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

func (h *Handler) AddRamKit(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, domainerr.ErrUnauthorized.Error())
		return
	}

	assemblyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || assemblyID <= 0 {
		response.Fail(w, http.StatusBadRequest, domainerr.ErrInvalidInput.Error())
		return
	}
	ramKitID, err := strconv.ParseInt(chi.URLParam(r, "ramKitId"), 10, 64)
	if err != nil || ramKitID <= 0 {
		response.Fail(w, http.StatusBadRequest, domainerr.ErrInvalidInput.Error())
		return
	}

	view, err := h.partsSvc.AddRamKit(r.Context(), userID, assemblyID, ramKitID)
	if err != nil {
		h.mapDBErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, view)
}

func (h *Handler) RemoveRamKit(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, domainerr.ErrUnauthorized.Error())
		return
	}

	assemblyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || assemblyID <= 0 {
		response.Fail(w, http.StatusBadRequest, domainerr.ErrInvalidInput.Error())
		return
	}
	ramKitID, err := strconv.ParseInt(chi.URLParam(r, "ramKitId"), 10, 64)
	if err != nil || ramKitID <= 0 {
		response.Fail(w, http.StatusBadRequest, domainerr.ErrInvalidInput.Error())
		return
	}

	view, err := h.partsSvc.RemoveRamKit(r.Context(), userID, assemblyID, ramKitID)
	if err != nil {
		h.mapDBErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, view)
}

func (h *Handler) AddDrive(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, domainerr.ErrUnauthorized.Error())
		return
	}

	assemblyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || assemblyID <= 0 {
		response.Fail(w, http.StatusBadRequest, domainerr.ErrInvalidInput.Error())
		return
	}
	driveID, err := strconv.ParseInt(chi.URLParam(r, "driveId"), 10, 64)
	if err != nil || driveID <= 0 {
		response.Fail(w, http.StatusBadRequest, domainerr.ErrInvalidInput.Error())
		return
	}

	var req addDriveReq
	_ = json.NewDecoder(r.Body).Decode(&req) // body optional

	view, err := h.partsSvc.AddDrive(r.Context(), userID, assemblyID, driveID, req.MountType)
	if err != nil {
		h.mapDBErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, view)
}

func (h *Handler) RemoveDrive(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, domainerr.ErrUnauthorized.Error())
		return
	}

	assemblyID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || assemblyID <= 0 {
		response.Fail(w, http.StatusBadRequest, domainerr.ErrInvalidInput.Error())
		return
	}
	driveID, err := strconv.ParseInt(chi.URLParam(r, "driveId"), 10, 64)
	if err != nil || driveID <= 0 {
		response.Fail(w, http.StatusBadRequest, domainerr.ErrInvalidInput.Error())
		return
	}

	view, err := h.partsSvc.RemoveDrive(r.Context(), userID, assemblyID, driveID)
	if err != nil {
		h.mapDBErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, view)
}

func (h *Handler) mapDBErr(w http.ResponseWriter, err error) {
	// доменные
	if errors.Is(err, domainerr.ErrNotFound) {
		response.Fail(w, http.StatusNotFound, "not found")
		return
	}

	// postgres
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23514":
			// check_violation -> несовместимость
			response.Fail(w, http.StatusConflict, pgErr.Message)
			return
		case "23503":
			// foreign_key_violation
			response.Fail(w, http.StatusBadRequest, "component reference not found")
			return
		case "23505":
			// unique violation (например, дубликат в M:N) — можно считать OK/идемпотентно
			response.Fail(w, http.StatusConflict, "already exists")
			return
		}
	}

	response.Fail(w, http.StatusInternalServerError, "internal error")
}
