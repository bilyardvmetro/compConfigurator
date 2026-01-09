package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"comp_config.com/services/api/adapters/search"
	searchpb "comp_config.com/services/proto/search"
)

type SearchHandler struct {
	log    *slog.Logger
	search *search.Client
}

func NewSearchHandler(address string, log *slog.Logger) (*SearchHandler, error) {
	searcher, err := search.NewClient(address, log)
	if err != nil {
		return &SearchHandler{}, err
	}
	return &SearchHandler{
		search: searcher,
		log:    log,
	}, nil
}

// Get methods handlers
func (h *SearchHandler) GetGPU(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid GPU ID")
		return
	}

	gpu, err := h.search.GetGPU(r.Context(), int32(id))
	if err != nil {
		h.log.Error("Failed to get GPU", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get GPU")
		return
	}

	h.respondJSON(w, http.StatusOK, gpu)
}

func (h *SearchHandler) GetCPU(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid CPU ID")
		return
	}

	cpu, err := h.search.GetCPU(r.Context(), int32(id))
	if err != nil {
		h.log.Error("Failed to get CPU", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get CPU")
		return
	}

	h.respondJSON(w, http.StatusOK, cpu)
}

func (h *SearchHandler) GetMotherboard(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid motherboard ID")
		return
	}

	mb, err := h.search.GetMotherboard(r.Context(), int32(id))
	if err != nil {
		h.log.Error("Failed to get motherboard", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get motherboard")
		return
	}

	h.respondJSON(w, http.StatusOK, mb)
}

func (h *SearchHandler) GetRAMKit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid RAM kit ID")
		return
	}

	ram, err := h.search.GetRAMKit(r.Context(), int32(id))
	if err != nil {
		h.log.Error("Failed to get RAM kit", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get RAM kit")
		return
	}

	h.respondJSON(w, http.StatusOK, ram)
}

func (h *SearchHandler) GetPSU(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid PSU ID")
		return
	}

	psu, err := h.search.GetPSU(r.Context(), int32(id))
	if err != nil {
		h.log.Error("Failed to get PSU", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get PSU")
		return
	}

	h.respondJSON(w, http.StatusOK, psu)
}

func (h *SearchHandler) GetCase(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}

	case_, err := h.search.GetCase(r.Context(), int32(id))
	if err != nil {
		h.log.Error("Failed to get case", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get case")
		return
	}

	h.respondJSON(w, http.StatusOK, case_)
}

func (h *SearchHandler) GetCpuCooler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid CPU cooler ID")
		return
	}

	cooler, err := h.search.GetCpuCooler(r.Context(), int32(id))
	if err != nil {
		h.log.Error("Failed to get CPU cooler", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get CPU cooler")
		return
	}

	h.respondJSON(w, http.StatusOK, cooler)
}

func (h *SearchHandler) GetStorageDrive(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid storage drive ID")
		return
	}

	drive, err := h.search.GetStorageDrive(r.Context(), int32(id))
	if err != nil {
		h.log.Error("Failed to get storage drive", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get storage drive")
		return
	}

	h.respondJSON(w, http.StatusOK, drive)
}

func (h *SearchHandler) GetCpuSocket(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		h.respondError(w, http.StatusBadRequest, "CPU socket code is required")
		return
	}

	socket, err := h.search.GetCpuSocket(r.Context(), code)
	if err != nil {
		h.log.Error("Failed to get CPU socket", "code", code, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get CPU socket")
		return
	}

	h.respondJSON(w, http.StatusOK, socket)
}

func (h *SearchHandler) GetMotherboardFormFactor(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		h.respondError(w, http.StatusBadRequest, "Form factor code is required")
		return
	}

	ff, err := h.search.GetMotherboardFormFactor(r.Context(), code)
	if err != nil {
		h.log.Error("Failed to get motherboard form factor", "code", code, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get motherboard form factor")
		return
	}

	h.respondJSON(w, http.StatusOK, ff)
}

