package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/handler"
	"github.com/luisfpires18/woo/internal/repository/sqlite"
	"github.com/luisfpires18/woo/internal/service"
	"github.com/luisfpires18/woo/internal/testutil"
)

// apiEnvelope mirrors handler.apiResponse for test decoding.
type apiEnvelope struct {
	Data  json.RawMessage `json:"data,omitempty"`
	Error string          `json:"error,omitempty"`
}

// newTestMapHandler creates a MapHandler backed by an in-memory DB with a generated world map.
func newTestMapHandler(t *testing.T) *handler.MapHandler {
	t.Helper()
	db := testutil.NewTestDB(t)
	worldMapRepo := sqlite.NewWorldMapRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)
	mapService := service.NewMapService(worldMapRepo, villageRepo)

	if err := mapService.GenerateMap(context.Background()); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	return handler.NewMapHandler(mapService)
}

// decodeChunkResponse decodes the apiResponse envelope and extracts a MapChunkResponse.
func decodeChunkResponse(t *testing.T, rec *httptest.ResponseRecorder) dto.MapChunkResponse {
	t.Helper()
	var env apiEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&env); err != nil {
		t.Fatalf("decode envelope: %v", err)
	}
	var resp dto.MapChunkResponse
	if err := json.Unmarshal(env.Data, &resp); err != nil {
		t.Fatalf("decode MapChunkResponse: %v", err)
	}
	return resp
}

// decodeTileResponse decodes the apiResponse envelope and extracts a MapTileInfo.
func decodeTileResponse(t *testing.T, rec *httptest.ResponseRecorder) dto.MapTileInfo {
	t.Helper()
	var env apiEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&env); err != nil {
		t.Fatalf("decode envelope: %v", err)
	}
	var tile dto.MapTileInfo
	if err := json.Unmarshal(env.Data, &tile); err != nil {
		t.Fatalf("decode MapTileInfo: %v", err)
	}
	return tile
}

func TestGetMapChunk_DefaultParams(t *testing.T) {
	h := newTestMapHandler(t)

	req := httptest.NewRequest("GET", "/api/map", nil)
	rec := httptest.NewRecorder()

	h.GetMapChunk(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	resp := decodeChunkResponse(t, rec)

	// Default center is (0,0), default range is 10
	if resp.CenterX != 0 || resp.CenterY != 0 {
		t.Errorf("center: got (%d,%d), want (0,0)", resp.CenterX, resp.CenterY)
	}
	if resp.Range != 10 {
		t.Errorf("range: got %d, want 10", resp.Range)
	}

	// 10 radius → 21×21 = 441 tiles
	if len(resp.Tiles) != 441 {
		t.Errorf("tile count: got %d, want 441", len(resp.Tiles))
	}
}

func TestGetMapChunk_WithCoordinates(t *testing.T) {
	h := newTestMapHandler(t)

	req := httptest.NewRequest("GET", "/api/map?x=10&y=-10&range=3", nil)
	rec := httptest.NewRecorder()

	h.GetMapChunk(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	resp := decodeChunkResponse(t, rec)

	if resp.CenterX != 10 || resp.CenterY != -10 {
		t.Errorf("center: got (%d,%d), want (10,-10)", resp.CenterX, resp.CenterY)
	}
	if resp.Range != 3 {
		t.Errorf("range: got %d, want 3", resp.Range)
	}

	// 3 radius → 7×7 = 49 tiles
	if len(resp.Tiles) != 49 {
		t.Errorf("tile count: got %d, want 49", len(resp.Tiles))
	}
}

func TestGetMapChunk_RangeClampedToMax(t *testing.T) {
	h := newTestMapHandler(t)

	req := httptest.NewRequest("GET", "/api/map?x=0&y=0&range=999", nil)
	rec := httptest.NewRecorder()

	h.GetMapChunk(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	resp := decodeChunkResponse(t, rec)

	// Should clamp to 40
	if resp.Range != 40 {
		t.Errorf("range: got %d, want 40 (clamped)", resp.Range)
	}
}

func TestGetMapChunk_TilesHaveRequiredFields(t *testing.T) {
	h := newTestMapHandler(t)

	req := httptest.NewRequest("GET", "/api/map?x=0&y=0&range=2", nil)
	rec := httptest.NewRecorder()

	h.GetMapChunk(rec, req)

	resp := decodeChunkResponse(t, rec)

	for _, tile := range resp.Tiles {
		if tile.Terrain == "" {
			t.Errorf("tile (%d,%d) has empty terrain", tile.X, tile.Y)
		}
	}
}

func TestGetTile_Success(t *testing.T) {
	h := newTestMapHandler(t)

	req := httptest.NewRequest("GET", "/api/map/tile?x=0&y=0", nil)
	rec := httptest.NewRecorder()

	h.GetTile(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	tile := decodeTileResponse(t, rec)

	if tile.X != 0 || tile.Y != 0 {
		t.Errorf("coords: got (%d,%d), want (0,0)", tile.X, tile.Y)
	}
	if tile.Terrain == "" {
		t.Error("expected non-empty terrain")
	}
}

func TestGetTile_MissingParams(t *testing.T) {
	h := newTestMapHandler(t)

	tests := []struct {
		name string
		url  string
	}{
		{"missing x", "/api/map/tile?y=0"},
		{"missing y", "/api/map/tile?x=0"},
		{"missing both", "/api/map/tile"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.url, nil)
			rec := httptest.NewRecorder()

			h.GetTile(rec, req)

			if rec.Code == http.StatusOK {
				t.Errorf("expected error status for %s, got 200", tt.url)
			}
		})
	}
}

func TestGetTile_NotFound(t *testing.T) {
	h := newTestMapHandler(t)

	req := httptest.NewRequest("GET", "/api/map/tile?x=999&y=999", nil)
	rec := httptest.NewRecorder()

	h.GetTile(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}
