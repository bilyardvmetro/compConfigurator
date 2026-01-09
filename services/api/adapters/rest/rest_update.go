package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"comp_config.com/services/api/adapters/update"
	"comp_config.com/services/api/core"
	updatepb "comp_config.com/services/proto/update"
)

type Handler struct {
	log    *slog.Logger
	update update.Client
}

func NewClientUpdate(address string, log *slog.Logger) (*Handler, error) {
	updater, err := update.NewClient(address, log)
	if err != nil {
		return &Handler{}, err
	}
	return &Handler{
		update: *updater,
		log:    log,
	}, nil
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	pingers := map[string]core.Pinger{
		"update": h.update,
	}

	response := PingResponse{}
	response.Answer = make(map[string]string)

	for name, pinger := range pingers {
		if err := pinger.Ping(r.Context()); err != nil {
			response.Answer[name] = "unavailable"
			h.log.Error("service is not available", "service", name, "error", err)
			continue
		}
		response.Answer[name] = "ok"
	}

	h.respondJSON(w, http.StatusOK, response)
}

// GPU handlers
func (h *Handler) AddGPU(w http.ResponseWriter, r *http.Request) {
	var req updatepb.GpuRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode GPU request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddGPU(r.Context(), &req); err != nil {
		h.log.Error("Failed to add GPU", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add GPU")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "GPU added successfully"})
}

func (h *Handler) UpdateGPU(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid GPU ID")
		return
	}

	var req updatepb.GpuRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode GPU request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updateReq := &updatepb.UpdateGpuRequest{
		Id:  id,
		Gpu: &req,
	}

	if err := h.update.UpdateGPU(r.Context(), updateReq); err != nil {
		h.log.Error("Failed to update GPU", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update GPU")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "GPU updated successfully"})
}

func (h *Handler) DeleteGPU(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid GPU ID")
		return
	}

	if err := h.update.DeleteGPU(r.Context(), id); err != nil {
		h.log.Error("Failed to delete GPU", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete GPU")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "GPU deleted successfully"})
}

// CPU Cooler handlers
func (h *Handler) AddCPUCooler(w http.ResponseWriter, r *http.Request) {
	var req updatepb.CpuCoolerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode CPU cooler request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddCPUCooler(r.Context(), &req); err != nil {
		h.log.Error("Failed to add CPU cooler", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add CPU cooler")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "CPU cooler added successfully"})
}

func (h *Handler) UpdateCPUCooler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid CPU cooler ID")
		return
	}

	var req updatepb.CpuCoolerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode CPU cooler request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updateReq := &updatepb.UpdateCpuCoolerRequest{
		Id:        id,
		CpuCooler: &req,
	}

	if err := h.update.UpdateCPUCooler(r.Context(), updateReq); err != nil {
		h.log.Error("Failed to update CPU cooler", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update CPU cooler")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "CPU cooler updated successfully"})
}

func (h *Handler) DeleteCPUCooler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid CPU cooler ID")
		return
	}

	if err := h.update.DeleteCPUCooler(r.Context(), id); err != nil {
		h.log.Error("Failed to delete CPU cooler", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete CPU cooler")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "CPU cooler deleted successfully"})
}

// Case handlers
func (h *Handler) AddCase(w http.ResponseWriter, r *http.Request) {
	var req updatepb.CaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode case request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddCase(r.Context(), &req); err != nil {
		h.log.Error("Failed to add case", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add case")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Case added successfully"})
}

func (h *Handler) UpdateCase(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}

	var req updatepb.CaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode case request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updateReq := &updatepb.UpdateCaseRequest{
		Id:   id,
		Case: &req,
	}

	if err := h.update.UpdateCase(r.Context(), updateReq); err != nil {
		h.log.Error("Failed to update case", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update case")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Case updated successfully"})
}

func (h *Handler) DeleteCase(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}

	if err := h.update.DeleteCase(r.Context(), id); err != nil {
		h.log.Error("Failed to delete case", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete case")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Case deleted successfully"})
}

// PSU handlers
func (h *Handler) AddPSU(w http.ResponseWriter, r *http.Request) {
	var req updatepb.PsuRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode PSU request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddPSU(r.Context(), &req); err != nil {
		h.log.Error("Failed to add PSU", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add PSU")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "PSU added successfully"})
}

func (h *Handler) UpdatePSU(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid PSU ID")
		return
	}

	var req updatepb.PsuRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode PSU request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updateReq := &updatepb.UpdatePsuRequest{
		Id:  id,
		Psu: &req,
	}

	if err := h.update.UpdatePSU(r.Context(), updateReq); err != nil {
		h.log.Error("Failed to update PSU", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update PSU")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "PSU updated successfully"})
}

func (h *Handler) DeletePSU(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid PSU ID")
		return
	}

	if err := h.update.DeletePSU(r.Context(), id); err != nil {
		h.log.Error("Failed to delete PSU", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete PSU")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "PSU deleted successfully"})
}

// CPU Socket handlers
func (h *Handler) AddCPUSocket(w http.ResponseWriter, r *http.Request) {
	var req updatepb.CpuSocketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode CPU socket request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddCPUSocket(r.Context(), &req); err != nil {
		h.log.Error("Failed to add CPU socket", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add CPU socket")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "CPU socket added successfully"})
}

func (h *Handler) UpdateCPUSocket(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		h.respondError(w, http.StatusBadRequest, "CPU socket code is required")
		return
	}

	var req updatepb.CpuSocketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode CPU socket request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updateReq := &updatepb.UpdateCpuSocketRequest{
		SocketCode: code,
		CpuSocket:  &req,
	}

	if err := h.update.UpdateCPUSocket(r.Context(), updateReq); err != nil {
		h.log.Error("Failed to update CPU socket", "code", code, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update CPU socket")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "CPU socket updated successfully"})
}

func (h *Handler) DeleteCPUSocket(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		h.respondError(w, http.StatusBadRequest, "CPU socket code is required")
		return
	}

	if err := h.update.DeleteCPUSocket(r.Context(), code); err != nil {
		h.log.Error("Failed to delete CPU socket", "code", code, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete CPU socket")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "CPU socket deleted successfully"})
}

// Motherboard Form Factor handlers
func (h *Handler) AddMotherboardFormFactor(w http.ResponseWriter, r *http.Request) {
	var req updatepb.MotherboardFormFactorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode motherboard form factor request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddMotherboardFormFactor(r.Context(), &req); err != nil {
		h.log.Error("Failed to add motherboard form factor", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add motherboard form factor")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Motherboard form factor added successfully"})
}

func (h *Handler) UpdateMotherboardFormFactor(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		h.respondError(w, http.StatusBadRequest, "Form factor code is required")
		return
	}

	var req updatepb.MotherboardFormFactorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode motherboard form factor request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updateReq := &updatepb.UpdateMotherboardFormFactorRequest{
		FormFactorCode:        code,
		MotherboardFormFactor: &req,
	}

	if err := h.update.UpdateMotherboardFormFactor(r.Context(), updateReq); err != nil {
		h.log.Error("Failed to update motherboard form factor", "code", code, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update motherboard form factor")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Motherboard form factor updated successfully"})
}

func (h *Handler) DeleteMotherboardFormFactor(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		h.respondError(w, http.StatusBadRequest, "Form factor code is required")
		return
	}

	if err := h.update.DeleteMotherboardFormFactor(r.Context(), code); err != nil {
		h.log.Error("Failed to delete motherboard form factor", "code", code, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete motherboard form factor")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Motherboard form factor deleted successfully"})
}

// CPU handlers
func (h *Handler) AddCPU(w http.ResponseWriter, r *http.Request) {
	var req updatepb.CpuRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode CPU request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddCPU(r.Context(), &req); err != nil {
		h.log.Error("Failed to add CPU", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add CPU")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "CPU added successfully"})
}

func (h *Handler) UpdateCPU(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid CPU ID")
		return
	}

	var req updatepb.CpuRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode CPU request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updateReq := &updatepb.UpdateCpuRequest{
		Id:  id,
		Cpu: &req,
	}

	if err := h.update.UpdateCPU(r.Context(), updateReq); err != nil {
		h.log.Error("Failed to update CPU", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update CPU")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "CPU updated successfully"})
}

func (h *Handler) DeleteCPU(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid CPU ID")
		return
	}

	if err := h.update.DeleteCPU(r.Context(), id); err != nil {
		h.log.Error("Failed to delete CPU", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete CPU")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "CPU deleted successfully"})
}

// Motherboard handlers
func (h *Handler) AddMotherboard(w http.ResponseWriter, r *http.Request) {
	var req updatepb.MotherboardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode motherboard request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddMotherboard(r.Context(), &req); err != nil {
		h.log.Error("Failed to add motherboard", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add motherboard")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Motherboard added successfully"})
}

func (h *Handler) UpdateMotherboard(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid motherboard ID")
		return
	}

	var req updatepb.MotherboardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode motherboard request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updateReq := &updatepb.UpdateMotherboardRequest{
		Id:          id,
		Motherboard: &req,
	}

	if err := h.update.UpdateMotherboard(r.Context(), updateReq); err != nil {
		h.log.Error("Failed to update motherboard", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update motherboard")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Motherboard updated successfully"})
}

func (h *Handler) DeleteMotherboard(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid motherboard ID")
		return
	}

	if err := h.update.DeleteMotherboard(r.Context(), id); err != nil {
		h.log.Error("Failed to delete motherboard", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete motherboard")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Motherboard deleted successfully"})
}

// RAM Kit handlers
func (h *Handler) AddRAMKit(w http.ResponseWriter, r *http.Request) {
	var req updatepb.RAMKitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode RAM kit request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddRAMKit(r.Context(), &req); err != nil {
		h.log.Error("Failed to add RAM kit", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add RAM kit")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "RAM kit added successfully"})
}

func (h *Handler) UpdateRAMKit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid RAM kit ID")
		return
	}

	var req updatepb.RAMKitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode RAM kit request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updateReq := &updatepb.UpdateRAMKitRequest{
		Id:     id,
		RamKit: &req,
	}

	if err := h.update.UpdateRAMKit(r.Context(), updateReq); err != nil {
		h.log.Error("Failed to update RAM kit", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update RAM kit")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "RAM kit updated successfully"})
}

func (h *Handler) DeleteRAMKit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid RAM kit ID")
		return
	}

	if err := h.update.DeleteRAMKit(r.Context(), id); err != nil {
		h.log.Error("Failed to delete RAM kit", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete RAM kit")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "RAM kit deleted successfully"})
}

// Storage Drive handlers
func (h *Handler) AddStorageDrive(w http.ResponseWriter, r *http.Request) {
	var req updatepb.StorageDriveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode storage drive request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddStorageDrive(r.Context(), &req); err != nil {
		h.log.Error("Failed to add storage drive", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add storage drive")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Storage drive added successfully"})
}

func (h *Handler) UpdateStorageDrive(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid storage drive ID")
		return
	}

	var req updatepb.StorageDriveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode storage drive request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updateReq := &updatepb.UpdateStorageDriveRequest{
		Id:           id,
		StorageDrive: &req,
	}

	if err := h.update.UpdateStorageDrive(r.Context(), updateReq); err != nil {
		h.log.Error("Failed to update storage drive", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update storage drive")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Storage drive updated successfully"})
}

func (h *Handler) DeleteStorageDrive(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid storage drive ID")
		return
	}

	if err := h.update.DeleteStorageDrive(r.Context(), id); err != nil {
		h.log.Error("Failed to delete storage drive", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete storage drive")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Storage drive deleted successfully"})
}

// User handlers
func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	var req updatepb.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode user request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddUser(r.Context(), &req); err != nil {
		h.log.Error("Failed to add user", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add user")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "User added successfully"})
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req updatepb.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode user request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updateReq := &updatepb.UpdateUserRequest{
		Id:   id,
		User: &req,
	}

	if err := h.update.UpdateUser(r.Context(), updateReq); err != nil {
		h.log.Error("Failed to update user", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.update.DeleteUser(r.Context(), id); err != nil {
		h.log.Error("Failed to delete user", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// Shop handlers
func (h *Handler) AddShop(w http.ResponseWriter, r *http.Request) {
	var req updatepb.ShopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode shop request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddShop(r.Context(), &req); err != nil {
		h.log.Error("Failed to add shop", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add shop")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Shop added successfully"})
}

func (h *Handler) UpdateShop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid shop ID")
		return
	}

	var req updatepb.ShopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode shop request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updateReq := &updatepb.UpdateShopRequest{
		Id:   id,
		Shop: &req,
	}

	if err := h.update.UpdateShop(r.Context(), updateReq); err != nil {
		h.log.Error("Failed to update shop", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update shop")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Shop updated successfully"})
}

func (h *Handler) DeleteShop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid shop ID")
		return
	}

	if err := h.update.DeleteShop(r.Context(), id); err != nil {
		h.log.Error("Failed to delete shop", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete shop")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Shop deleted successfully"})
}

// Product Offer handlers
func (h *Handler) AddProductOffer(w http.ResponseWriter, r *http.Request) {
	var req updatepb.ProductOfferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode product offer request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddProductOffer(r.Context(), &req); err != nil {
		h.log.Error("Failed to add product offer", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add product offer")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Product offer added successfully"})
}

func (h *Handler) UpdateProductOffer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid product offer ID")
		return
	}

	var req updatepb.ProductOfferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode product offer request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updateReq := &updatepb.UpdateProductOfferRequest{
		Id:           id,
		ProductOffer: &req,
	}

	if err := h.update.UpdateProductOffer(r.Context(), updateReq); err != nil {
		h.log.Error("Failed to update product offer", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update product offer")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Product offer updated successfully"})
}

func (h *Handler) DeleteProductOffer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid product offer ID")
		return
	}

	if err := h.update.DeleteProductOffer(r.Context(), id); err != nil {
		h.log.Error("Failed to delete product offer", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete product offer")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Product offer deleted successfully"})
}

// Assembly handlers
func (h *Handler) AddAssembly(w http.ResponseWriter, r *http.Request) {
	var req updatepb.AssemblyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode assembly request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddAssembly(r.Context(), &req); err != nil {
		h.log.Error("Failed to add assembly", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add assembly")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Assembly added successfully"})
}

func (h *Handler) UpdateAssembly(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid assembly ID")
		return
	}

	var req updatepb.AssemblyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode assembly request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updateReq := &updatepb.UpdateAssemblyRequest{
		Id:       id,
		Assembly: &req,
	}

	if err := h.update.UpdateAssembly(r.Context(), updateReq); err != nil {
		h.log.Error("Failed to update assembly", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to update assembly")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Assembly updated successfully"})
}

func (h *Handler) DeleteAssembly(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid assembly ID")
		return
	}

	if err := h.update.DeleteAssembly(r.Context(), id); err != nil {
		h.log.Error("Failed to delete assembly", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete assembly")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Assembly deleted successfully"})
}

// Cooler Socket handlers (many-to-many)
func (h *Handler) AddCoolerSocket(w http.ResponseWriter, r *http.Request) {
	var req updatepb.CoolerSocketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode cooler socket request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddCoolerSocket(r.Context(), &req); err != nil {
		h.log.Error("Failed to add cooler socket", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add cooler socket")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Cooler socket added successfully"})
}

func (h *Handler) DeleteCoolerSocket(w http.ResponseWriter, r *http.Request) {
	var req updatepb.DeleteCoolerSocketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode delete cooler socket request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.DeleteCoolerSocket(r.Context(), &req); err != nil {
		h.log.Error("Failed to delete cooler socket", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete cooler socket")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Cooler socket deleted successfully"})
}

func (h *Handler) DeleteCoolerSocketsByCooler(w http.ResponseWriter, r *http.Request) {
	coolerID, err := strconv.ParseInt(r.PathValue("cooler_id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid cooler ID")
		return
	}

	if err := h.update.DeleteCoolerSocketsByCooler(r.Context(), coolerID); err != nil {
		h.log.Error("Failed to delete cooler sockets by cooler", "cooler_id", coolerID, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete cooler sockets")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Cooler sockets deleted successfully"})
}

// Case Form Factor Support handlers (many-to-many)
func (h *Handler) AddCaseFormFactorSupport(w http.ResponseWriter, r *http.Request) {
	var req updatepb.CaseFormFactorSupportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode case form factor support request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddCaseFormFactorSupport(r.Context(), &req); err != nil {
		h.log.Error("Failed to add case form factor support", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add case form factor support")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Case form factor support added successfully"})
}

func (h *Handler) DeleteCaseFormFactorSupport(w http.ResponseWriter, r *http.Request) {
	var req updatepb.DeleteCaseFormFactorSupportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode delete case form factor support request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.DeleteCaseFormFactorSupport(r.Context(), &req); err != nil {
		h.log.Error("Failed to delete case form factor support", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete case form factor support")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Case form factor support deleted successfully"})
}

func (h *Handler) DeleteCaseFormFactorSupportsByCase(w http.ResponseWriter, r *http.Request) {
	caseID, err := strconv.ParseInt(r.PathValue("case_id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}

	if err := h.update.DeleteCaseFormFactorSupportsByCase(r.Context(), caseID); err != nil {
		h.log.Error("Failed to delete case form factor supports by case", "case_id", caseID, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete case form factor supports")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Case form factor supports deleted successfully"})
}

// Assembly RAM Kit handlers (many-to-many)
func (h *Handler) AddAssemblyRAMKit(w http.ResponseWriter, r *http.Request) {
	var req updatepb.AssemblyRAMKitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode assembly RAM kit request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddAssemblyRAMKit(r.Context(), &req); err != nil {
		h.log.Error("Failed to add assembly RAM kit", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add assembly RAM kit")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Assembly RAM kit added successfully"})
}

func (h *Handler) DeleteAssemblyRAMKit(w http.ResponseWriter, r *http.Request) {
	var req updatepb.DeleteAssemblyRAMKitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode delete assembly RAM kit request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.DeleteAssemblyRAMKit(r.Context(), &req); err != nil {
		h.log.Error("Failed to delete assembly RAM kit", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete assembly RAM kit")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Assembly RAM kit deleted successfully"})
}

func (h *Handler) DeleteAssemblyRAMKitsByAssembly(w http.ResponseWriter, r *http.Request) {
	assemblyID, err := strconv.ParseInt(r.PathValue("assembly_id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid assembly ID")
		return
	}

	if err := h.update.DeleteAssemblyRAMKitsByAssembly(r.Context(), assemblyID); err != nil {
		h.log.Error("Failed to delete assembly RAM kits by assembly", "assembly_id", assemblyID, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete assembly RAM kits")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Assembly RAM kits deleted successfully"})
}

// Assembly Drive handlers (many-to-many)
func (h *Handler) AddAssemblyDrive(w http.ResponseWriter, r *http.Request) {
	var req updatepb.AssemblyDriveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode assembly drive request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.AddAssemblyDrive(r.Context(), &req); err != nil {
		h.log.Error("Failed to add assembly drive", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to add assembly drive")
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"message": "Assembly drive added successfully"})
}

func (h *Handler) DeleteAssemblyDrive(w http.ResponseWriter, r *http.Request) {
	var req updatepb.DeleteAssemblyDriveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode delete assembly drive request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.update.DeleteAssemblyDrive(r.Context(), &req); err != nil {
		h.log.Error("Failed to delete assembly drive", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete assembly drive")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Assembly drive deleted successfully"})
}

func (h *Handler) DeleteAssemblyDrivesByAssembly(w http.ResponseWriter, r *http.Request) {
	assemblyID, err := strconv.ParseInt(r.PathValue("assembly_id"), 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid assembly ID")
		return
	}

	if err := h.update.DeleteAssemblyDrivesByAssembly(r.Context(), assemblyID); err != nil {
		h.log.Error("Failed to delete assembly drives by assembly", "assembly_id", assemblyID, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to delete assembly drives")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"message": "Assembly drives deleted successfully"})
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("Failed to encode response", "error", err)
	}
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}

type PingResponse struct {
	Answer map[string]string `json:"replies"`
}

func NewPingHandler(log *slog.Logger, pingers map[string]core.Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := PingResponse{}
		response.Answer = make(map[string]string)

		for name, pinger := range pingers {
			if err := pinger.Ping(r.Context()); err != nil {
				response.Answer[name] = "unavailable"
				log.Error("service is not available", "service", name)
				continue
			}
			response.Answer[name] = "ok"
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("Response cannot be encoded", "error", err)
		}

	}
}
