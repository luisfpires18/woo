package handler

import (
	"errors"
	"net/http"

	"github.com/luisfpires18/woo/internal/middleware"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/service"
)

// PlayerHandler handles player-related HTTP endpoints.
type PlayerHandler struct {
	playerService *service.PlayerService
	seasonService *service.SeasonService
}

// NewPlayerHandler creates a new PlayerHandler.
func NewPlayerHandler(playerService *service.PlayerService, seasonService *service.SeasonService) *PlayerHandler {
	return &PlayerHandler{playerService: playerService, seasonService: seasonService}
}

// ChooseKingdom handles PUT /api/player/kingdom.
// Sets the player's kingdom (one-time, post-registration) and creates their first village.
func (h *PlayerHandler) ChooseKingdom(w http.ResponseWriter, r *http.Request) {
	playerID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req struct {
		Kingdom string `json:"kingdom"`
	}
	if !decodeJSON(w, r, &req) {
		return
	}

	playerInfo, villageID, err := h.playerService.ChooseKingdom(r.Context(), playerID, req.Kingdom)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidKingdom):
			writeError(w, http.StatusBadRequest, "invalid kingdom")
		case errors.Is(err, service.ErrKingdomAlreadyChosen):
			writeError(w, http.StatusConflict, "kingdom already chosen")
		case errors.Is(err, model.ErrNotFound):
			writeError(w, http.StatusNotFound, "player not found")
		default:
			writeError(w, http.StatusInternalServerError, "village creation failed — is a map template applied?")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"player":     playerInfo,
		"village_id": villageID,
	})
}

// GetMe handles GET /api/player/me — returns the current player's info.
func (h *PlayerHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	playerID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	playerInfo, err := h.playerService.GetMe(r.Context(), playerID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "player not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch player")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"player": playerInfo,
	})
}

// GetProfile handles GET /api/player/profile — returns cross-season profile data.
func (h *PlayerHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	playerID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	profile, err := h.seasonService.GetPlayerProfile(r.Context(), playerID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "player not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to fetch profile")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"profile": profile,
	})
}
