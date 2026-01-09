package me

import (
	domainErr "compConfigurator/internal/domain/errors"
	"compConfigurator/internal/httpapi/middleware"
	"compConfigurator/internal/httpapi/response"
	"compConfigurator/internal/repo"
	"errors"
	"net/http"
)

type Handler struct {
	users *repo.UsersRepo
}

func New(users *repo.UsersRepo) *Handler {
	return &Handler{users: users}
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Fail(w, http.StatusUnauthorized, domainErr.ErrUnauthorized.Error())
		return
	}

	u, err := h.users.GetByID(r.Context(), userID)

	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrNotFound):
			response.Fail(w, http.StatusNotFound, "user not found")
		default:
			response.Fail(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.JSON(w, http.StatusOK, u)
}
