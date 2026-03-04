package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/luisfpires18/woo/internal/middleware"
	"github.com/luisfpires18/woo/internal/service"
)

// VillageHandler handles village HTTP endpoints.
type VillageHandler struct {
	villageService *service.VillageService
}

// NewVillageHandler creates a new VillageHandler.
func NewVillageHandler(villageService *service.VillageService) *VillageHandler {
	return &VillageHandler{villageService: villageService}
}

// RegisterRoutes registers village routes on the given mux.
// All village routes require authentication middleware to be applied externally.
func (h *VillageHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/villages", h.ListVillages)
	mux.HandleFunc("GET /api/villages/{id}", h.GetVillage)
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

	writeJSON(w, http.StatusOK, resp)
}
