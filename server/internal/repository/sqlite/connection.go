package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// NewConnection opens a SQLite database and applies recommended PRAGMAs.
func NewConnection(dbPath string) *sql.DB {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(fmt.Sprintf("failed to create database directory: %v", err))
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		panic(fmt.Sprintf("failed to open database: %v", err))
	}

	// Enable WAL mode for concurrent reads
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		panic(fmt.Sprintf("failed to set WAL mode: %v", err))
	}

	// Enable foreign key enforcement
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		panic(fmt.Sprintf("failed to enable foreign keys: %v", err))
	}

	// Increase cache size for better performance (64 MB)
	if _, err := db.Exec("PRAGMA cache_size=-64000"); err != nil {
		panic(fmt.Sprintf("failed to set cache size: %v", err))
	}

	// Synchronous mode: NORMAL is safe with WAL
	if _, err := db.Exec("PRAGMA synchronous=NORMAL"); err != nil {
		panic(fmt.Sprintf("failed to set synchronous mode: %v", err))
	}

	// SQLite only supports 1 writer
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	return db
}

// RunMigrations applies all pending SQL migration files from the given directory.
func RunMigrations(db *sql.DB, migrationsDir string) error {
	// Create schema_migrations table if not exists
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TEXT NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	// Get current version
	var currentVersion int
	db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&currentVersion)

	// Read and apply pending migrations
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}

		version := parseVersion(f.Name())
		if version <= currentVersion {
			continue
		}

		sqlBytes, err := os.ReadFile(filepath.Join(migrationsDir, f.Name()))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", f.Name(), err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin transaction for %s: %w", f.Name(), err)
		}

		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", f.Name(), err)
		}

		tx.Exec("INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)",
			version, time.Now().UTC().Format(time.RFC3339))

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", f.Name(), err)
		}

		slog.Info("applied migration", "version", version, "file", f.Name())
	}

	return nil
}

// parseVersion extracts the numeric prefix from a migration filename (e.g. "001_create_players.sql" → 1).
func parseVersion(filename string) int {
	parts := strings.SplitN(filename, "_", 2)
	if len(parts) == 0 {
		return 0
	}
	v, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}
	return v
}
