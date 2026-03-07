package handler

import (
	"net/http"
	"strconv"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/service"
)

// MapHandler handles world map HTTP endpoints.
type MapHandler struct {
	mapService *service.MapService
}

// NewMapHandler creates a new MapHandler.
func NewMapHandler(mapService *service.MapService) *MapHandler {
	return &MapHandler{mapService: mapService}
}

// RegisterRoutes registers map routes on the given mux.
func (h *MapHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/map", h.GetMapChunk)
	mux.HandleFunc("GET /api/map/tile", h.GetTile)
}

// GetMapChunk handles GET /api/map?x={x}&y={y}&range={r}
// Returns a grid of tiles centered on (x, y) with the given range.
func (h *MapHandler) GetMapChunk(w http.ResponseWriter, r *http.Request) {
	cx, err := strconv.Atoi(r.URL.Query().Get("x"))
	if err != nil {
		cx = 0
	}
	cy, err := strconv.Atoi(r.URL.Query().Get("y"))
	if err != nil {
		cy = 0
	}
	radius, err := strconv.Atoi(r.URL.Query().Get("range"))
	if err != nil || radius < 1 {
		radius = 10
	}
	if radius > 40 {
		radius = 40
	}

	tiles, err := h.mapService.GetMapChunk(r.Context(), cx, cy, radius)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load map")
		return
	}

	tileInfos := make([]dto.MapTileInfo, 0, len(tiles))
	for _, t := range tiles {
		tileInfos = append(tileInfos, dto.MapTileFromModel(t))
	}

	resp := dto.MapChunkResponse{
		CenterX: cx,
		CenterY: cy,
		Range:   radius,
		Tiles:   tileInfos,
	}

	writeJSON(w, http.StatusOK, resp)
}

// GetTile handles GET /api/map/tile?x={x}&y={y}
// Returns a single tile's details.
func (h *MapHandler) GetTile(w http.ResponseWriter, r *http.Request) {
	x, err := strconv.Atoi(r.URL.Query().Get("x"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid x coordinate")
		return
	}
	y, err := strconv.Atoi(r.URL.Query().Get("y"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid y coordinate")
		return
	}

	tile, err := h.mapService.GetTile(r.Context(), x, y)
	if err != nil {
		writeError(w, http.StatusNotFound, "tile not found")
		return
	}

	writeJSON(w, http.StatusOK, dto.MapTileFromModel(tile))
}
