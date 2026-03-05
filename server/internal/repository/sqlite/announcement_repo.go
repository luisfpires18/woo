package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/luisfpires18/woo/internal/model"
)

type announcementRepo struct {
	db *sql.DB
}

// NewAnnouncementRepo creates a new SQLite-backed AnnouncementRepository.
func NewAnnouncementRepo(db *sql.DB) *announcementRepo {
	return &announcementRepo{db: db}
}

func (r *announcementRepo) Create(ctx context.Context, a *model.Announcement) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO announcements (title, content, author_id, created_at, expires_at)
		 VALUES (?, ?, ?, datetime('now'), ?)`,
		a.Title, a.Content, a.AuthorID, nullableTime(a.ExpiresAt),
	)
	if err != nil {
		return fmt.Errorf("insert announcement: %w", err)
	}
	id, _ := result.LastInsertId()
	a.ID = id
	return nil
}

func (r *announcementRepo) ListActive(ctx context.Context) ([]*model.Announcement, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, content, author_id, created_at, expires_at
		 FROM announcements
		 WHERE expires_at IS NULL OR expires_at > datetime('now')
		 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list active announcements: %w", err)
	}
	defer rows.Close()

	var announcements []*model.Announcement
	for rows.Next() {
		var a model.Announcement
		var createdAtStr string
		var expiresAtStr nullableTimeStr
		if err := rows.Scan(&a.ID, &a.Title, &a.Content, &a.AuthorID, &createdAtStr, &expiresAtStr); err != nil {
			return nil, fmt.Errorf("scan announcement: %w", err)
		}
		a.CreatedAt, _ = parseTime(createdAtStr)
		a.ExpiresAt = expiresAtStr.Time()
		announcements = append(announcements, &a)
	}
	return announcements, rows.Err()
}

func (r *announcementRepo) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx,
		`DELETE FROM announcements WHERE id = ?`, id,
	)
	if err != nil {
		return fmt.Errorf("delete announcement %d: %w", id, err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return model.ErrNotFound
	}
	return nil
}

// nullableTime converts a *time.Time to a value suitable for SQLite (string or nil).
func nullableTime(t *time.Time) any {
	if t == nil {
		return nil
	}
	return t.UTC().Format("2006-01-02 15:04:05")
}
