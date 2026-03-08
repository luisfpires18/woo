package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
	"github.com/luisfpires18/woo/internal/service"
)

// TemplateHandler handles admin map template endpoints.
type TemplateHandler struct {
	templateService *service.TemplateService
}

// NewTemplateHandler creates a new TemplateHandler.
func NewTemplateHandler(templateService *service.TemplateService) *TemplateHandler {
	return &TemplateHandler{templateService: templateService}
}

// RegisterRoutes registers template routes on the given mux.
// These are mounted under /api/admin/templates (with StripPrefix).
func (h *TemplateHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /templates", h.ListTemplates)
	mux.HandleFunc("POST /templates", h.CreateTemplate)
	mux.HandleFunc("GET /templates/{name}", h.GetTemplate)
	mux.HandleFunc("DELETE /templates/{name}", h.DeleteTemplate)
	mux.HandleFunc("PUT /templates/{name}/resize", h.ResizeTemplate)
	mux.HandleFunc("PUT /templates/{name}/terrain", h.UpdateTerrain)
	mux.HandleFunc("PUT /templates/{name}/zones", h.UpdateZones)
	mux.HandleFunc("POST /templates/{name}/apply", h.ApplyTemplate)
	mux.HandleFunc("GET /templates/{name}/export", h.ExportTemplate)
	mux.HandleFunc("POST /templates/import", h.ImportTemplate)
}

// ListTemplates handles GET /api/admin/templates.
func (h *TemplateHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	infos, err := h.templateService.ListTemplates()
	if err != nil {
		slog.Error("failed to list templates", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to list templates")
		return
	}

	if infos == nil {
		infos = []repository.TemplateInfo{}
	}

	writeJSON(w, http.StatusOK, map[string]any{"templates": infos})
}

// CreateTemplate handles POST /api/admin/templates.
func (h *TemplateHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTemplateRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	tmpl, err := h.templateService.CreateTemplate(req.Name, req.Description, req.MapSize)
	if err != nil {
		if err == service.ErrTemplateExists {
			writeError(w, http.StatusConflict, "template with this name already exists")
			return
		}
		slog.Error("failed to create template", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to create template")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"message":  "template created",
		"template": templateToSummary(tmpl),
	})
}

// GetTemplate handles GET /api/admin/templates/{name}.
func (h *TemplateHandler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "template name is required")
		return
	}

	tmpl, err := h.templateService.GetTemplate(name)
	if err != nil {
		writeError(w, http.StatusNotFound, "template not found")
		return
	}

	writeJSON(w, http.StatusOK, tmpl)
}

// DeleteTemplate handles DELETE /api/admin/templates/{name}.
func (h *TemplateHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "template name is required")
		return
	}

	if err := h.templateService.DeleteTemplate(name); err != nil {
		writeError(w, http.StatusNotFound, "template not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "template deleted"})
}

// UpdateTerrain handles PUT /api/admin/templates/{name}/terrain.
func (h *TemplateHandler) UpdateTerrain(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "template name is required")
		return
	}

	var req dto.UpdateTemplateTerrainRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	updates := make([]model.TileTerrainUpdate, len(req.Tiles))
	for i, t := range req.Tiles {
		updates[i] = model.TileTerrainUpdate{X: t.X, Y: t.Y, TerrainType: t.TerrainType}
	}

	if err := h.templateService.UpdateTemplateTerrain(name, updates); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "terrain updated"})
}

// UpdateZones handles PUT /api/admin/templates/{name}/zones.
func (h *TemplateHandler) UpdateZones(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "template name is required")
		return
	}

	var req dto.UpdateTemplateZonesRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	updates := make([]model.TileZoneUpdate, len(req.Tiles))
	for i, t := range req.Tiles {
		updates[i] = model.TileZoneUpdate{X: t.X, Y: t.Y, KingdomZone: t.KingdomZone}
	}

	if err := h.templateService.UpdateTemplateZones(name, updates); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "zones updated"})
}

// ApplyTemplate handles POST /api/admin/templates/{name}/apply.
func (h *TemplateHandler) ApplyTemplate(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "template name is required")
		return
	}

	var req dto.ApplyTemplateRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if !req.Confirm {
		writeError(w, http.StatusBadRequest, "confirm must be true to apply template")
		return
	}

	if err := h.templateService.ApplyTemplate(r.Context(), name); err != nil {
		slog.Error("failed to apply template", "name", name, "error", err)
		writeError(w, http.StatusInternalServerError, "failed to apply template: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "template applied to live map"})
}

// ExportTemplate handles GET /api/admin/templates/{name}/export — downloads the template as JSON.
func (h *TemplateHandler) ExportTemplate(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "template name is required")
		return
	}

	tmpl, err := h.templateService.GetTemplate(name)
	if err != nil {
		writeError(w, http.StatusNotFound, "template not found")
		return
	}

	data, err := json.MarshalIndent(tmpl, "", "  ")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to serialize template")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+name+".json\"")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// ImportTemplate handles POST /api/admin/templates/import — upload a .json file as a new template.
func (h *TemplateHandler) ImportTemplate(w http.ResponseWriter, r *http.Request) {
	// Max 10MB upload
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read file")
		return
	}

	var tmpl model.MapTemplate
	if err := json.Unmarshal(data, &tmpl); err != nil {
		writeError(w, http.StatusBadRequest, "invalid template JSON: "+err.Error())
		return
	}

	if tmpl.Name == "" {
		writeError(w, http.StatusBadRequest, "template name is required in JSON")
		return
	}

	if len(tmpl.Tiles) == 0 {
		writeError(w, http.StatusBadRequest, "template has no tiles")
		return
	}

	// Validate all tiles
	mapHalf := (tmpl.MapSize - 1) / 2
	for _, t := range tmpl.Tiles {
		if !service.ValidTerrainTypes[t.TerrainType] {
			writeError(w, http.StatusBadRequest, "invalid terrain type: "+t.TerrainType)
			return
		}
		if !model.ValidKingdomZones[t.KingdomZone] {
			writeError(w, http.StatusBadRequest, "invalid zone: "+t.KingdomZone)
			return
		}
		if t.X < -mapHalf || t.X > mapHalf || t.Y < -mapHalf || t.Y > mapHalf {
			writeError(w, http.StatusBadRequest, "tile out of bounds")
			return
		}
	}

	if err := h.templateService.SaveTemplate(&tmpl); err != nil {
		slog.Error("failed to import template", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to save imported template")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"message":  "template imported",
		"template": templateToSummary(&tmpl),
	})
}

// ResizeTemplate handles PUT /api/admin/templates/{name}/resize.
func (h *TemplateHandler) ResizeTemplate(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "template name is required")
		return
	}

	var req dto.ResizeTemplateRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if req.MapSize < 3 {
		writeError(w, http.StatusBadRequest, "map_size must be at least 3")
		return
	}

	tmpl, err := h.templateService.ResizeTemplate(name, req.MapSize)
	if err != nil {
		slog.Error("failed to resize template", "name", name, "error", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message":  "template resized",
		"template": templateToSummary(tmpl),
	})
}

// templateToSummary converts a template to a summary (without tiles).
func templateToSummary(tmpl *model.MapTemplate) map[string]any {
	return map[string]any{
		"name":        tmpl.Name,
		"description": tmpl.Description,
		"map_size":    tmpl.MapSize,
		"tile_count":  len(tmpl.Tiles),
		"created_at":  tmpl.CreatedAt.Format("2006-01-02T15:04:05Z"),
		"updated_at":  tmpl.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
