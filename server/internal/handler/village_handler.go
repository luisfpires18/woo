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

// VillageHandler handles village HTTP endpoints.
type VillageHandler struct {
	villageService  *service.VillageService
	buildingService *service.BuildingService
}

// NewVillageHandler creates a new VillageHandler.
func NewVillageHandler(villageService *service.VillageService, buildingService *service.BuildingService) *VillageHandler {
	return &VillageHandler{
		villageService:  villageService,
		buildingService: buildingService,
	}
}

// RegisterRoutes registers village routes on the given mux.
// All village routes require authentication middleware to be applied externally.
func (h *VillageHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/villages", h.ListVillages)
	mux.HandleFunc("GET /api/villages/{id}", h.GetVillage)
	mux.HandleFunc("PUT /api/villages/{id}/name", h.RenameVillage)
	mux.HandleFunc("POST /api/villages/{id}/upgrade", h.StartUpgrade)
	mux.HandleFunc("GET /api/villages/{id}/upgrade/cost", h.GetUpgradeCost)
	mux.HandleFunc("DELETE /api/villages/{id}/upgrade/{queueId}", h.CancelUpgrade)
}

// ListVillages handles GET /api/villages.
func (h *VillageHandler) ListVillages(w http.ResponseWriter, r *http.Request) {
	playerID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	villages, err := h.villageService.ListVillages(r.Context(), playerID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list villages")
		return
	}

	writeJSON(w, http.StatusOK, villages)
}

// GetVillage handles GET /api/villages/{id}.
func (h *VillageHandler) GetVillage(w http.ResponseWriter, r *http.Request) {
	playerID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	idStr := r.PathValue("id")
	villageID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid village id")
		return
	}

	resp, err := h.villageService.GetVillage(r.Context(), villageID, playerID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrVillageNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrNotOwner):
			writeError(w, http.StatusForbidden, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to get village")
		}
		return
	}

	// Attach build queue
	queue, err := h.buildingService.GetBuildQueue(r.Context(), villageID)
	if err == nil {
		resp.BuildQueue = queue
	}

	writeJSON(w, http.StatusOK, resp)
}

// StartUpgrade handles POST /api/villages/{id}/upgrade.
func (h *VillageHandler) StartUpgrade(w http.ResponseWriter, r *http.Request) {
	playerID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	villageID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid village id")
		return
	}

	var req dto.StartUpgradeRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	if req.BuildingType == "" {
		writeError(w, http.StatusBadRequest, "building_type is required")
		return
	}

	resp, err := h.buildingService.StartUpgrade(r.Context(), playerID, villageID, req.BuildingType)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrVillageNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrNotOwner):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, service.ErrUnknownBuilding):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, model.ErrMaxLevelReached):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, model.ErrPrerequisitesNotMet):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, model.ErrInsufficientResources):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, model.ErrBuildingInProgress):
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to start upgrade")
		}
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// GetUpgradeCost handles GET /api/villages/{id}/upgrade/cost?building_type=X.
func (h *VillageHandler) GetUpgradeCost(w http.ResponseWriter, r *http.Request) {
	playerID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	villageID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid village id")
		return
	}

	buildingType := r.URL.Query().Get("building_type")
	if buildingType == "" {
		writeError(w, http.StatusBadRequest, "building_type query parameter is required")
		return
	}

	resp, err := h.buildingService.GetUpgradeCost(r.Context(), playerID, villageID, buildingType)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrVillageNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrNotOwner):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, service.ErrUnknownBuilding):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, model.ErrMaxLevelReached):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to get upgrade cost")
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// CancelUpgrade handles DELETE /api/villages/{id}/upgrade/{queueId}.
func (h *VillageHandler) CancelUpgrade(w http.ResponseWriter, r *http.Request) {
	playerID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	villageID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid village id")
		return
	}

	queueID, err := strconv.ParseInt(r.PathValue("queueId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid queue id")
		return
	}

	err = h.buildingService.CancelUpgrade(r.Context(), playerID, villageID, queueID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrVillageNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrNotOwner):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, model.ErrNotFound):
			writeError(w, http.StatusNotFound, "queue item not found")
		default:
			writeError(w, http.StatusInternalServerError, "failed to cancel upgrade")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "upgrade cancelled"})
}

// RenameVillage handles PUT /api/villages/{id}/name.
func (h *VillageHandler) RenameVillage(w http.ResponseWriter, r *http.Request) {
	playerID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	villageID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid village id")
		return
	}

	var req dto.RenameVillageRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	resp, err := h.villageService.RenameVillage(r.Context(), villageID, playerID, req.Name)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrVillageNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrNotOwner):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, service.ErrInvalidName):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to rename village")
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
