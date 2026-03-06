package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/middleware"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
	"github.com/luisfpires18/woo/internal/service"
)

// PlayerHandler handles player-related HTTP endpoints.
type PlayerHandler struct {
	playerRepo     repository.PlayerRepository
	authService    *service.AuthService
	villageService *service.VillageService
}

// NewPlayerHandler creates a new PlayerHandler.
func NewPlayerHandler(
	playerRepo repository.PlayerRepository,
	authService *service.AuthService,
	villageService *service.VillageService,
) *PlayerHandler {
	return &PlayerHandler{
		playerRepo:     playerRepo,
		authService:    authService,
		villageService: villageService,
	}
}

// ChooseKingdom handles PUT /api/player/kingdom.
// Sets the player's kingdom (one-time, post-registration) and creates their first village.
func (h *PlayerHandler) ChooseKingdom(w http.ResponseWriter, r *http.Request) {
	playerID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req dto.ChooseKingdomRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	// Validate kingdom value
	if !service.IsValidKingdom(req.Kingdom) {
		writeError(w, http.StatusBadRequest, "invalid kingdom")
		return
	}

	// Get current player
	player, err := h.playerRepo.GetByID(r.Context(), playerID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "player not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch player")
		return
	}

	// Reject if kingdom already set
	if player.Kingdom != "" {
		writeError(w, http.StatusConflict, "kingdom already chosen")
		return
	}

	// Set kingdom
	if err := h.playerRepo.UpdateKingdom(r.Context(), playerID, req.Kingdom); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update kingdom")
		return
	}

	// Create first village
	village, err := h.villageService.CreateFirstVillage(r.Context(), playerID, req.Kingdom, player.Username)
	if err != nil {
		slog.Error("failed to create first village after kingdom selection", "player_id", playerID, "error", err)
		writeError(w, http.StatusInternalServerError, "kingdom set but village creation failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"player": &dto.PlayerInfo{
			ID:       player.ID,
			Username: player.Username,
			Email:    player.Email,
			Kingdom:  req.Kingdom,
			Role:     player.Role,
		},
		"village_id": village.ID,
	})
}
