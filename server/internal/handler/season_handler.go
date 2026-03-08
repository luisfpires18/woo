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

// SeasonHandler handles season-related HTTP endpoints.
type SeasonHandler struct {
	seasonService *service.SeasonService
}

// NewSeasonHandler creates a new SeasonHandler.
func NewSeasonHandler(seasonService *service.SeasonService) *SeasonHandler {
	return &SeasonHandler{seasonService: seasonService}
}

// RegisterAdminRoutes registers admin season routes on the given mux.
// These are mounted under /api/admin (after prefix stripping).
func (h *SeasonHandler) RegisterAdminRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /seasons", h.AdminListSeasons)
	mux.HandleFunc("POST /seasons", h.AdminCreateSeason)
	mux.HandleFunc("GET /seasons/{id}", h.AdminGetSeason)
	mux.HandleFunc("PUT /seasons/{id}", h.AdminUpdateSeason)
	mux.HandleFunc("DELETE /seasons/{id}", h.AdminDeleteSeason)
	mux.HandleFunc("POST /seasons/{id}/launch", h.AdminLaunchSeason)
	mux.HandleFunc("POST /seasons/{id}/end", h.AdminEndSeason)
	mux.HandleFunc("POST /seasons/{id}/archive", h.AdminArchiveSeason)
}

// RegisterPlayerRoutes registers player-facing season routes.
// These are mounted under /api/seasons (protected, after prefix stripping).
func (h *SeasonHandler) RegisterPlayerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.ListSeasons)
	mux.HandleFunc("GET /my", h.GetMySeasons)
	mux.HandleFunc("GET /{id}", h.GetSeason)
	mux.HandleFunc("POST /{id}/join", h.JoinSeason)
}

// ── Admin handlers ───────────────────────────────────────────────────────────

func (h *SeasonHandler) AdminListSeasons(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")
	seasons, err := h.seasonService.ListSeasons(r.Context(), statusFilter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list seasons")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"seasons": seasons})
}

func (h *SeasonHandler) AdminCreateSeason(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSeasonRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	season, err := h.seasonService.CreateSeason(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrSeasonNameRequired) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create season")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"season": season})
}

func (h *SeasonHandler) AdminGetSeason(w http.ResponseWriter, r *http.Request) {
	id, err := parseSeasonID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid season id")
		return
	}

	playerID, _ := middleware.PlayerIDFromContext(r.Context())
	season, err := h.seasonService.GetSeason(r.Context(), id, playerID)
	if err != nil {
		if errors.Is(err, service.ErrSeasonNotFound) {
			writeError(w, http.StatusNotFound, "season not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get season")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"season": season})
}

func (h *SeasonHandler) AdminUpdateSeason(w http.ResponseWriter, r *http.Request) {
	id, err := parseSeasonID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid season id")
		return
	}

	var req dto.UpdateSeasonRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	season, err := h.seasonService.UpdateSeason(r.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSeasonNotFound):
			writeError(w, http.StatusNotFound, "season not found")
		case errors.Is(err, service.ErrSeasonNotUpcoming):
			writeError(w, http.StatusConflict, "can only edit upcoming seasons")
		default:
			writeError(w, http.StatusInternalServerError, "failed to update season")
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"season": season})
}

func (h *SeasonHandler) AdminDeleteSeason(w http.ResponseWriter, r *http.Request) {
	id, err := parseSeasonID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid season id")
		return
	}

	if err := h.seasonService.DeleteSeason(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, service.ErrSeasonNotFound):
			writeError(w, http.StatusNotFound, "season not found")
		case errors.Is(err, service.ErrSeasonNotUpcoming):
			writeError(w, http.StatusConflict, "can only delete upcoming seasons")
		default:
			writeError(w, http.StatusInternalServerError, "failed to delete season")
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"deleted": true})
}

func (h *SeasonHandler) AdminLaunchSeason(w http.ResponseWriter, r *http.Request) {
	id, err := parseSeasonID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid season id")
		return
	}

	season, err := h.seasonService.LaunchSeason(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSeasonNotFound):
			writeError(w, http.StatusNotFound, "season not found")
		case errors.Is(err, service.ErrSeasonNotUpcoming):
			writeError(w, http.StatusConflict, "season must be in upcoming status to launch")
		default:
			writeError(w, http.StatusInternalServerError, "failed to launch season")
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"season": season})
}

func (h *SeasonHandler) AdminEndSeason(w http.ResponseWriter, r *http.Request) {
	id, err := parseSeasonID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid season id")
		return
	}

	season, err := h.seasonService.EndSeason(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSeasonNotFound):
			writeError(w, http.StatusNotFound, "season not found")
		case errors.Is(err, service.ErrSeasonNotActive):
			writeError(w, http.StatusConflict, "season must be active to end")
		default:
			writeError(w, http.StatusInternalServerError, "failed to end season")
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"season": season})
}

func (h *SeasonHandler) AdminArchiveSeason(w http.ResponseWriter, r *http.Request) {
	id, err := parseSeasonID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid season id")
		return
	}

	season, err := h.seasonService.ArchiveSeason(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSeasonNotFound):
			writeError(w, http.StatusNotFound, "season not found")
		case errors.Is(err, service.ErrSeasonNotEnded):
			writeError(w, http.StatusConflict, "season must be ended to archive")
		default:
			writeError(w, http.StatusInternalServerError, "failed to archive season")
		}
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"season": season})
}

// ── Player handlers ──────────────────────────────────────────────────────────

func (h *SeasonHandler) ListSeasons(w http.ResponseWriter, r *http.Request) {
	statusFilter := r.URL.Query().Get("status")
	seasons, err := h.seasonService.ListSeasons(r.Context(), statusFilter)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list seasons")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"seasons": seasons})
}

func (h *SeasonHandler) GetSeason(w http.ResponseWriter, r *http.Request) {
	id, err := parseSeasonID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid season id")
		return
	}

	playerID, _ := middleware.PlayerIDFromContext(r.Context())
	season, err := h.seasonService.GetSeason(r.Context(), id, playerID)
	if err != nil {
		if errors.Is(err, service.ErrSeasonNotFound) {
			writeError(w, http.StatusNotFound, "season not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get season")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"season": season})
}

func (h *SeasonHandler) JoinSeason(w http.ResponseWriter, r *http.Request) {
	id, err := parseSeasonID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid season id")
		return
	}

	playerID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req dto.JoinSeasonRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	detail, villageID, err := h.seasonService.JoinSeason(r.Context(), id, playerID, req.Kingdom)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidKingdom):
			writeError(w, http.StatusBadRequest, "invalid kingdom")
		case errors.Is(err, service.ErrSeasonNotFound):
			writeError(w, http.StatusNotFound, "season not found")
		case errors.Is(err, service.ErrSeasonNotActive):
			writeError(w, http.StatusConflict, "season is not active")
		case errors.Is(err, service.ErrAlreadyJoined):
			writeError(w, http.StatusConflict, "already joined this season")
		case errors.Is(err, model.ErrNotFound):
			writeError(w, http.StatusNotFound, "player not found")
		default:
			writeError(w, http.StatusInternalServerError, "failed to join season")
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"season":     detail,
		"village_id": villageID,
	})
}

func (h *SeasonHandler) GetMySeasons(w http.ResponseWriter, r *http.Request) {
	playerID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	seasons, err := h.seasonService.GetMySeasons(r.Context(), playerID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get your seasons")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"seasons": seasons})
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func parseSeasonID(r *http.Request) (int64, error) {
	idStr := r.PathValue("id")
	return strconv.ParseInt(idStr, 10, 64)
}
