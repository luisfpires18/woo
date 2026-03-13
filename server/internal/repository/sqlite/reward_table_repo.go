package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type rewardTableRepo struct {
	db *sql.DB
}

// NewRewardTableRepo creates a new SQLite-backed RewardTableRepository.
func NewRewardTableRepo(db *sql.DB) *rewardTableRepo {
	return &rewardTableRepo{db: db}
}

func (r *rewardTableRepo) Create(ctx context.Context, rt *model.RewardTable) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO reward_tables (name, updated_by) VALUES (?, ?)`,
		rt.Name, rt.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("create reward table: %w", err)
	}
	id, _ := result.LastInsertId()
	rt.ID = id
	return nil
}

func (r *rewardTableRepo) GetByID(ctx context.Context, id int64) (*model.RewardTable, error) {
	var rt model.RewardTable
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, created_at, updated_at, updated_by FROM reward_tables WHERE id = ?`, id,
	).Scan(&rt.ID, &rt.Name, &rt.CreatedAt, &rt.UpdatedAt, &rt.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get reward table %d: %w", id, err)
	}
	return &rt, nil
}

func (r *rewardTableRepo) GetAll(ctx context.Context) ([]*model.RewardTable, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, created_at, updated_at, updated_by FROM reward_tables ORDER BY id ASC`)
	if err != nil {
		return nil, fmt.Errorf("list reward tables: %w", err)
	}
	defer rows.Close()

	var tables []*model.RewardTable
	for rows.Next() {
		var rt model.RewardTable
		if err := rows.Scan(&rt.ID, &rt.Name, &rt.CreatedAt, &rt.UpdatedAt, &rt.UpdatedBy); err != nil {
			return nil, fmt.Errorf("scan reward table: %w", err)
		}
		tables = append(tables, &rt)
	}
	return tables, rows.Err()
}

func (r *rewardTableRepo) Update(ctx context.Context, rt *model.RewardTable) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE reward_tables SET name = ?, updated_at = datetime('now'), updated_by = ? WHERE id = ?`,
		rt.Name, rt.UpdatedBy, rt.ID,
	)
	if err != nil {
		return fmt.Errorf("update reward table %d: %w", rt.ID, err)
	}
	return nil
}

func (r *rewardTableRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM reward_tables WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete reward table %d: %w", id, err)
	}
	return nil
}

// ── Reward Table Entries ──────────────────────────────────────────────────────

type rewardTableEntryRepo struct {
	db *sql.DB
}

// NewRewardTableEntryRepo creates a new SQLite-backed RewardTableEntryRepository.
func NewRewardTableEntryRepo(db *sql.DB) *rewardTableEntryRepo {
	return &rewardTableEntryRepo{db: db}
}

func (r *rewardTableEntryRepo) Create(ctx context.Context, entry *model.RewardTableEntry) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO reward_table_entries (reward_table_id, reward_type, min_amount, max_amount, drop_chance_pct)
		 VALUES (?, ?, ?, ?, ?)`,
		entry.RewardTableID, entry.RewardType, entry.MinAmount, entry.MaxAmount, entry.DropChancePct,
	)
	if err != nil {
		return fmt.Errorf("create reward table entry: %w", err)
	}
	id, _ := result.LastInsertId()
	entry.ID = id
	return nil
}

func (r *rewardTableEntryRepo) GetByRewardTableID(ctx context.Context, rewardTableID int64) ([]*model.RewardTableEntry, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, reward_table_id, reward_type, min_amount, max_amount, drop_chance_pct, created_at
		 FROM reward_table_entries WHERE reward_table_id = ? ORDER BY id ASC`, rewardTableID,
	)
	if err != nil {
		return nil, fmt.Errorf("list reward table entries for table %d: %w", rewardTableID, err)
	}
	defer rows.Close()

	var entries []*model.RewardTableEntry
	for rows.Next() {
		var e model.RewardTableEntry
		if err := rows.Scan(&e.ID, &e.RewardTableID, &e.RewardType, &e.MinAmount, &e.MaxAmount, &e.DropChancePct, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan reward table entry: %w", err)
		}
		entries = append(entries, &e)
	}
	return entries, rows.Err()
}

func (r *rewardTableEntryRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM reward_table_entries WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete reward table entry %d: %w", id, err)
	}
	return nil
}

func (r *rewardTableEntryRepo) DeleteByRewardTableID(ctx context.Context, rewardTableID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM reward_table_entries WHERE reward_table_id = ?`, rewardTableID)
	if err != nil {
		return fmt.Errorf("delete reward table entries for table %d: %w", rewardTableID, err)
	}
	return nil
}
