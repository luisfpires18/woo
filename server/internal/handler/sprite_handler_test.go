package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/luisfpires18/woo/internal/handler"
)

// setupSpriteDir creates a temp uploads dir with sprite files for testing.
func setupSpriteDir(t *testing.T, kingdom string, filenames []string) string {
	t.Helper()
	uploadsDir := t.TempDir()
	dir := filepath.Join(uploadsDir, "sprites", "kingdoms", kingdom, "buildings")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	for _, name := range filenames {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("fake-png"), 0o644); err != nil {
			t.Fatalf("WriteFile %s: %v", name, err)
		}
	}
	return uploadsDir
}

func TestResolveBuildingSprite_Found(t *testing.T) {
	uploadsDir := setupSpriteDir(t, "arkazia", []string{
		"arkazia_food_1_herdstead.png",
		"arkazia_food_2_slopefarm.png",
		"arkazia_stone_3_deepmine.png",
	})

	h := handler.NewSpriteHandler(uploadsDir)
	mux := http.NewServeMux()
	h.RegisterPublicRoutes(mux)

	req := httptest.NewRequest("GET", "/api/sprites/building/arkazia/food_1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	expected := "/uploads/sprites/kingdoms/arkazia/buildings/arkazia_food_1_herdstead.png"
	if loc != expected {
		t.Errorf("expected redirect to %s, got %s", expected, loc)
	}
}

func TestResolveBuildingSprite_NotFound(t *testing.T) {
	uploadsDir := setupSpriteDir(t, "arkazia", []string{
		"arkazia_food_1_herdstead.png",
	})

	h := handler.NewSpriteHandler(uploadsDir)
	mux := http.NewServeMux()
	h.RegisterPublicRoutes(mux)

	req := httptest.NewRequest("GET", "/api/sprites/building/arkazia/food_3", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestResolveBuildingSprite_InvalidKingdom(t *testing.T) {
	uploadsDir := t.TempDir()
	h := handler.NewSpriteHandler(uploadsDir)
	mux := http.NewServeMux()
	h.RegisterPublicRoutes(mux)

	req := httptest.NewRequest("GET", "/api/sprites/building/invalid_kingdom/food_1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestResolveBuildingSprite_InvalidKey(t *testing.T) {
	uploadsDir := t.TempDir()
	h := handler.NewSpriteHandler(uploadsDir)
	mux := http.NewServeMux()
	h.RegisterPublicRoutes(mux)

	req := httptest.NewRequest("GET", "/api/sprites/building/arkazia/invalid", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestResolveBuildingSprite_MissingKingdomDir(t *testing.T) {
	uploadsDir := t.TempDir()
	h := handler.NewSpriteHandler(uploadsDir)
	mux := http.NewServeMux()
	h.RegisterPublicRoutes(mux)

	req := httptest.NewRequest("GET", "/api/sprites/building/sylvara/food_1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestListKingdomBuildingSprites_Success(t *testing.T) {
	uploadsDir := setupSpriteDir(t, "arkazia", []string{
		"arkazia_food_1_herdstead.png",
		"arkazia_water_2_dam.png",
		"arkazia_stone_3_deepmine.png",
		"readme.txt", // should be ignored
	})

	h := handler.NewSpriteHandler(uploadsDir)
	mux := http.NewServeMux()
	h.RegisterAdminRoutes(mux)

	req := httptest.NewRequest("GET", "/sprites/buildings/arkazia", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	// Decode the envelope
	var env struct {
		Data struct {
			Sprites []handler.BuildingSpriteInfo `json:"sprites"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&env); err != nil {
		t.Fatalf("decode: %v", err)
	}

	sprites := env.Data.Sprites
	if len(sprites) != 3 {
		t.Fatalf("expected 3 sprites, got %d", len(sprites))
	}

	// Check that all expected sprites are present
	found := map[string]bool{}
	for _, s := range sprites {
		found[s.Filename] = true
		if s.URL == "" {
			t.Errorf("sprite %s has empty URL", s.Filename)
		}
	}
	for _, expected := range []string{
		"arkazia_food_1_herdstead.png",
		"arkazia_water_2_dam.png",
		"arkazia_stone_3_deepmine.png",
	} {
		if !found[expected] {
			t.Errorf("missing expected sprite: %s", expected)
		}
	}
}

func TestListKingdomBuildingSprites_EmptyDir(t *testing.T) {
	uploadsDir := t.TempDir()
	h := handler.NewSpriteHandler(uploadsDir)
	mux := http.NewServeMux()
	h.RegisterAdminRoutes(mux)

	req := httptest.NewRequest("GET", "/sprites/buildings/veridor", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var env struct {
		Data struct {
			Sprites []handler.BuildingSpriteInfo `json:"sprites"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&env); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(env.Data.Sprites) != 0 {
		t.Errorf("expected 0 sprites, got %d", len(env.Data.Sprites))
	}
}

func TestListKingdomBuildingSprites_InvalidKingdom(t *testing.T) {
	uploadsDir := t.TempDir()
	h := handler.NewSpriteHandler(uploadsDir)
	mux := http.NewServeMux()
	h.RegisterAdminRoutes(mux)

	req := httptest.NewRequest("GET", "/sprites/buildings/hackme", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestListKingdomBuildingSprites_IgnoresWrongKingdomPrefix(t *testing.T) {
	uploadsDir := setupSpriteDir(t, "arkazia", []string{
		"arkazia_food_1_farm.png",
		"sylvara_food_1_grove.png", // wrong kingdom prefix — should be ignored
	})

	h := handler.NewSpriteHandler(uploadsDir)
	mux := http.NewServeMux()
	h.RegisterAdminRoutes(mux)

	req := httptest.NewRequest("GET", "/sprites/buildings/arkazia", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var env struct {
		Data struct {
			Sprites []handler.BuildingSpriteInfo `json:"sprites"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&env); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(env.Data.Sprites) != 1 {
		t.Fatalf("expected 1 sprite, got %d", len(env.Data.Sprites))
	}
	if env.Data.Sprites[0].Filename != "arkazia_food_1_farm.png" {
		t.Errorf("unexpected filename: %s", env.Data.Sprites[0].Filename)
	}
}

func TestSyncSpriteManifest_GeneratesTemplate(t *testing.T) {
	uploadsDir := t.TempDir()

	h := handler.NewSpriteHandler(uploadsDir)
	if err := h.SyncSpriteManifest(); err != nil {
		t.Fatalf("SyncSpriteManifest: %v", err)
	}

	manifestPath := filepath.Join(uploadsDir, "sprites", "sprites.txt")
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("ReadFile sprites.txt: %v", err)
	}

	text := string(content)

	// Check header
	if !strings.Contains(text, "# Sprite Manifest — Template & Guideline") {
		t.Fatalf("missing header in manifest")
	}
	if !strings.Contains(text, "RECOMMENDED") {
		t.Fatalf("missing 'RECOMMENDED' in manifest")
	}

	// Check all 4 resource types are listed with 256x256 recommended
	for _, resource := range []string{"food", "water", "lumber", "stone"} {
		expected := "uploads/sprites/resources/" + resource + ".png | 256x256"
		if !strings.Contains(text, expected) {
			t.Errorf("missing resource sprite entry: %s", expected)
		}
	}

	// Check sample building sprites from each kingdom
	expectedKingdoms := []string{"arkazia", "sylvara", "veridor"}
	for _, kingdom := range expectedKingdoms {
		expected := "uploads/sprites/kingdoms/" + kingdom + "/buildings/" + kingdom + "_food_1_"
		if !strings.Contains(text, expected) {
			t.Errorf("missing building sprite entry for kingdom: %s", kingdom)
		}
	}

	// Should have 8 kingdoms × 4 resources × 3 slots = 96 building entries + 4 resource entries = 100 entries
	entries := strings.Count(text, "| 256x256")
	if entries < 100 {
		t.Errorf("expected at least 100 sprite entries (96 buildings + 4 resources), got %d", entries)
	}
}
