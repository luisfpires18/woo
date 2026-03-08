package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
)

type seasonRepo struct {
	db *sql.DB
}

// NewSeasonRepo creates a new SQLite-backed SeasonRepository.
func NewSeasonRepo(db *sql.DB) repository.SeasonRepository {
	return &seasonRepo{db: db}
}

func (r *seasonRepo) Create(ctx context.Context, season *model.Season) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO seasons (name, description, status, start_date, map_template_name,
		 game_speed, resource_multiplier, max_villages_per_player, weapons_of_chaos_count,
		 map_width, map_height)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		season.Name, season.Description, season.Status, season.StartDate,
		season.MapTemplateName, season.GameSpeed, season.ResourceMultiplier,
		season.MaxVillagesPerPlayer, season.WeaponsOfChaosCount,
		season.MapWidth, season.MapHeight,
	)
	if err != nil {
		return fmt.Errorf("insert season: %w", err)
	}
	id, _ := result.LastInsertId()
	season.ID = id
	return nil
}

func (r *seasonRepo) GetByID(ctx context.Context, id int64) (*model.Season, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, description, status, start_date, started_at, ended_at,
		 map_template_name, game_speed, resource_multiplier, max_villages_per_player,
		 weapons_of_chaos_count, map_width, map_height, created_at, updated_at
		 FROM seasons WHERE id = ?`, id)
	return scanSeason(row)
}

func (r *seasonRepo) List(ctx context.Context, statusFilter string) ([]*model.Season, error) {
	var rows *sql.Rows
	var err error
	if statusFilter != "" {
		rows, err = r.db.QueryContext(ctx,
			`SELECT id, name, description, status, start_date, started_at, ended_at,
			 map_template_name, game_speed, resource_multiplier, max_villages_per_player,
			 weapons_of_chaos_count, map_width, map_height, created_at, updated_at
			 FROM seasons WHERE status = ? ORDER BY created_at DESC`, statusFilter)
	} else {
		rows, err = r.db.QueryContext(ctx,
			`SELECT id, name, description, status, start_date, started_at, ended_at,
			 map_template_name, game_speed, resource_multiplier, max_villages_per_player,
			 weapons_of_chaos_count, map_width, map_height, created_at, updated_at
			 FROM seasons ORDER BY created_at DESC`)
	}
	if err != nil {
		return nil, fmt.Errorf("list seasons: %w", err)
	}
	defer rows.Close()

	var seasons []*model.Season
	for rows.Next() {
		s, err := scanSeasonRow(rows)
		if err != nil {
			return nil, err
		}
		seasons = append(seasons, s)
	}
	return seasons, rows.Err()
}

func (r *seasonRepo) Update(ctx context.Context, season *model.Season) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE seasons SET name = ?, description = ?, start_date = ?,
		 map_template_name = ?, game_speed = ?, resource_multiplier = ?,
		 max_villages_per_player = ?, weapons_of_chaos_count = ?,
		 map_width = ?, map_height = ?, updated_at = datetime('now')
		 WHERE id = ?`,
		season.Name, season.Description, season.StartDate,
		season.MapTemplateName, season.GameSpeed, season.ResourceMultiplier,
		season.MaxVillagesPerPlayer, season.WeaponsOfChaosCount,
		season.MapWidth, season.MapHeight, season.ID,
	)
	if err != nil {
		return fmt.Errorf("update season %d: %w", season.ID, err)
	}
	return nil
}

func (r *seasonRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	var query string
	switch status {
	case model.SeasonStatusActive:
		query = `UPDATE seasons SET status = ?, started_at = datetime('now'), updated_at = datetime('now') WHERE id = ?`
	case model.SeasonStatusEnded:
		query = `UPDATE seasons SET status = ?, ended_at = datetime('now'), updated_at = datetime('now') WHERE id = ?`
	default:
		query = `UPDATE seasons SET status = ?, updated_at = datetime('now') WHERE id = ?`
	}
	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("update season status %d: %w", id, err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *seasonRepo) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM seasons WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete season %d: %w", id, err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return model.ErrNotFound
	}
	return nil
}

// ── Season Players ───────────────────────────────────────────────────────────

func (r *seasonRepo) AddPlayer(ctx context.Context, sp *model.SeasonPlayer) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO season_players (season_id, player_id, kingdom) VALUES (?, ?, ?)`,
		sp.SeasonID, sp.PlayerID, sp.Kingdom,
	)
	if err != nil {
		return fmt.Errorf("add season player: %w", err)
	}
	id, _ := result.LastInsertId()
	sp.ID = id
	return nil
}

func (r *seasonRepo) RemovePlayer(ctx context.Context, seasonID, playerID int64) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM season_players WHERE season_id = ? AND player_id = ?`,
		seasonID, playerID,
	)
	if err != nil {
		return fmt.Errorf("remove season player: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *seasonRepo) GetSeasonPlayer(ctx context.Context, seasonID, playerID int64) (*model.SeasonPlayer, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, season_id, player_id, kingdom, joined_at
		 FROM season_players WHERE season_id = ? AND player_id = ?`,
		seasonID, playerID,
	)
	var sp model.SeasonPlayer
	var joinedAtStr string
	err := row.Scan(&sp.ID, &sp.SeasonID, &sp.PlayerID, &sp.Kingdom, &joinedAtStr)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan season player: %w", err)
	}
	sp.JoinedAt, _ = parseTime(joinedAtStr)
	return &sp, nil
}

func (r *seasonRepo) ListSeasonPlayers(ctx context.Context, seasonID int64) ([]*model.SeasonPlayer, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, season_id, player_id, kingdom, joined_at
		 FROM season_players WHERE season_id = ? ORDER BY joined_at`, seasonID)
	if err != nil {
		return nil, fmt.Errorf("list season players: %w", err)
	}
	defer rows.Close()

	var players []*model.SeasonPlayer
	for rows.Next() {
		var sp model.SeasonPlayer
		var joinedAtStr string
		if err := rows.Scan(&sp.ID, &sp.SeasonID, &sp.PlayerID, &sp.Kingdom, &joinedAtStr); err != nil {
			return nil, fmt.Errorf("scan season player row: %w", err)
		}
		sp.JoinedAt, _ = parseTime(joinedAtStr)
		players = append(players, &sp)
	}
	return players, rows.Err()
}

