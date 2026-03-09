package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/middleware"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/service"
)

// TrainingHandler handles troop training HTTP endpoints.
type TrainingHandler struct {
	trainingService *service.TrainingService
}

// NewTrainingHandler creates a new TrainingHandler.
func NewTrainingHandler(trainingService *service.TrainingService) *TrainingHandler {
	return &TrainingHandler{trainingService: trainingService}
}

// RegisterRoutes registers training routes on the given mux.
// These routes share the /api/villages/{id}/ prefix and require auth middleware externally.
func (h *TrainingHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/villages/{id}/train", h.StartTraining)
	mux.HandleFunc("GET /api/villages/{id}/train", h.GetTrainingQueue)
	mux.HandleFunc("GET /api/villages/{id}/train/cost", h.GetTrainingCost)
	mux.HandleFunc("DELETE /api/villages/{id}/train/{queueId}", h.CancelTraining)
	mux.HandleFunc("GET /api/villages/{id}/troops", h.ListTroops)
}

// RegisterAdminRoutes registers admin-only training routes.
func (h *TrainingHandler) RegisterAdminRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /training/{queueId}/complete", h.InstantComplete)
}

// StartTraining handles POST /api/villages/{id}/train.
func (h *TrainingHandler) StartTraining(w http.ResponseWriter, r *http.Request) {
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

	var req dto.StartTrainingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.trainingService.StartTraining(r.Context(), playerID, villageID, req.TroopType, req.Quantity)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUnknownTroop):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrInvalidQuantity):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrVillageNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrNotOwner):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, service.ErrTrainingBuildingReq):
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, model.ErrInsufficientResources):
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to start training")
		}
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// GetTrainingCost handles GET /api/villages/{id}/train/cost?troop_type=X&quantity=N.
func (h *TrainingHandler) GetTrainingCost(w http.ResponseWriter, r *http.Request) {
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

	troopType := r.URL.Query().Get("troop_type")
	if troopType == "" {
		writeError(w, http.StatusBadRequest, "troop_type is required")
		return
	}

	quantityStr := r.URL.Query().Get("quantity")
	quantity := 1
	if quantityStr != "" {
		quantity, err = strconv.Atoi(quantityStr)
		if err != nil || quantity < 1 {
			writeError(w, http.StatusBadRequest, "invalid quantity")
			return
		}
	}

	resp, err := h.trainingService.GetTrainingCost(r.Context(), playerID, villageID, troopType, quantity)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUnknownTroop):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrVillageNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrNotOwner):
			writeError(w, http.StatusForbidden, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to get training cost")
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// GetTrainingQueue handles GET /api/villages/{id}/train.
func (h *TrainingHandler) GetTrainingQueue(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	villageID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid village id")
		return
	}

	queue, err := h.trainingService.GetTrainingQueue(r.Context(), villageID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get training queue")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"queue": queue})
}

// CancelTraining handles DELETE /api/villages/{id}/train/{queueId}.
func (h *TrainingHandler) CancelTraining(w http.ResponseWriter, r *http.Request) {
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

	if err := h.trainingService.CancelTraining(r.Context(), playerID, villageID, queueID); err != nil {
		switch {
		case errors.Is(err, service.ErrVillageNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrNotOwner):
			writeError(w, http.StatusForbidden, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to cancel training")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListTroops handles GET /api/villages/{id}/troops.
func (h *TrainingHandler) ListTroops(w http.ResponseWriter, r *http.Request) {
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

	troops, err := h.trainingService.GetTroops(r.Context(), playerID, villageID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrVillageNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrNotOwner):
			writeError(w, http.StatusForbidden, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to list troops")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"troops": troops})
}

// InstantComplete handles POST /api/admin/training/{queueId}/complete.
// Sets the training queue item's completes_at to now so the game loop processes it instantly.
func (h *TrainingHandler) InstantComplete(w http.ResponseWriter, r *http.Request) {
	queueID, err := strconv.ParseInt(r.PathValue("queueId"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid queue id")
		return
	}

	if err := h.trainingService.InstantCompleteTraining(r.Context(), queueID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to instant complete training")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
