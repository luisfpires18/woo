package testutil

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

// NewTestDB creates an in-memory SQLite database with all migrations applied.
func NewTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	// Enable foreign keys
	db.Exec("PRAGMA foreign_keys=ON")

	// Find migrations directory: walk up from cwd until we find server/migrations
	migrationsDir := findMigrationsDir(t)

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("read migrations dir %s: %v", migrationsDir, err)
	}

	// Sort files to ensure order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(migrationsDir, f.Name()))
		if err != nil {
			t.Fatalf("read migration %s: %v", f.Name(), err)
		}
		if _, err := db.Exec(string(content)); err != nil {
			t.Fatalf("apply migration %s: %v", f.Name(), err)
		}
	}

	return db
}

// findMigrationsDir walks up directories from cwd to find server/migrations.
func findMigrationsDir(t *testing.T) string {
	t.Helper()

	// Try from environment variable first
	if dir := os.Getenv("MIGRATIONS_PATH"); dir != "" {
		return dir
	}

	// Walk up from current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}

	dir := cwd
	for {
		candidate := filepath.Join(dir, "migrations")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	t.Fatalf("could not find migrations directory starting from %s", cwd)
	return ""
}
