package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type refreshTokenRepo struct {
	db *sql.DB
}

// NewRefreshTokenRepo creates a new SQLite-backed RefreshTokenRepository.
func NewRefreshTokenRepo(db *sql.DB) *refreshTokenRepo {
	return &refreshTokenRepo{db: db}
}

func (r *refreshTokenRepo) Create(ctx context.Context, token *model.RefreshToken) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO refresh_tokens (player_id, token_hash, expires_at, created_at)
		 VALUES (?, ?, ?, ?)`,
		token.PlayerID, token.TokenHash, token.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"), token.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	)
	if err != nil {
		return fmt.Errorf("insert refresh token: %w", err)
	}
	id, _ := result.LastInsertId()
	token.ID = id
	return nil
}

func (r *refreshTokenRepo) GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, player_id, token_hash, expires_at, created_at
		 FROM refresh_tokens WHERE token_hash = ?`, tokenHash,
	)
	var t model.RefreshToken
	var expiresAtStr, createdAtStr string
	err := row.Scan(&t.ID, &t.PlayerID, &t.TokenHash, &expiresAtStr, &createdAtStr)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan refresh token: %w", err)
	}
	t.ExpiresAt, _ = parseTime(expiresAtStr)
	t.CreatedAt, _ = parseTime(createdAtStr)
	return &t, nil
}

func (r *refreshTokenRepo) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM refresh_tokens WHERE token_hash = ?`, tokenHash,
	)
	if err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}
	return nil
}

func (r *refreshTokenRepo) DeleteAllByPlayerID(ctx context.Context, playerID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM refresh_tokens WHERE player_id = ?`, playerID,
	)
	if err != nil {
		return fmt.Errorf("delete refresh tokens for player %d: %w", playerID, err)
	}
	return nil
}
