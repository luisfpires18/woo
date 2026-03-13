package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type adminAuditLogRepo struct {
	db *sql.DB
}

// NewAdminAuditLogRepo creates a new SQLite-backed AdminAuditLogRepository.
func NewAdminAuditLogRepo(db *sql.DB) *adminAuditLogRepo {
	return &adminAuditLogRepo{db: db}
}

func (r *adminAuditLogRepo) Create(ctx context.Context, entry *model.AdminAuditLog) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO admin_audit_log (admin_player_id, action, entity_type, entity_id, old_value_json, new_value_json)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		entry.AdminPlayerID, entry.Action, entry.EntityType, entry.EntityID,
		entry.OldValueJSON, entry.NewValueJSON,
	)
	if err != nil {
		return fmt.Errorf("create audit log entry: %w", err)
	}
	id, _ := result.LastInsertId()
	entry.ID = id
	return nil
}

func (r *adminAuditLogRepo) List(ctx context.Context, entityType string, limit, offset int) ([]*model.AdminAuditLog, error) {
	query := `SELECT id, admin_player_id, action, entity_type, entity_id, old_value_json, new_value_json, created_at FROM admin_audit_log`
	args := []interface{}{}

	if entityType != "" {
		query += ` WHERE entity_type = ?`
		args = append(args, entityType)
	}
	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list audit log: %w", err)
	}
	defer rows.Close()

	var entries []*model.AdminAuditLog
	for rows.Next() {
		var e model.AdminAuditLog
		if err := rows.Scan(&e.ID, &e.AdminPlayerID, &e.Action, &e.EntityType, &e.EntityID,
			&e.OldValueJSON, &e.NewValueJSON, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan audit log entry: %w", err)
		}
		entries = append(entries, &e)
	}
	return entries, rows.Err()
}
