package assemblies

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	domainerr "compConfigurator/internal/domain/errors"
	"compConfigurator/internal/httpapi/middleware"
	"compConfigurator/internal/httpapi/response"
	"compConfigurator/internal/service/pricing"

	"github.com/go-chi/chi/v5"
)

type PricingHandlerDeps interface {
	GetForOwner(ctx context.Context, userID, assemblyID int64) (pricing.PricingResult, error)
}

// если у тебя Handler уже есть — просто добавь поле pricingSvc *pricing.Service
func (h *Handler) Pricing(w http.ResponseWriter, r *http.Request) {
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

	out, err := h.pricingSvc.GetForOwner(r.Context(), userID, assemblyID)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			response.Fail(w, http.StatusNotFound, "assembly not found")
			return
		}
		response.Fail(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusOK, out)
}
