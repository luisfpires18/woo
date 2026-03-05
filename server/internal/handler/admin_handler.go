package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/middleware"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/service"
)

// AdminHandler handles admin HTTP endpoints.
type AdminHandler struct {
	adminService *service.AdminService
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// RegisterRoutes registers admin routes on the given mux.
// All admin routes require auth + admin middleware applied externally.
func (h *AdminHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /players", h.ListPlayers)
	mux.HandleFunc("PATCH /players/{id}/role", h.UpdatePlayerRole)
	mux.HandleFunc("GET /config", h.GetWorldConfig)
	mux.HandleFunc("PUT /config/{key}", h.SetWorldConfig)
	mux.HandleFunc("GET /stats", h.GetStats)
	mux.HandleFunc("GET /announcements", h.ListAnnouncements)
	mux.HandleFunc("POST /announcements", h.CreateAnnouncement)
	mux.HandleFunc("DELETE /announcements/{id}", h.DeleteAnnouncement)
}

// ListPlayers handles GET /api/admin/players?offset=0&limit=20.
func (h *AdminHandler) ListPlayers(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}

	resp, err := h.adminService.ListPlayers(r.Context(), offset, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list players")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// UpdatePlayerRole handles PATCH /api/admin/players/{id}/role.
func (h *AdminHandler) UpdatePlayerRole(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid player id")
		return
	}

	var req dto.UpdateRoleRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if err := h.adminService.UpdatePlayerRole(r.Context(), id, req.Role); err != nil {
		if errors.Is(err, service.ErrInvalidRole) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "player not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update role")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "role updated"})
}

// GetWorldConfig handles GET /api/admin/config.
func (h *AdminHandler) GetWorldConfig(w http.ResponseWriter, r *http.Request) {
	resp, err := h.adminService.GetWorldConfig(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get config")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// SetWorldConfig handles PUT /api/admin/config/{key}.
func (h *AdminHandler) SetWorldConfig(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		writeError(w, http.StatusBadRequest, "missing config key")
		return
	}

	var req dto.SetConfigRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if err := h.adminService.SetWorldConfig(r.Context(), key, req.Value); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "config key not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to set config")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "config updated"})
}

// GetStats handles GET /api/admin/stats.
func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	resp, err := h.adminService.GetStats(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get stats")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// ListAnnouncements handles GET /api/admin/announcements.
func (h *AdminHandler) ListAnnouncements(w http.ResponseWriter, r *http.Request) {
	list, err := h.adminService.ListAnnouncements(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list announcements")
		return
	}

	writeJSON(w, http.StatusOK, list)
}

// CreateAnnouncement handles POST /api/admin/announcements.
func (h *AdminHandler) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	authorID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req dto.CreateAnnouncementRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	resp, err := h.adminService.CreateAnnouncement(r.Context(), &req, authorID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// DeleteAnnouncement handles DELETE /api/admin/announcements/{id}.
func (h *AdminHandler) DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid announcement id")
		return
	}

	if err := h.adminService.DeleteAnnouncement(r.Context(), id); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "announcement not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete announcement")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "announcement deleted"})
}
