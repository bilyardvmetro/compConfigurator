package public

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	domainErr "compConfigurator/internal/domain/errors"
	"compConfigurator/internal/httpapi/middleware"
	"compConfigurator/internal/httpapi/response"

	"github.com/go-chi/chi/v5"
)

type cloneReq struct {
	Name string `json:"name"`
}

func (h *Handler) ClonePublicAssembly(w http.ResponseWriter, r *http.Request) {
	// 1. assembly_id из URL
	srcStr := chi.URLParam(r, "id")
	srcID, err := strconv.ParseInt(srcStr, 10, 64)
	if err != nil || srcID <= 0 {
		response.Fail(w, http.StatusBadRequest, domainErr.ErrInvalidInput.Error())
		return
	}

	// 2. user_id из твоего helper'а
	userID, ok := middleware.GetUserID(r)
	if !ok || userID <= 0 {
		response.Fail(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// 3. тело запроса (опционально)
	var req cloneReq
	_ = json.NewDecoder(r.Body).Decode(&req)

	// 4. клонирование
	newID, err := h.clone.ClonePublicAssembly(
		r.Context(),
		srcID,
		userID,
		req.Name,
	)
	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrNotFound):
			response.Fail(w, http.StatusNotFound, "public assembly not found")
		default:
			// тут могут прилетать ошибки из БД:
			// - триггеры совместимости
			// - FK
			response.Fail(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	// 5. ответ
	response.JSON(w, http.StatusCreated, map[string]any{
		"assembly_id": newID,
	})
}
