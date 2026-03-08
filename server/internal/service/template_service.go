package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
)

// Template service errors.
var (
	ErrTemplateNotFound    = fmt.Errorf("template not found")
	ErrTemplateExists      = fmt.Errorf("template already exists")
	ErrInvalidZone         = fmt.Errorf("invalid kingdom zone")
	ErrInvalidTemplateTile = fmt.Errorf("invalid template tile")
)

// TemplateService manages map templates — creating, editing, saving, and applying to the live game.
type TemplateService struct {
	templateRepo *repository.TemplateRepository
	mapRepo      repository.WorldMapRepository
}

// NewTemplateService creates a new TemplateService.
func NewTemplateService(templateRepo *repository.TemplateRepository, mapRepo repository.WorldMapRepository) *TemplateService {
	return &TemplateService{
		templateRepo: templateRepo,
		mapRepo:      mapRepo,
	}
}

// ListTemplates returns metadata for all saved templates.
func (s *TemplateService) ListTemplates() ([]repository.TemplateInfo, error) {
	return s.templateRepo.List()
}

// GetTemplate loads a full template by name (including all tiles).
func (s *TemplateService) GetTemplate(name string) (*model.MapTemplate, error) {
	tmpl, err := s.templateRepo.Load(name)
	if err != nil {
		return nil, fmt.Errorf("load template %q: %w", name, err)
	}
	return tmpl, nil
}