func (h *SearchHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.search.GetUser(r.Context(), int32(id))
	if err != nil {
		h.log.Error("Failed to get user", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	h.respondJSON(w, http.StatusOK, user)
}

func (h *SearchHandler) GetShop(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid shop ID")
		return
	}

	shop, err := h.search.GetShop(r.Context(), int32(id))
	if err != nil {
		h.log.Error("Failed to get shop", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get shop")
		return
	}

	h.respondJSON(w, http.StatusOK, shop)
}

func (h *SearchHandler) GetProductOffer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid product offer ID")
		return
	}

	offer, err := h.search.GetProductOffer(r.Context(), int32(id))
	if err != nil {
		h.log.Error("Failed to get product offer", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get product offer")
		return
	}

	h.respondJSON(w, http.StatusOK, offer)
}

func (h *SearchHandler) GetAssembly(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid assembly ID")
		return
	}

	assembly, err := h.search.GetAssembly(r.Context(), int32(id))
	if err != nil {
		h.log.Error("Failed to get assembly", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get assembly")
		return
	}

	h.respondJSON(w, http.StatusOK, assembly)
}

// List methods handlers
func (h *SearchHandler) ListGPUs(w http.ResponseWriter, r *http.Request) {
	req, err := h.parseListRequest(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.search.ListGPUs(r.Context(), req)
	if err != nil {
		h.log.Error("Failed to list GPUs", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to list GPUs")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) ListCPUs(w http.ResponseWriter, r *http.Request) {
	req, err := h.parseListRequest(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.search.ListCPUs(r.Context(), req)
	if err != nil {
		h.log.Error("Failed to list CPUs", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to list CPUs")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) ListMotherboards(w http.ResponseWriter, r *http.Request) {
	req, err := h.parseListRequest(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.search.ListMotherboards(r.Context(), req)
	if err != nil {
		h.log.Error("Failed to list motherboards", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to list motherboards")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) ListRAMKits(w http.ResponseWriter, r *http.Request) {
	req, err := h.parseListRequest(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.search.ListRAMKits(r.Context(), req)
	if err != nil {
		h.log.Error("Failed to list RAM kits", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to list RAM kits")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) ListPSUs(w http.ResponseWriter, r *http.Request) {
	req, err := h.parseListRequest(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.search.ListPSUs(r.Context(), req)
	if err != nil {
		h.log.Error("Failed to list PSUs", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to list PSUs")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) ListCases(w http.ResponseWriter, r *http.Request) {
	req, err := h.parseListRequest(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.search.ListCases(r.Context(), req)
	if err != nil {
		h.log.Error("Failed to list cases", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to list cases")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) ListCpuCoolers(w http.ResponseWriter, r *http.Request) {
	req, err := h.parseListRequest(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.search.ListCpuCoolers(r.Context(), req)
	if err != nil {
		h.log.Error("Failed to list CPU coolers", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to list CPU coolers")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) ListStorageDrives(w http.ResponseWriter, r *http.Request) {
	req, err := h.parseListRequest(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.search.ListStorageDrives(r.Context(), req)
	if err != nil {
		h.log.Error("Failed to list storage drives", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to list storage drives")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) ListCpuSockets(w http.ResponseWriter, r *http.Request) {
	req, err := h.parseListRequest(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.search.ListCpuSockets(r.Context(), req)
	if err != nil {
		h.log.Error("Failed to list CPU sockets", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to list CPU sockets")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) ListMotherboardFormFactors(w http.ResponseWriter, r *http.Request) {
	req, err := h.parseListRequest(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.search.ListMotherboardFormFactors(r.Context(), req)
	if err != nil {
		h.log.Error("Failed to list motherboard form factors", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to list motherboard form factors")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	req, err := h.parseListRequest(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.search.ListUsers(r.Context(), req)
	if err != nil {
		h.log.Error("Failed to list users", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to list users")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) ListShops(w http.ResponseWriter, r *http.Request) {
	req, err := h.parseListRequest(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.search.ListShops(r.Context(), req)
	if err != nil {
		h.log.Error("Failed to list shops", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to list shops")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) ListProductOffers(w http.ResponseWriter, r *http.Request) {
	req, err := h.parseListRequest(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.search.ListProductOffers(r.Context(), req)
	if err != nil {
		h.log.Error("Failed to list product offers", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to list product offers")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) ListAssemblies(w http.ResponseWriter, r *http.Request) {
	req, err := h.parseListRequest(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.search.ListAssemblies(r.Context(), req)
	if err != nil {
		h.log.Error("Failed to list assemblies", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to list assemblies")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) ListAssembliesByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	req, err := h.parseListRequest(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.search.ListAssembliesByUser(r.Context(), int32(userID), req)
	if err != nil {
		h.log.Error("Failed to list assemblies by user", "user_id", userID, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to list assemblies")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Search handlers
func (h *SearchHandler) SearchComponents(w http.ResponseWriter, r *http.Request) {
	var req searchpb.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode search request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Query == "" {
		h.respondError(w, http.StatusBadRequest, "Search query is required")
		return
	}

	result, err := h.search.SearchComponents(r.Context(), &req)
	if err != nil {
		h.log.Error("Failed to search components", "query", req.Query, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to search components")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) SearchAssemblies(w http.ResponseWriter, r *http.Request) {
	var req searchpb.SearchAssembliesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode search assemblies request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.search.SearchAssemblies(r.Context(), &req)
	if err != nil {
		h.log.Error("Failed to search assemblies", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to search assemblies")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Compatibility handlers
func (h *SearchHandler) CheckCompatibility(w http.ResponseWriter, r *http.Request) {
	assemblyID, err := strconv.ParseInt(r.PathValue("assembly_id"), 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid assembly ID")
		return
	}

	result, err := h.search.CheckCompatibility(r.Context(), int32(assemblyID))
	if err != nil {
		h.log.Error("Failed to check compatibility", "assembly_id", assemblyID, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to check compatibility")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) GetCompatibleComponents(w http.ResponseWriter, r *http.Request) {
	var req searchpb.CompatibilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode compatibility request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	result, err := h.search.GetCompatibleComponents(r.Context(), &req)
	if err != nil {
		h.log.Error("Failed to get compatible components", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get compatible components")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) RecommendComponents(w http.ResponseWriter, r *http.Request) {
	var req searchpb.RecommendationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("Failed to decode recommendation request", "error", err)
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Purpose == "" {
		h.respondError(w, http.StatusBadRequest, "Purpose is required")
		return
	}

	result, err := h.search.RecommendComponents(r.Context(), &req)
	if err != nil {
		h.log.Error("Failed to recommend components", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to recommend components")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Price handlers
func (h *SearchHandler) GetComponentPrices(w http.ResponseWriter, r *http.Request) {
	componentType := r.PathValue("component_type")
	if componentType == "" {
		h.respondError(w, http.StatusBadRequest, "Component type is required")
		return
	}

	componentID, err := strconv.ParseInt(r.PathValue("component_id"), 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid component ID")
		return
	}

	shopIDStr := r.URL.Query().Get("shop_id")
	var shopID *int32
	if shopIDStr != "" {
		id, err := strconv.ParseInt(shopIDStr, 10, 32)
		if err != nil {
			h.respondError(w, http.StatusBadRequest, "Invalid shop ID")
			return
		}
		shopIDVal := int32(id)
		shopID = &shopIDVal
	}

	result, err := h.search.GetComponentPrices(r.Context(), componentType, int32(componentID), *shopID)
	if err != nil {
		h.log.Error("Failed to get component prices", "type", componentType, "id", componentID, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get component prices")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

func (h *SearchHandler) GetBestOffers(w http.ResponseWriter, r *http.Request) {
	componentType := r.PathValue("component_type")
	if componentType == "" {
		h.respondError(w, http.StatusBadRequest, "Component type is required")
		return
	}

	componentID, err := strconv.ParseInt(r.PathValue("component_id"), 10, 32)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid component ID")
		return
	}

	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 32)
	if err != nil || limit <= 0 {
		limit = 10
	}

	cheapestFirst := r.URL.Query().Get("cheapest") == "true"

	result, err := h.search.GetBestOffers(r.Context(), componentType, int32(componentID), int32(limit), cheapestFirst)
	if err != nil {
		h.log.Error("Failed to get best offers", "type", componentType, "id", componentID, "error", err)
		h.respondError(w, http.StatusInternalServerError, "Failed to get best offers")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// Helper methods
func (h *SearchHandler) parseListRequest(r *http.Request) (*searchpb.ListRequest, error) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 32)
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 32)
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := strings.ToLower(r.URL.Query().Get("sort_order"))
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	searchQuery := r.URL.Query().Get("search")

	// Parse filters from query parameters
	var filters []*searchpb.Filter
	for key, values := range r.URL.Query() {
		if strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]") {
			field := strings.TrimSuffix(strings.TrimPrefix(key, "filter["), "]")
			if len(values) > 0 {
				// Simple filter with "=" operator
				filters = append(filters, &searchpb.Filter{
					Field:    field,
					Operator: "=",
					Value:    values[0],
				})
			}
		}
	}

	// Parse complex filters from JSON
	filterJSON := r.URL.Query().Get("filters")
	if filterJSON != "" {
		var complexFilters []*searchpb.Filter
		if err := json.Unmarshal([]byte(filterJSON), &complexFilters); err != nil {
			h.log.Warn("Failed to parse filters JSON", "error", err)
		} else {
			filters = append(filters, complexFilters...)
		}
	}

	return &searchpb.ListRequest{
		Page:        int32(page),
		PageSize:    int32(pageSize),
		SortBy:      sortBy,
		SortOrder:   sortOrder,
		Filters:     filters,
		SearchQuery: &searchQuery,
	}, nil
}

func (h *SearchHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.Error("Failed to encode response", "error", err)
	}
}

func (h *SearchHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}