func (r *seasonRepo) ListPlayerSeasons(ctx context.Context, playerID int64) ([]*model.SeasonPlayer, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, season_id, player_id, kingdom, joined_at
		 FROM season_players WHERE player_id = ? ORDER BY joined_at DESC`, playerID)
	if err != nil {
		return nil, fmt.Errorf("list player seasons: %w", err)
	}
	defer rows.Close()

	var entries []*model.SeasonPlayer
	for rows.Next() {
		var sp model.SeasonPlayer
		var joinedAtStr string
		if err := rows.Scan(&sp.ID, &sp.SeasonID, &sp.PlayerID, &sp.Kingdom, &joinedAtStr); err != nil {
			return nil, fmt.Errorf("scan player season row: %w", err)
		}
		sp.JoinedAt, _ = parseTime(joinedAtStr)
		entries = append(entries, &sp)
	}
	return entries, rows.Err()
}

func (r *seasonRepo) GetSeasonPlayerCount(ctx context.Context, seasonID int64) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM season_players WHERE season_id = ?`, seasonID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count season players: %w", err)
	}
	return count, nil
}

func (r *seasonRepo) GetPlayerSeasonHistory(ctx context.Context, playerID int64) ([]model.SeasonHistoryRow, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT s.id, s.name, s.status, sp.kingdom, sp.joined_at,
		 (SELECT COUNT(*) FROM villages v WHERE v.player_id = sp.player_id AND v.season_id = s.id) AS village_count
		 FROM season_players sp
		 JOIN seasons s ON s.id = sp.season_id
		 WHERE sp.player_id = ?
		 ORDER BY sp.joined_at DESC`, playerID)
	if err != nil {
		return nil, fmt.Errorf("get player season history: %w", err)
	}
	defer rows.Close()

	var history []model.SeasonHistoryRow
	for rows.Next() {
		var h model.SeasonHistoryRow
		if err := rows.Scan(&h.SeasonID, &h.SeasonName, &h.SeasonStatus, &h.Kingdom, &h.JoinedAt, &h.VillageCount); err != nil {
			return nil, fmt.Errorf("scan season history row: %w", err)
		}
		history = append(history, h)
	}
	return history, rows.Err()
}

// ── Scan helpers ─────────────────────────────────────────────────────────────

func scanSeason(row *sql.Row) (*model.Season, error) {
	var s model.Season
	var startDate sql.NullString
	var startedAt, endedAt nullableTimeStr
	var createdAtStr, updatedAtStr string

	err := row.Scan(
		&s.ID, &s.Name, &s.Description, &s.Status,
		&startDate, &startedAt, &endedAt,
		&s.MapTemplateName, &s.GameSpeed, &s.ResourceMultiplier,
		&s.MaxVillagesPerPlayer, &s.WeaponsOfChaosCount,
		&s.MapWidth, &s.MapHeight,
		&createdAtStr, &updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan season: %w", err)
	}
	if startDate.Valid {
		s.StartDate = &startDate.String
	}
	s.StartedAt = startedAt.Time()
	s.EndedAt = endedAt.Time()
	s.CreatedAt, _ = parseTime(createdAtStr)
	s.UpdatedAt, _ = parseTime(updatedAtStr)
	return &s, nil
}

func scanSeasonRow(rows *sql.Rows) (*model.Season, error) {
	var s model.Season
	var startDate sql.NullString
	var startedAt, endedAt nullableTimeStr
	var createdAtStr, updatedAtStr string

	err := rows.Scan(
		&s.ID, &s.Name, &s.Description, &s.Status,
		&startDate, &startedAt, &endedAt,
		&s.MapTemplateName, &s.GameSpeed, &s.ResourceMultiplier,
		&s.MaxVillagesPerPlayer, &s.WeaponsOfChaosCount,
		&s.MapWidth, &s.MapHeight,
		&createdAtStr, &updatedAtStr,
	)
	if err != nil {
		return nil, fmt.Errorf("scan season row: %w", err)
	}
	if startDate.Valid {
		s.StartDate = &startDate.String
	}
	s.StartedAt = startedAt.Time()
	s.EndedAt = endedAt.Time()
	s.CreatedAt, _ = parseTime(createdAtStr)
	s.UpdatedAt, _ = parseTime(updatedAtStr)
	return &s, nil
}
