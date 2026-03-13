package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/middleware"
	"github.com/luisfpires18/woo/internal/service"
)

// ExpeditionHandler handles camp viewing and expedition endpoints.
type ExpeditionHandler struct {
	expeditionService *service.ExpeditionService
	campService       *service.CampService
}

// NewExpeditionHandler creates a new ExpeditionHandler.
func NewExpeditionHandler(expeditionService *service.ExpeditionService, campService *service.CampService) *ExpeditionHandler {
	return &ExpeditionHandler{
		expeditionService: expeditionService,
		campService:       campService,
	}
}

// RegisterRoutes registers expedition and camp routes on the given mux.
func (h *ExpeditionHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/camps", h.ListCamps)
	mux.HandleFunc("GET /api/camps/{id}", h.GetCamp)
	mux.HandleFunc("POST /api/villages/{id}/expeditions", h.DispatchExpedition)
	mux.HandleFunc("GET /api/expeditions", h.ListExpeditions)
	mux.HandleFunc("GET /api/battles/{id}", h.GetBattleReport)
	mux.HandleFunc("GET /api/battles/{id}/replay", h.GetBattleReplay)
}

// ListCamps handles GET /api/camps — list all active camps on the map.
func (h *ExpeditionHandler) ListCamps(w http.ResponseWriter, r *http.Request) {
	camps, err := h.campService.ListActiveCamps(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list camps")
		return
	}
	writeJSON(w, http.StatusOK, camps)
}

// GetCamp handles GET /api/camps/{id} — get a specific camp with beast details.
func (h *ExpeditionHandler) GetCamp(w http.ResponseWriter, r *http.Request) {
	campID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid camp id")
		return
	}

	camp, beasts, tmpl, err := h.campService.GetCampWithBeasts(r.Context(), campID)
	if err != nil {
		writeError(w, http.StatusNotFound, "camp not found")
		return
	}

	// Group beasts by template ID for the response
	type beastGroup struct {
		dto   dto.CampBeastResponse
		count int
	}
	groups := make(map[int64]*beastGroup)
	var order []int64
	for _, b := range beasts {
		g, ok := groups[b.BeastTemplateID]
		if !ok {
			order = append(order, b.BeastTemplateID)
			groups[b.BeastTemplateID] = &beastGroup{
				dto: dto.CampBeastResponse{
					BeastTemplateID:   b.BeastTemplateID,
					Name:              b.Name,
					SpriteKey:         b.SpriteKey,
					HP:                b.HP,
					MaxHP:             b.MaxHP,
					AttackPower:       b.AttackPower,
					AttackInterval:    b.AttackInterval,
					DefensePercent:    b.DefensePercent,
					CritChancePercent: b.CritChancePercent,
				},
				count: 1,
			}
		} else {
			g.count++
		}
	}
	beastDTOs := make([]dto.CampBeastResponse, 0, len(groups))
	for _, bid := range order {
		g := groups[bid]
		g.dto.Count = g.count
		beastDTOs = append(beastDTOs, g.dto)
	}

	resp := dto.CampResponse{
		ID:           camp.ID,
		TemplateName: tmpl.Name,
		Tier:         tmpl.Tier,
		SpriteKey:    tmpl.SpriteKey,
		TileX:        camp.TileX,
		TileY:        camp.TileY,
		Status:       camp.Status,
		Beasts:       beastDTOs,
	}
	if t, err2 := time.Parse(time.RFC3339, camp.SpawnedAt); err2 == nil {
		resp.SpawnedAt = t
	}
	writeJSON(w, http.StatusOK, resp)
}

// DispatchExpedition handles POST /api/villages/{id}/expeditions — send troops to a camp.
func (h *ExpeditionHandler) DispatchExpedition(w http.ResponseWriter, r *http.Request) {
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

	var req dto.DispatchExpeditionRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	resp, err := h.expeditionService.DispatchExpedition(r.Context(), playerID, villageID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrVillageNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrNotOwner):
			writeError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, service.ErrCampNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrCampNotActive):
			writeError(w, http.StatusConflict, err.Error())
		case errors.Is(err, service.ErrNoTroopsSent):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrInsufficientTroops):
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to dispatch expedition")
		}
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

// ListExpeditions handles GET /api/expeditions — list all expeditions for the player.
func (h *ExpeditionHandler) ListExpeditions(w http.ResponseWriter, r *http.Request) {
	playerID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	exps, err := h.expeditionService.GetExpeditionsByPlayer(r.Context(), playerID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list expeditions")
		return
	}

	if exps == nil {
		exps = []dto.ExpeditionResponse{}
	}
	writeJSON(w, http.StatusOK, exps)
}

// GetBattleReport handles GET /api/battles/{id} — get the battle summary/report.
func (h *ExpeditionHandler) GetBattleReport(w http.ResponseWriter, r *http.Request) {
	playerID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	battleID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid battle id")
		return
	}

	report, err := h.expeditionService.GetBattleReport(r.Context(), playerID, battleID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBattleNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrNotOwner):
			writeError(w, http.StatusForbidden, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to get battle report")
		}
		return
	}

	writeJSON(w, http.StatusOK, report)
}

// GetBattleReplay handles GET /api/battles/{id}/replay — get the replay data blob.
func (h *ExpeditionHandler) GetBattleReplay(w http.ResponseWriter, r *http.Request) {
	playerID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	battleID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid battle id")
		return
	}

	replay, err := h.expeditionService.GetBattleReplay(r.Context(), playerID, battleID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBattleNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrNotOwner):
			writeError(w, http.StatusForbidden, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "failed to get replay")
		}
		return
	}

	var replayData any
	if err := json.Unmarshal(replay, &replayData); err != nil {
		writeError(w, http.StatusInternalServerError, "corrupt replay data")
		return
	}
	writeJSON(w, http.StatusOK, replayData)
}