// CreateTemplate creates a new blank template with all-plains/wilderness tiles.
// If mapSize is 0, defaults to model.MapSize (51).
func (s *TemplateService) CreateTemplate(name, description string, mapSize int) (*model.MapTemplate, error) {
	if s.templateRepo.Exists(name) {
		return nil, ErrTemplateExists
	}

	if mapSize <= 0 {
		mapSize = model.MapSize
	}
	if mapSize%2 == 0 {
		mapSize++ // ensure odd
	}
	if mapSize < 3 {
		mapSize = 3
	}
	if mapSize > 201 {
		return nil, fmt.Errorf("map_size too large (max 201)")
	}

	mapHalf := (mapSize - 1) / 2
	now := time.Now().UTC()

	tiles := make([]model.TemplateTile, 0, mapSize*mapSize)
	for y := -mapHalf; y <= mapHalf; y++ {
		for x := -mapHalf; x <= mapHalf; x++ {
			tiles = append(tiles, model.TemplateTile{
				X:           x,
				Y:           y,
				TerrainType: model.TerrainPlains,
				KingdomZone: model.ZoneWilderness,
			})
		}
	}

	tmpl := &model.MapTemplate{
		Name:        name,
		Description: description,
		MapSize:     mapSize,
		Tiles:       tiles,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.templateRepo.Save(tmpl); err != nil {
		return nil, fmt.Errorf("save new template: %w", err)
	}

	slog.Info("map template created", "name", name, "size", mapSize, "tiles", len(tiles))
	return tmpl, nil
}

// ResizeTemplate changes the size of an existing template.
// Tiles within the new bounds are preserved; new tiles are plains/wilderness; out-of-bounds tiles are dropped.
func (s *TemplateService) ResizeTemplate(name string, newSize int) (*model.MapTemplate, error) {
	tmpl, err := s.templateRepo.Load(name)
	if err != nil {
		return nil, fmt.Errorf("load template: %w", err)
	}

	if newSize%2 == 0 {
		newSize++ // ensure odd
	}
	if newSize < 3 {
		newSize = 3
	}
	if newSize > 201 {
		return nil, fmt.Errorf("map_size too large (max 201)")
	}

	newHalf := (newSize - 1) / 2

	// Build lookup of existing tiles
	existing := make(map[[2]int]model.TemplateTile, len(tmpl.Tiles))
	for _, t := range tmpl.Tiles {
		existing[[2]int{t.X, t.Y}] = t
	}

	// Generate the new tile grid
	newTiles := make([]model.TemplateTile, 0, newSize*newSize)
	for y := -newHalf; y <= newHalf; y++ {
		for x := -newHalf; x <= newHalf; x++ {
			if t, ok := existing[[2]int{x, y}]; ok {
				newTiles = append(newTiles, t)
			} else {
				newTiles = append(newTiles, model.TemplateTile{
					X:           x,
					Y:           y,
					TerrainType: model.TerrainPlains,
					KingdomZone: model.ZoneWilderness,
				})
			}
		}
	}

	tmpl.MapSize = newSize
	tmpl.Tiles = newTiles
	tmpl.UpdatedAt = time.Now().UTC()

	if err := s.templateRepo.Save(tmpl); err != nil {
		return nil, fmt.Errorf("save resized template: %w", err)
	}

	slog.Info("map template resized", "name", name, "new_size", newSize, "tiles", len(newTiles))
	return tmpl, nil
}

// DeleteTemplate removes a template by name.
func (s *TemplateService) DeleteTemplate(name string) error {
	return s.templateRepo.Delete(name)
}

// UpdateTemplateTerrain applies terrain changes to a template (same as admin paint tool but on template, not live map).
func (s *TemplateService) UpdateTemplateTerrain(name string, updates []model.TileTerrainUpdate) error {
	tmpl, err := s.templateRepo.Load(name)
	if err != nil {
		return fmt.Errorf("load template: %w", err)
	}

	if len(updates) == 0 {
		return fmt.Errorf("no tiles provided")
	}
	if len(updates) > 2601 {
		return fmt.Errorf("too many tiles (max %d)", model.MapSize*model.MapSize)
	}

	// Build a quick lookup of the update set
	mapHalf := (tmpl.MapSize - 1) / 2
	updateMap := make(map[[2]int]string, len(updates))
	for _, u := range updates {
		if !ValidTerrainTypes[u.TerrainType] {
			return fmt.Errorf("invalid terrain type: %s", u.TerrainType)
		}
		if u.X < -mapHalf || u.X > mapHalf || u.Y < -mapHalf || u.Y > mapHalf {
			return fmt.Errorf("tile (%d,%d) out of map bounds", u.X, u.Y)
		}
		updateMap[[2]int{u.X, u.Y}] = u.TerrainType
	}

	// Apply
	for i, t := range tmpl.Tiles {
		if newTerrain, ok := updateMap[[2]int{t.X, t.Y}]; ok {
			tmpl.Tiles[i].TerrainType = newTerrain
		}
	}

	tmpl.UpdatedAt = time.Now().UTC()
	return s.templateRepo.Save(tmpl)
}

// UpdateTemplateZones applies zone changes to a template.
func (s *TemplateService) UpdateTemplateZones(name string, updates []model.TileZoneUpdate) error {
	tmpl, err := s.templateRepo.Load(name)
	if err != nil {
		return fmt.Errorf("load template: %w", err)
	}

	if len(updates) == 0 {
		return fmt.Errorf("no tiles provided")
	}
	if len(updates) > 2601 {
		return fmt.Errorf("too many tiles (max %d)", model.MapSize*model.MapSize)
	}

	mapHalf := (tmpl.MapSize - 1) / 2
	updateMap := make(map[[2]int]string, len(updates))
	for _, u := range updates {
		if !model.ValidKingdomZones[u.KingdomZone] {
			return fmt.Errorf("invalid kingdom zone: %s", u.KingdomZone)
		}
		if u.X < -mapHalf || u.X > mapHalf || u.Y < -mapHalf || u.Y > mapHalf {
			return fmt.Errorf("tile (%d,%d) out of map bounds", u.X, u.Y)
		}
		updateMap[[2]int{u.X, u.Y}] = u.KingdomZone
	}

	for i, t := range tmpl.Tiles {
		if newZone, ok := updateMap[[2]int{t.X, t.Y}]; ok {
			tmpl.Tiles[i].KingdomZone = newZone
		}
	}

	tmpl.UpdatedAt = time.Now().UTC()
	return s.templateRepo.Save(tmpl)
}

// SaveTemplate overwrites a template with new data (full replacement).
func (s *TemplateService) SaveTemplate(tmpl *model.MapTemplate) error {
	tmpl.UpdatedAt = time.Now().UTC()
	return s.templateRepo.Save(tmpl)
}

// ApplyTemplate applies a saved template to the live world_map table.
// It clears the existing map and inserts all tiles from the template.
// Owner/village data is NOT preserved — this is a full map reset.
func (s *TemplateService) ApplyTemplate(ctx context.Context, name string) error {
	tmpl, err := s.templateRepo.Load(name)
	if err != nil {
		return fmt.Errorf("load template for apply: %w", err)
	}

	if len(tmpl.Tiles) == 0 {
		return fmt.Errorf("template has no tiles")
	}

	// Clear the existing world map (handles size mismatches cleanly).
	if err := s.mapRepo.DeleteAll(ctx); err != nil {
		return fmt.Errorf("clear world map: %w", err)
	}

	// Insert all tiles from the template.
	mapTiles := make([]*model.MapTile, 0, len(tmpl.Tiles))
	for _, t := range tmpl.Tiles {
		mapTiles = append(mapTiles, &model.MapTile{
			X:           t.X,
			Y:           t.Y,
			TerrainType: t.TerrainType,
			KingdomZone: t.KingdomZone,
		})
	}
	if err := s.mapRepo.InsertBatch(ctx, mapTiles); err != nil {
		return fmt.Errorf("insert template tiles: %w", err)
	}

	slog.Info("template applied to live map", "name", name, "tiles", len(mapTiles))
	return nil
}
