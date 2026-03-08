package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/luisfpires18/woo/internal/model"
)

// TemplateRepository manages map templates stored as JSON files on disk.
type TemplateRepository struct {
	dir string // e.g. "data/templates"
}

// NewTemplateRepository creates a new file-system backed template repository.
// It ensures the templates directory exists.
func NewTemplateRepository(dir string) (*TemplateRepository, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create templates directory: %w", err)
	}
	return &TemplateRepository{dir: dir}, nil
}

// TemplateInfo contains metadata about a template without the full tile set.
type TemplateInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	MapSize     int    `json:"map_size"`
	TileCount   int    `json:"tile_count"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// List returns metadata for all saved templates.
func (r *TemplateRepository) List() ([]TemplateInfo, error) {
	entries, err := os.ReadDir(r.dir)
	if err != nil {
		return nil, fmt.Errorf("read templates directory: %w", err)
	}

	var infos []TemplateInfo
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		tmpl, err := r.loadFile(entry.Name())
		if err != nil {
			continue // skip corrupt files
		}

		infos = append(infos, TemplateInfo{
			Name:        tmpl.Name,
			Description: tmpl.Description,
			MapSize:     tmpl.MapSize,
			TileCount:   len(tmpl.Tiles),
			CreatedAt:   tmpl.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:   tmpl.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return infos, nil
}

// Load reads a template by name.
func (r *TemplateRepository) Load(name string) (*model.MapTemplate, error) {
	filename := r.filename(name)
	return r.loadFile(filename)
}

// Save writes a template to disk.
func (r *TemplateRepository) Save(tmpl *model.MapTemplate) error {
	data, err := json.MarshalIndent(tmpl, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal template: %w", err)
	}

	path := filepath.Join(r.dir, r.filename(tmpl.Name))
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write template file: %w", err)
	}

	return nil
}

// Delete removes a template file.
func (r *TemplateRepository) Delete(name string) error {
	path := filepath.Join(r.dir, r.filename(name))
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("template %q not found", name)
		}
		return fmt.Errorf("delete template: %w", err)
	}
	return nil
}

// Exists returns true if a template with the given name exists.
func (r *TemplateRepository) Exists(name string) bool {
	path := filepath.Join(r.dir, r.filename(name))
	_, err := os.Stat(path)
	return err == nil
}

// loadFile reads and parses a template JSON file.
func (r *TemplateRepository) loadFile(filename string) (*model.MapTemplate, error) {
	path := filepath.Join(r.dir, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("template not found")
		}
		return nil, fmt.Errorf("read template file: %w", err)
	}

	var tmpl model.MapTemplate
	if err := json.Unmarshal(data, &tmpl); err != nil {
		return nil, fmt.Errorf("parse template file: %w", err)
	}

	return &tmpl, nil
}

// filename converts a template name to a safe filename.
func (r *TemplateRepository) filename(name string) string {
	// Sanitize: only allow alphanumeric, dash, underscore
	safe := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, name)

	if safe == "" {
		safe = "unnamed"
	}

	return safe + ".json"
}
