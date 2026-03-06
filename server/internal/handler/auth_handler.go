package handler

import (
	"errors"
	"net/http"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/service"
)

// AuthHandler handles authentication HTTP endpoints.
type AuthHandler struct {
	authService    *service.AuthService
	villageService *service.VillageService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *service.AuthService, villageService *service.VillageService) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		villageService: villageService,
	}
}

// RegisterRoutes registers auth routes on the given mux.
func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/register", h.Register)
	mux.HandleFunc("POST /api/auth/login", h.Login)
	mux.HandleFunc("POST /api/auth/refresh", h.Refresh)
	mux.HandleFunc("POST /api/auth/logout", h.Logout)
}

// Register handles POST /api/auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	resp, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		status := mapAuthError(err)
		writeError(w, status, err.Error())
		return
	}

	// Village is created later when the player chooses a kingdom
	writeJSON(w, http.StatusCreated, resp)
}

// Login handles POST /api/auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	resp, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		status := mapAuthError(err)
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// Refresh handles POST /api/auth/refresh.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	resp, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		status := mapAuthError(err)
		writeError(w, status, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// Logout handles POST /api/auth/logout.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if err := h.authService.Logout(r.Context(), req.RefreshToken); err != nil {
		writeError(w, http.StatusInternalServerError, "logout failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

// mapAuthError maps service errors to HTTP status codes.
func mapAuthError(err error) int {
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		return http.StatusUnauthorized
	case errors.Is(err, service.ErrEmailTaken),
		errors.Is(err, service.ErrUsernameTaken):
		return http.StatusConflict
	case errors.Is(err, service.ErrInvalidKingdom),
		errors.Is(err, service.ErrWeakPassword),
		errors.Is(err, service.ErrInvalidEmail),
		errors.Is(err, service.ErrInvalidUsername):
		return http.StatusBadRequest
	case errors.Is(err, service.ErrInvalidRefreshToken):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
