package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/models/sub"
	"github.com/bestruirui/bestsub/internal/utils/log"
)

type SubRepository struct {
	db *DB
}

func (db *DB) Sub() interfaces.SubRepository {
	return &SubRepository{db: db}
}

func (r *SubRepository) Create(ctx context.Context, link *sub.Data) error {
	log.Debugf("Create sub")
	query := `INSERT INTO sub (enable, cron_expr, name, config, result, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.db.ExecContext(ctx, query,
		link.Enable,
		link.CronExpr,
		link.Name,
		link.Config,
		link.Result,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create sub: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get sub id: %w", err)
	}

	link.ID = uint16(id)
	link.CreatedAt = now
	link.UpdatedAt = now

	return nil
}

func (r *SubRepository) GetByID(ctx context.Context, id uint16) (*sub.Data, error) {
	log.Debugf("Get sub by id")
	query := `SELECT id, enable, cron_expr, name, config, result, created_at, updated_at
	          FROM sub WHERE id = ?`

	var s sub.Data
	var enable bool
	err := r.db.db.QueryRowContext(ctx, query, id).Scan(
		&s.ID,
		&enable,
		&s.CronExpr,
		&s.Name,
		&s.Config,
		&s.Result,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get sub by id: %w", err)
	}

	return &s, nil
}

func (r *SubRepository) Update(ctx context.Context, link *sub.Data) error {
	log.Debugf("Update sub")
	query := `UPDATE sub SET enable = ?, cron_expr = ?, name = ?, config = ?, result = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query,
		link.Enable,
		link.CronExpr,
		link.Name,
		link.Config,
		link.Result,
		time.Now(),
		link.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update sub: %w", err)
	}

	return nil
}

func (r *SubRepository) Delete(ctx context.Context, id uint16) error {
	log.Debugf("Delete sub")
	query := `DELETE FROM sub WHERE id = ?`

	_, err := r.db.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete sub: %w", err)
	}

	return nil
}

func (r *SubRepository) List(ctx context.Context) (*[]sub.Data, error) {
	log.Debugf("List sub")
	query := `SELECT id, enable, cron_expr, name, config, result, created_at, updated_at
	          FROM sub ORDER BY created_at DESC`

	rows, err := r.db.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list sub: %w", err)
	}
	defer rows.Close()

	var subs []sub.Data
	for rows.Next() {
		var s sub.Data
		var enable bool
		err := rows.Scan(
			&s.ID,
			&enable,
			&s.CronExpr,
			&s.Name,
			&s.Config,
			&s.Result,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sub: %w", err)
		}
		subs = append(subs, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate subs: %w", err)
	}

	return &subs, nil
}
