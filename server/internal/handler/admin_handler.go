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
	mapService   *service.MapService
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(adminService *service.AdminService, mapService *service.MapService) *AdminHandler {
	return &AdminHandler{adminService: adminService, mapService: mapService}
}

// RegisterRoutes registers admin routes on the given mux.
// All admin routes require auth + admin middleware applied externally.
func (h *AdminHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /players", h.ListPlayers)
	mux.HandleFunc("PATCH /players/{id}/role", h.UpdatePlayerRole)
	mux.HandleFunc("GET /stats", h.GetStats)
	mux.HandleFunc("GET /announcements", h.ListAnnouncements)
	mux.HandleFunc("POST /announcements", h.CreateAnnouncement)
	mux.HandleFunc("DELETE /announcements/{id}", h.DeleteAnnouncement)
	mux.HandleFunc("GET /assets", h.ListAssets)
	mux.HandleFunc("POST /assets", h.CreateAsset)
	mux.HandleFunc("DELETE /assets/{id}", h.DeleteAsset)
	mux.HandleFunc("GET /building-displays", h.ListBuildingDisplayConfigs)
	mux.HandleFunc("GET /building-displays/{id}", h.GetBuildingDisplayConfig)
	mux.HandleFunc("PUT /building-displays/{id}", h.UpdateBuildingDisplayConfig)
	mux.HandleFunc("GET /troop-displays", h.ListTroopDisplayConfigs)
	mux.HandleFunc("GET /troop-displays/{id}", h.GetTroopDisplayConfig)
	mux.HandleFunc("PUT /troop-displays/{id}", h.UpdateTroopDisplayConfig)
	mux.HandleFunc("GET /resource-buildings", h.ListResourceBuildingConfigs)
	mux.HandleFunc("GET /resource-buildings/{id}", h.GetResourceBuildingConfig)
	mux.HandleFunc("PUT /resource-buildings/{id}", h.UpdateResourceBuildingConfig)
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

// ListAssets handles GET /api/admin/assets.
func (h *AdminHandler) ListAssets(w http.ResponseWriter, r *http.Request) {
	resp, err := h.adminService.ListGameAssets(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list assets")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// CreateAsset handles POST /api/admin/assets.
// Creates a new game asset row (used for adding variants of zone/terrain tiles).
func (h *AdminHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateGameAssetRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	asset := &model.GameAsset{
		ID:          req.ID,
		Category:    req.Category,
		DisplayName: req.DisplayName,
		DefaultIcon: req.DefaultIcon,
	}

	if err := h.adminService.CreateGameAsset(r.Context(), asset); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Fetch the created asset for full DTO response.
	created, err := h.adminService.GetGameAsset(r.Context(), req.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to refetch asset")
		return
	}

	writeJSON(w, http.StatusCreated, dto.GameAssetDTO{
		ID:          created.ID,
		Category:    created.Category,
		DisplayName: created.DisplayName,
		DefaultIcon: created.DefaultIcon,
		UpdatedAt:   created.UpdatedAt,
	})
}

// DeleteAsset handles DELETE /api/admin/assets/{id}.
// Removes a game asset row.
func (h *AdminHandler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	assetID := r.PathValue("id")
	if assetID == "" {
		writeError(w, http.StatusBadRequest, "missing asset id")
		return
	}

	// Delete the DB row.
	if err := h.adminService.DeleteGameAsset(r.Context(), assetID); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "asset not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete asset")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "asset deleted"})
}

// --- Building display configs ---

