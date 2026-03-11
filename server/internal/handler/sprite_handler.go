package handler

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// allKingdoms includes all 8 kingdoms (playable + NPC) that may have sprites.
var allKingdoms = map[string]bool{
	"veridor":   true,
	"sylvara":   true,
	"arkazia":   true,
	"draxys":    true,
	"nordalh":   true,
	"zandres":   true,
	"lumus":     true,
	"drakanith": true,
}

// validResources lists the four resource types.
var validResources = map[string]bool{
	"food":   true,
	"water":  true,
	"lumber": true,
	"stone":  true,
}

// validSlots lists valid slot numbers.
var validSlots = map[string]bool{"1": true, "2": true, "3": true}

// spriteKeyRe validates the resource_slot format: e.g. "food_1", "stone_3".
var spriteKeyRe = regexp.MustCompile(`^(food|water|lumber|stone)_([1-3])$`)

// buildingSpriteRe parses a filename like "arkazia_food_1_herdstead.png".
// Groups: 1=kingdom, 2=resource, 3=slot, 4=optional_name (without leading _).
var buildingSpriteRe = regexp.MustCompile(`^([a-z]+)_(food|water|lumber|stone)_([1-3])(?:_(.+))?\.png$`)

// SpriteHandler handles sprite-related HTTP endpoints.
type SpriteHandler struct {
	uploadsDir string // base directory for uploads (e.g. "uploads")
}

// NewSpriteHandler creates a new SpriteHandler.
// uploadsDir is the root uploads directory (e.g. "uploads").
func NewSpriteHandler(uploadsDir string) *SpriteHandler {
	return &SpriteHandler{uploadsDir: uploadsDir}
}

// SyncSpriteManifest regenerates uploads/sprites/sprites.txt with sprite paths and dimensions.
func (h *SpriteHandler) SyncSpriteManifest() error {
	return h.writeSpriteManifest()
}

// RegisterPublicRoutes registers public sprite routes on the given mux.
func (h *SpriteHandler) RegisterPublicRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/sprites/building/{kingdom}/{key}", h.ResolveBuildingSprite)
}

// RegisterAdminRoutes registers admin-only sprite routes on the given mux.
func (h *SpriteHandler) RegisterAdminRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /sprites/buildings/{kingdom}", h.ListKingdomBuildingSprites)
}

// buildingsDir returns the path to a kingdom's buildings sprite folder.
func (h *SpriteHandler) buildingsDir(kingdom string) string {
	return filepath.Join(h.uploadsDir, "sprites", "kingdoms", kingdom, "buildings")
}

// ResolveBuildingSprite handles GET /api/sprites/building/{kingdom}/{key}.
// It scans the kingdom's buildings folder for a file matching the prefix
// {kingdom}_{key}*.png and redirects to the static file URL.
// key must be in format "{resource}_{slot}" e.g. "food_1".
func (h *SpriteHandler) ResolveBuildingSprite(w http.ResponseWriter, r *http.Request) {
	_ = h.writeSpriteManifest()

	kingdom := r.PathValue("kingdom")
	key := r.PathValue("key")

	// Validate kingdom
	if !allKingdoms[kingdom] {
		writeError(w, http.StatusBadRequest, "invalid kingdom")
		return
	}

	// Validate key format
	if !spriteKeyRe.MatchString(key) {
		writeError(w, http.StatusBadRequest, "invalid key format, expected {resource}_{slot} e.g. food_1")
		return
	}

	// Build prefix to search for: e.g. "arkazia_food_1"
	prefix := kingdom + "_" + key

	dir := h.buildingsDir(kingdom)
	entries, err := os.ReadDir(dir)
	if err != nil {
		// Directory doesn't exist — no sprites for this kingdom
		http.NotFound(w, r)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, ".png") {
			// Found a match — redirect to the static file URL
			url := "/uploads/sprites/kingdoms/" + kingdom + "/buildings/" + name
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
	}

	// No matching sprite found
	http.NotFound(w, r)
}

// BuildingSpriteInfo describes a parsed building sprite file.
type BuildingSpriteInfo struct {
	Filename     string `json:"filename"`
	ResourceType string `json:"resource_type"`
	Slot         int    `json:"slot"`
	Name         string `json:"name"`
	URL          string `json:"url"`
}

// ListKingdomBuildingSprites handles GET /api/admin/sprites/buildings/{kingdom}.
// It scans the kingdom's buildings folder and returns parsed sprite metadata.
func (h *SpriteHandler) ListKingdomBuildingSprites(w http.ResponseWriter, r *http.Request) {
	_ = h.writeSpriteManifest()

	kingdom := r.PathValue("kingdom")

	if !allKingdoms[kingdom] {
		writeError(w, http.StatusBadRequest, "invalid kingdom")
		return
	}

	dir := h.buildingsDir(kingdom)
	entries, err := os.ReadDir(dir)
	if err != nil {
		// Directory doesn't exist — return empty list
		writeJSON(w, http.StatusOK, map[string]any{"sprites": []BuildingSpriteInfo{}})
		return
	}

	sprites := make([]BuildingSpriteInfo, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		matches := buildingSpriteRe.FindStringSubmatch(name)
		if matches == nil {
			continue
		}

		// matches: [full, kingdom, resource, slot, optional_name]
		fileKingdom := matches[1]
		if fileKingdom != kingdom {
			continue // skip files that don't match the expected kingdom prefix
		}

		resource := matches[2]
		slotStr := matches[3]
		slot := 0
		switch slotStr {
		case "1":
			slot = 1
		case "2":
			slot = 2
		case "3":
			slot = 3
		}

		spriteName := matches[4] // may be empty if no name suffix

		sprites = append(sprites, BuildingSpriteInfo{
			Filename:     name,
			ResourceType: resource,
			Slot:         slot,
			Name:         spriteName,
			URL:          "/uploads/sprites/kingdoms/" + kingdom + "/buildings/" + name,
		})
	}

	writeJSON(w, http.StatusOK, map[string]any{"sprites": sprites})
}

