package handler

import (
	"errors"
	"fmt"
	"image/png"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/middleware"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/service"
)

// AdminHandler handles admin HTTP endpoints.
type AdminHandler struct {
	adminService *service.AdminService
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// RegisterRoutes registers admin routes on the given mux.
// All admin routes require auth + admin middleware applied externally.
func (h *AdminHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /players", h.ListPlayers)
	mux.HandleFunc("PATCH /players/{id}/role", h.UpdatePlayerRole)
	mux.HandleFunc("GET /config", h.GetWorldConfig)
	mux.HandleFunc("PUT /config/{key}", h.SetWorldConfig)
	mux.HandleFunc("GET /stats", h.GetStats)
	mux.HandleFunc("GET /announcements", h.ListAnnouncements)
	mux.HandleFunc("POST /announcements", h.CreateAnnouncement)
	mux.HandleFunc("DELETE /announcements/{id}", h.DeleteAnnouncement)
	mux.HandleFunc("GET /assets", h.ListAssets)
	mux.HandleFunc("POST /assets/{id}/sprite", h.UploadSprite)
	mux.HandleFunc("DELETE /assets/{id}/sprite", h.DeleteSprite)
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

// GetWorldConfig handles GET /api/admin/config.
func (h *AdminHandler) GetWorldConfig(w http.ResponseWriter, r *http.Request) {
	resp, err := h.adminService.GetWorldConfig(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get config")
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// SetWorldConfig handles PUT /api/admin/config/{key}.
func (h *AdminHandler) SetWorldConfig(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		writeError(w, http.StatusBadRequest, "missing config key")
		return
	}

	var req dto.SetConfigRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if err := h.adminService.SetWorldConfig(r.Context(), key, req.Value); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "config key not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to set config")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "config updated"})
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

// UploadSprite handles POST /api/admin/assets/{id}/sprite.
// Expects multipart/form-data with a "file" field containing a PNG image.
func (h *AdminHandler) UploadSprite(w http.ResponseWriter, r *http.Request) {
	assetID := r.PathValue("id")
	if assetID == "" {
		writeError(w, http.StatusBadRequest, "missing asset id")
		return
	}

	// Fetch the asset to validate it exists and get category.
	asset, err := h.adminService.GetGameAsset(r.Context(), assetID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "asset not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get asset")
		return
	}

	// Get expected dimensions and max size for this category.
	expectedDims, ok := model.AssetSpriteDimensions[asset.Category]
	if !ok {
		writeError(w, http.StatusBadRequest, "unknown asset category")
		return
	}
	maxBytes, _ := model.AssetMaxSpriteBytes[asset.Category]

	// Parse multipart (limit to maxBytes + some overhead for headers).
	if err := r.ParseMultipartForm(maxBytes + 4096); err != nil {
		writeError(w, http.StatusBadRequest, "file too large or invalid multipart")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "missing file field")
		return
	}
	defer file.Close()

	// Validate content type.
	ct := header.Header.Get("Content-Type")
	if ct != "image/png" {
		writeError(w, http.StatusBadRequest, "only PNG images are accepted")
		return
	}

	// Validate file size.
	if header.Size > maxBytes {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("file exceeds max size of %d KB", maxBytes/1024))
		return
	}

	// Decode the PNG to validate dimensions.
	img, err := png.Decode(file)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid PNG image")
		return
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width != expectedDims[0] || height != expectedDims[1] {
		writeError(w, http.StatusBadRequest, fmt.Sprintf(
			"image must be %dx%d pixels, got %dx%d",
			expectedDims[0], expectedDims[1], width, height,
		))
		return
	}

	// Save file to disk.
	relDir := filepath.Join("sprites", asset.Category)
	absDir := filepath.Join("uploads", relDir)
	if err := os.MkdirAll(absDir, 0o755); err != nil {
		slog.Error("failed to create sprite directory", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	filename := assetID + ".png"
	absPath := filepath.Join(absDir, filename)

	// Seek back to start since png.Decode consumed the reader.
	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	} else {
		writeError(w, http.StatusInternalServerError, "cannot re-read file")
		return
	}

	dst, err := os.Create(absPath)
	if err != nil {
		slog.Error("failed to create sprite file", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		slog.Error("failed to write sprite file", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to save file")
		return
	}

	// Use forward slashes in DB path for URL compatibility.
	dbPath := relDir + "/" + filename
	if err := h.adminService.UpdateGameAssetSprite(r.Context(), assetID, &dbPath); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update asset")
		return
	}

	spriteURL := "/uploads/" + dbPath
	writeJSON(w, http.StatusOK, map[string]string{
		"message":    "sprite uploaded",
		"sprite_url": spriteURL,
	})
}

// DeleteSprite handles DELETE /api/admin/assets/{id}/sprite.
func (h *AdminHandler) DeleteSprite(w http.ResponseWriter, r *http.Request) {
	assetID := r.PathValue("id")
	if assetID == "" {
		writeError(w, http.StatusBadRequest, "missing asset id")
		return
	}

	// Fetch asset to get the sprite path for file deletion.
	asset, err := h.adminService.GetGameAsset(r.Context(), assetID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			writeError(w, http.StatusNotFound, "asset not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get asset")
		return
	}

	// Delete file from disk if it exists.
	if asset.SpritePath != nil {
		absPath := filepath.Join("uploads", *asset.SpritePath)
		os.Remove(absPath) // ignore error — file may not exist
	}

	// Clear sprite_path in DB.
	if err := h.adminService.UpdateGameAssetSprite(r.Context(), assetID, nil); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update asset")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "sprite deleted"})
}