// ListBuildingDisplayConfigs handles GET /api/admin/building-displays?kingdom=X.
func (h *AdminHandler) ListBuildingDisplayConfigs(w http.ResponseWriter, r *http.Request) {
	kingdom := r.URL.Query().Get("kingdom")

	resp, err := h.adminService.ListBuildingDisplayConfigs(r.Context(), kingdom)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list building display configs")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetBuildingDisplayConfig handles GET /api/admin/building-displays/{id}.
func (h *AdminHandler) GetBuildingDisplayConfig(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid building display config id")
		return
	}

	cfg, err := h.adminService.GetBuildingDisplayConfig(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "building display config not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get building display config")
		return
	}

	writeJSON(w, http.StatusOK, &dto.BuildingDisplayConfigDTO{
		ID:           cfg.ID,
		BuildingType: cfg.BuildingType,
		Kingdom:      cfg.Kingdom,
		DisplayName:  cfg.DisplayName,
		Description:  cfg.Description,
		DefaultIcon:  cfg.DefaultIcon,
		UpdatedAt:    cfg.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// UpdateBuildingDisplayConfig handles PUT /api/admin/building-displays/{id}.
func (h *AdminHandler) UpdateBuildingDisplayConfig(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid building display config id")
		return
	}

	var req dto.UpdateBuildingDisplayConfigRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if err := h.adminService.UpdateBuildingDisplayConfig(r.Context(), id, &req); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "building display config not found")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "config updated"})
}

// ── Troop display config handlers ───────────────────────────────────────────

// ListTroopDisplayConfigs handles GET /api/admin/troop-displays?kingdom=X.
func (h *AdminHandler) ListTroopDisplayConfigs(w http.ResponseWriter, r *http.Request) {
	kingdom := r.URL.Query().Get("kingdom")

	resp, err := h.adminService.ListTroopDisplayConfigs(r.Context(), kingdom)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list troop display configs")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetTroopDisplayConfig handles GET /api/admin/troop-displays/{id}.
func (h *AdminHandler) GetTroopDisplayConfig(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid troop display config id")
		return
	}

	cfg, err := h.adminService.GetTroopDisplayConfig(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "troop display config not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get troop display config")
		return
	}

	writeJSON(w, http.StatusOK, &dto.TroopDisplayConfigDTO{
		ID:               cfg.ID,
		TroopType:        cfg.TroopType,
		Kingdom:          cfg.Kingdom,
		TrainingBuilding: cfg.TrainingBuilding,
		DisplayName:      cfg.DisplayName,
		Description:      cfg.Description,
		DefaultIcon:      cfg.DefaultIcon,
		UpdatedAt:        cfg.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// UpdateTroopDisplayConfig handles PUT /api/admin/troop-displays/{id}.
func (h *AdminHandler) UpdateTroopDisplayConfig(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid troop display config id")
		return
	}

	var req dto.UpdateTroopDisplayConfigRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if err := h.adminService.UpdateTroopDisplayConfig(r.Context(), id, &req); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "troop display config not found")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "config updated"})
}

// ── Resource building config handlers ───────────────────────────────────────

// ListResourceBuildingConfigs handles GET /api/admin/resource-buildings?kingdom=X.
func (h *AdminHandler) ListResourceBuildingConfigs(w http.ResponseWriter, r *http.Request) {
	kingdom := r.URL.Query().Get("kingdom")

	resp, err := h.adminService.ListResourceBuildingConfigs(r.Context(), kingdom)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list resource building configs")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetResourceBuildingConfig handles GET /api/admin/resource-buildings/{id}.
func (h *AdminHandler) GetResourceBuildingConfig(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid resource building config id")
		return
	}

	cfg, err := h.adminService.GetResourceBuildingConfig(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "resource building config not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get resource building config")
		return
	}

	writeJSON(w, http.StatusOK, &dto.ResourceBuildingConfigDTO{
		ID:           cfg.ID,
		ResourceType: cfg.ResourceType,
		Slot:         cfg.Slot,
		Kingdom:      cfg.Kingdom,
		DisplayName:  cfg.DisplayName,
		Description:  cfg.Description,
		DefaultIcon:  cfg.DefaultIcon,
		UpdatedAt:    cfg.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

// UpdateResourceBuildingConfig handles PUT /api/admin/resource-buildings/{id}.
func (h *AdminHandler) UpdateResourceBuildingConfig(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid resource building config id")
		return
	}

	var req dto.UpdateResourceBuildingConfigRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if err := h.adminService.UpdateResourceBuildingConfig(r.Context(), id, &req); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "resource building config not found")
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "config updated"})
}