func (h *SpriteHandler) writeSpriteManifest() error {
	root := filepath.Join(h.uploadsDir, "sprites")
	if err := os.MkdirAll(root, 0o755); err != nil {
		return fmt.Errorf("mkdir sprites root: %w", err)
	}

	kingdoms := []string{"arkazia", "draxys", "drakanith", "lumus", "nordalh", "sylvara", "veridor", "zandres"}
	resources := []string{"food", "water", "lumber", "stone"}
	slots := []string{"1", "2", "3"}
	militaryBuildings := []string{"barracks", "stable", "archery", "workshop", "special"}

	actualLines := make([]string, 0)
	templateLines := make([]string, 0)

	// Track which kingdom_resource_slot combinations exist
	existingCombos := make(map[string]bool)

	// Walk directory to find ACTUAL sprites
	walkErr := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if d.IsDir() || strings.ToLower(d.Name()) == "sprites.txt" {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" {
			relPath, _ := filepath.Rel(h.uploadsDir, path)
			fullPath := filepath.ToSlash(filepath.Join("uploads", relPath))
			actualLines = append(actualLines, fmt.Sprintf("%s | 256x256 (recommended)", fullPath))

			// Track kingdom_resource_slot combinations for building sprites
			filename := d.Name()
			matches := buildingSpriteRe.FindStringSubmatch(filename)
			if matches != nil {
				// matches: [full, kingdom, resource, slot, optional_name]
				kingdom := matches[1]
				resource := matches[2]
				slot := matches[3]
				combo := fmt.Sprintf("%s_%s_%s", kingdom, resource, slot)
				existingCombos[combo] = true
			}
		}
		return nil
	})
	_ = walkErr

	// Resource type sprites (template only if missing)
	for _, resource := range resources {
		path := fmt.Sprintf("uploads/sprites/resources/%s.png", resource)
		found := false
		for _, line := range actualLines {
			if strings.Contains(line, path) {
				found = true
				break
			}
		}
		if !found {
			templateLines = append(templateLines, fmt.Sprintf("%s | 256x256 (recommended)", path))
		}
	}

	// Kingdom RESOURCE building sprites (template only if the combination doesn't exist)
	for _, kingdom := range kingdoms {
		for _, resource := range resources {
			for _, slot := range slots {
				combo := fmt.Sprintf("%s_%s_%s", kingdom, resource, slot)
				if !existingCombos[combo] {
					path := fmt.Sprintf("uploads/sprites/kingdoms/%s/buildings/%s_%s_%s_[name].png", kingdom, kingdom, resource, slot)
					templateLines = append(templateLines, fmt.Sprintf("%s | 256x256 (recommended)", path))
				}
			}
		}
	}

	// Kingdom MILITARY building sprites (barracks, stable, archery, workshop, special)
	for _, kingdom := range kingdoms {
		for _, building := range militaryBuildings {
			path := fmt.Sprintf("uploads/sprites/kingdoms/%s/buildings/%s_%s.png", kingdom, kingdom, building)
			found := false
			for _, line := range actualLines {
				if strings.Contains(line, path) {
					found = true
					break
				}
			}
			if !found {
				templateLines = append(templateLines, fmt.Sprintf("%s | 256x256 (recommended)", path))
			}
		}
	}

	// Troop unit sprites (approximate 20 troops per kingdom)
	// Format: uploads/sprites/kingdoms/{kingdom}/troops/{kingdom}_{troop_name}.png
	for _, kingdom := range kingdoms {
		// Each kingdom has approximately 4-5 troops per building across 5 buildings = 20 total
		for _, building := range militaryBuildings {
			for i := 1; i <= 5; i++ {
				path := fmt.Sprintf("uploads/sprites/kingdoms/%s/troops/%s_%s_[troop_name]_%d.png", kingdom, kingdom, building, i)
				templateLines = append(templateLines, fmt.Sprintf("%s | 256x256 (recommended)", path))
			}
		}
	}

	content := "# Sprite Manifest\n" +
		"# This file lists all sprite locations with RECOMMENDED sizes.\n" +
		"# Format: <path> | <recommended-size>\n" +
		"# Paths marked [name] are optional display name suffixes (e.g. arkazia_food_1_herdstead.png).\n" +
		"# Paths marked [troop_name] are troop unit names (e.g. sylvara_barracks_rootguard_spearmen.png).\n" +
		"# Auto-generated by server on startup.\n\n"

	if len(actualLines) > 0 {
		content += "## EXISTING SPRITES\n\n"
		content += strings.Join(actualLines, "\n") + "\n\n"
	}

	if len(templateLines) > 0 {
		content += "## TEMPLATE (for future sprites)\n\n"
		content += strings.Join(templateLines, "\n") + "\n"
	}

	manifestPath := filepath.Join(root, "sprites.txt")
	if err := os.WriteFile(manifestPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write sprites.txt: %w", err)
	}

	return nil
}
