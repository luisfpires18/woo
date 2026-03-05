package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type playerRepo struct {
	db *sql.DB
}

// NewPlayerRepo creates a new SQLite-backed PlayerRepository.
func NewPlayerRepo(db *sql.DB) *playerRepo {
	return &playerRepo{db: db}
}

func (r *playerRepo) Create(ctx context.Context, player *model.Player) error {
	if player.Role == "" {
		player.Role = model.RolePlayer
	}
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO players (username, email, password_hash, kingdom, role, oauth_provider, oauth_id, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		player.Username, player.Email, player.PasswordHash,
		player.Kingdom, player.Role, player.OAuthProvider, player.OAuthID, player.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert player: %w", err)
	}
	id, _ := result.LastInsertId()
	player.ID = id
	return nil
}

func (r *playerRepo) GetByID(ctx context.Context, id int64) (*model.Player, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, username, email, password_hash, kingdom, role, oauth_provider, oauth_id, created_at, last_login_at
		 FROM players WHERE id = ?`, id,
	)
	return scanPlayer(row)
}

func (r *playerRepo) GetByEmail(ctx context.Context, email string) (*model.Player, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, username, email, password_hash, kingdom, role, oauth_provider, oauth_id, created_at, last_login_at
		 FROM players WHERE email = ?`, email,
	)
	return scanPlayer(row)
}

func (r *playerRepo) GetByOAuth(ctx context.Context, provider, oauthID string) (*model.Player, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, username, email, password_hash, kingdom, role, oauth_provider, oauth_id, created_at, last_login_at
		 FROM players WHERE oauth_provider = ? AND oauth_id = ?`, provider, oauthID,
	)
	return scanPlayer(row)
}

func (r *playerRepo) UpdateLastLogin(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE players SET last_login_at = datetime('now') WHERE id = ?`, id,
	)
	if err != nil {
		return fmt.Errorf("update last login for player %d: %w", id, err)
	}
	return nil
}

func (r *playerRepo) UpdateRole(ctx context.Context, id int64, role string) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE players SET role = ? WHERE id = ?`, role, id,
	)
	if err != nil {
		return fmt.Errorf("update role for player %d: %w", id, err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *playerRepo) ListAll(ctx context.Context, offset, limit int) ([]*model.Player, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, username, email, password_hash, kingdom, role, oauth_provider, oauth_id, created_at, last_login_at
		 FROM players ORDER BY id LIMIT ? OFFSET ?`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list players: %w", err)
	}
	defer rows.Close()

	var players []*model.Player
	for rows.Next() {
		p, err := scanPlayerRow(rows)
		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}
	return players, rows.Err()
}

func (r *playerRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM players`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count players: %w", err)
	}
	return count, nil
}

// scanPlayer scans a single *sql.Row into a Player.
func scanPlayer(row *sql.Row) (*model.Player, error) {
	var p model.Player
	var createdAtStr string
	var lastLoginStr nullableTimeStr
	var oauthProvider, oauthID sql.NullString
	err := row.Scan(&p.ID, &p.Username, &p.Email, &p.PasswordHash,
		&p.Kingdom, &p.Role, &oauthProvider, &oauthID, &createdAtStr, &lastLoginStr)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan player: %w", err)
	}
	p.OAuthProvider = oauthProvider.String
	p.OAuthID = oauthID.String
	p.CreatedAt, _ = parseTime(createdAtStr)
	p.LastLoginAt = lastLoginStr.Time()
	return &p, nil
}

// scanPlayerRow scans a single *sql.Rows into a Player.
func scanPlayerRow(rows *sql.Rows) (*model.Player, error) {
	var p model.Player
	var createdAtStr string
	var lastLoginStr nullableTimeStr
	var oauthProvider, oauthID sql.NullString
	err := rows.Scan(&p.ID, &p.Username, &p.Email, &p.PasswordHash,
		&p.Kingdom, &p.Role, &oauthProvider, &oauthID, &createdAtStr, &lastLoginStr)
	if err != nil {
		return nil, fmt.Errorf("scan player row: %w", err)
	}
	p.OAuthProvider = oauthProvider.String
	p.OAuthID = oauthID.String
	p.CreatedAt, _ = parseTime(createdAtStr)
	p.LastLoginAt = lastLoginStr.Time()
	return &p, nil
}
