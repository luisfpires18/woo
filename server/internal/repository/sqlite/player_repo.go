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
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO players (username, email, password_hash, kingdom, oauth_provider, oauth_id, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		player.Username, player.Email, player.PasswordHash,
		player.Kingdom, player.OAuthProvider, player.OAuthID, player.CreatedAt,
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
		`SELECT id, username, email, password_hash, kingdom, oauth_provider, oauth_id, created_at, last_login_at
		 FROM players WHERE id = ?`, id,
	)
	return scanPlayer(row)
}

func (r *playerRepo) GetByEmail(ctx context.Context, email string) (*model.Player, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, username, email, password_hash, kingdom, oauth_provider, oauth_id, created_at, last_login_at
		 FROM players WHERE email = ?`, email,
	)
	return scanPlayer(row)
}

func (r *playerRepo) GetByOAuth(ctx context.Context, provider, oauthID string) (*model.Player, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, username, email, password_hash, kingdom, oauth_provider, oauth_id, created_at, last_login_at
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

func scanPlayer(row *sql.Row) (*model.Player, error) {
	var p model.Player
	var createdAtStr string
	var lastLoginStr nullableTimeStr
	err := row.Scan(&p.ID, &p.Username, &p.Email, &p.PasswordHash,
		&p.Kingdom, &p.OAuthProvider, &p.OAuthID, &createdAtStr, &lastLoginStr)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan player: %w", err)
	}
	p.CreatedAt, _ = parseTime(createdAtStr)
	p.LastLoginAt = lastLoginStr.Time()
	return &p, nil
}
