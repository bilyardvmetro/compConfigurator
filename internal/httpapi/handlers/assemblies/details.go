package assemblies

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	domainerr "compConfigurator/internal/domain/errors"
	"compConfigurator/internal/httpapi/middleware"
	"compConfigurator/internal/httpapi/response"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) Details(w http.ResponseWriter, r *http.Request) {
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

	// h.detailsSvc нужно добавить в Handler (см. router шаг)
	out, err := h.detailsSvc.Get(r.Context(), userID, assemblyID)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, domainerr.ErrNotFound) {
			response.Fail(w, http.StatusNotFound, "assembly not found")
			return
		}
		response.Fail(w, http.StatusInternalServerError, "internal error")
		return
	}

	response.JSON(w, http.StatusOK, out)
}
