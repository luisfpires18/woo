package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/middleware"
	"github.com/luisfpires18/woo/internal/service"
)

// CampAdminHandler handles admin CRUD endpoints for camp configuration.
type CampAdminHandler struct {
	campAdminService *service.CampAdminService
}

// NewCampAdminHandler creates a new CampAdminHandler.
func NewCampAdminHandler(campAdminService *service.CampAdminService) *CampAdminHandler {
	return &CampAdminHandler{campAdminService: campAdminService}
}

// RegisterRoutes registers admin camp routes. These are mounted under /api/admin/ with prefix stripping.
func (h *CampAdminHandler) RegisterRoutes(mux *http.ServeMux) {
	// Beast templates
	mux.HandleFunc("GET /beast-templates", h.ListBeastTemplates)
	mux.HandleFunc("POST /beast-templates", h.CreateBeastTemplate)
	mux.HandleFunc("PUT /beast-templates/{id}", h.UpdateBeastTemplate)
	mux.HandleFunc("DELETE /beast-templates/{id}", h.DeleteBeastTemplate)

	// Camp templates
	mux.HandleFunc("GET /camp-templates", h.ListCampTemplates)
	mux.HandleFunc("POST /camp-templates", h.CreateCampTemplate)
	mux.HandleFunc("PUT /camp-templates/{id}", h.UpdateCampTemplate)
	mux.HandleFunc("DELETE /camp-templates/{id}", h.DeleteCampTemplate)

	// Spawn rules
	mux.HandleFunc("GET /spawn-rules", h.ListSpawnRules)
	mux.HandleFunc("POST /spawn-rules", h.CreateSpawnRule)
	mux.HandleFunc("PUT /spawn-rules/{id}", h.UpdateSpawnRule)
	mux.HandleFunc("DELETE /spawn-rules/{id}", h.DeleteSpawnRule)

	// Reward tables
	mux.HandleFunc("GET /reward-tables", h.ListRewardTables)
	mux.HandleFunc("POST /reward-tables", h.CreateRewardTable)
	mux.HandleFunc("DELETE /reward-tables/{id}", h.DeleteRewardTable)

	// Battle tuning
	mux.HandleFunc("GET /battle-tuning", h.GetBattleTuning)
	mux.HandleFunc("PUT /battle-tuning", h.UpdateBattleTuning)
}

// ── Beast Templates ──────────────────────────────────────────────────────────

func (h *CampAdminHandler) ListBeastTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.campAdminService.ListBeastTemplates(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list beast templates")
		return
	}
	writeJSON(w, http.StatusOK, templates)
}

func (h *CampAdminHandler) CreateBeastTemplate(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req dto.CreateBeastTemplateRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	resp, err := h.campAdminService.CreateBeastTemplate(r.Context(), adminID, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create beast template")
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (h *CampAdminHandler) UpdateBeastTemplate(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.UpdateBeastTemplateRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	resp, err := h.campAdminService.UpdateBeastTemplate(r.Context(), adminID, id, req)
	if err != nil {
		if errors.Is(err, service.ErrBeastTemplateNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update beast template")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *CampAdminHandler) DeleteBeastTemplate(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.campAdminService.DeleteBeastTemplate(r.Context(), adminID, id); err != nil {
		if errors.Is(err, service.ErrBeastTemplateNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete beast template")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// ── Camp Templates ───────────────────────────────────────────────────────────

func (h *CampAdminHandler) ListCampTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.campAdminService.ListCampTemplates(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list camp templates")
		return
	}
	writeJSON(w, http.StatusOK, templates)
}

func (h *CampAdminHandler) CreateCampTemplate(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req dto.CreateCampTemplateRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	resp, err := h.campAdminService.CreateCampTemplate(r.Context(), adminID, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create camp template")
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (h *CampAdminHandler) UpdateCampTemplate(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.UpdateCampTemplateRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	resp, err := h.campAdminService.UpdateCampTemplate(r.Context(), adminID, id, req)
	if err != nil {
		if errors.Is(err, service.ErrCampTemplateNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update camp template")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *CampAdminHandler) DeleteCampTemplate(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.campAdminService.DeleteCampTemplate(r.Context(), adminID, id); err != nil {
		if errors.Is(err, service.ErrCampTemplateNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete camp template")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// ── Spawn Rules ──────────────────────────────────────────────────────────────

func (h *CampAdminHandler) ListSpawnRules(w http.ResponseWriter, r *http.Request) {
	rules, err := h.campAdminService.ListSpawnRules(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list spawn rules")
		return
	}
	writeJSON(w, http.StatusOK, rules)
}

func (h *CampAdminHandler) CreateSpawnRule(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req dto.CreateSpawnRuleRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	resp, err := h.campAdminService.CreateSpawnRule(r.Context(), adminID, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create spawn rule")
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (h *CampAdminHandler) UpdateSpawnRule(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req dto.UpdateSpawnRuleRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	resp, err := h.campAdminService.UpdateSpawnRule(r.Context(), adminID, id, req)
	if err != nil {
		if errors.Is(err, service.ErrSpawnRuleNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update spawn rule")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *CampAdminHandler) DeleteSpawnRule(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.campAdminService.DeleteSpawnRule(r.Context(), adminID, id); err != nil {
		if errors.Is(err, service.ErrSpawnRuleNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete spawn rule")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// ── Reward Tables ────────────────────────────────────────────────────────────

func (h *CampAdminHandler) ListRewardTables(w http.ResponseWriter, r *http.Request) {
	tables, err := h.campAdminService.ListRewardTables(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list reward tables")
		return
	}
	writeJSON(w, http.StatusOK, tables)
}

func (h *CampAdminHandler) CreateRewardTable(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req dto.CreateRewardTableRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	resp, err := h.campAdminService.CreateRewardTable(r.Context(), adminID, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create reward table")
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (h *CampAdminHandler) DeleteRewardTable(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.campAdminService.DeleteRewardTable(r.Context(), adminID, id); err != nil {
		if errors.Is(err, service.ErrRewardTableNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete reward table")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// ── Battle Tuning ────────────────────────────────────────────────────────────

func (h *CampAdminHandler) GetBattleTuning(w http.ResponseWriter, r *http.Request) {
	resp, err := h.campAdminService.GetBattleTuning(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get battle tuning")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *CampAdminHandler) UpdateBattleTuning(w http.ResponseWriter, r *http.Request) {
	adminID, ok := middleware.PlayerIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req dto.UpdateBattleTuningRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	resp, err := h.campAdminService.UpdateBattleTuning(r.Context(), adminID, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update battle tuning")
		return
	}
	writeJSON(w, http.StatusOK, resp)
}
