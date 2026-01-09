package auth

import (
	domainErr "compConfigurator/internal/domain/errors"
	"compConfigurator/internal/httpapi/response"
	"compConfigurator/internal/service"
	"encoding/json"
	"errors"
	"net/http"
)

type Handler struct {
	auth *service.AuthService
}

func New(auth *service.AuthService) *Handler {
	return &Handler{auth: auth}
}

type registerRequest struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Nickname *string `json:"nickname"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
	User  any    `json:"user"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid json")
		return
	}

	u, token, err := h.auth.Register(r.Context(), req.Email, req.Password, req.Nickname)

	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrInvalidInput):
			response.Fail(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, domainErr.ErrConflict):
			response.Fail(w, http.StatusConflict, "email already exists")
		default:
			response.Fail(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.JSON(w, http.StatusCreated, authResponse{
		Token: token,
		User:  u,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid json")
		return
	}

	u, token, err := h.auth.Login(r.Context(), req.Email, req.Password)

	if err != nil {
		switch {
		case errors.Is(err, domainErr.ErrInvalidInput):
			response.Fail(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, domainErr.ErrUnauthorized):
			response.Fail(w, http.StatusUnauthorized, "invalid credentials")
		default:
			response.Fail(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	response.JSON(w, http.StatusCreated, authResponse{
		Token: token,
		User:  u,
	})
}
