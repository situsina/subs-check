package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/sqlite/database"
	"github.com/bestruirui/bestsub/internal/models/notify"
	"github.com/bestruirui/bestsub/internal/utils/local"
)

// NotifyRepository 通知渠道数据访问实现
type NotifyRepository struct {
	db *database.Database
}

// newNotificationChannelRepository 创建通知渠道仓库
func newNotifyRepository(db *database.Database) interfaces.NotifyRepository {
	return &NotifyRepository{db: db}
}

// Create 创建通知渠道
func (r *NotifyRepository) Create(ctx context.Context, channel *notify.Data) error {
	query := `INSERT INTO notify_config (enable, name, description, type, config, test_result, last_test, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := local.Time()
	result, err := r.db.ExecContext(ctx, query,
		channel.Enable,
		channel.Name,
		channel.Description,
		channel.Type,
		channel.Config,
		channel.TestResult,
		channel.LastTest,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create notification channel: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get notification channel id: %w", err)
	}

	channel.ID = id
	channel.CreatedAt = now
	channel.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取通知渠道
func (r *NotifyRepository) GetByID(ctx context.Context, id int64) (*notify.Data, error) {
	query := `SELECT id, enable, name, description, type, config, test_result, last_test, created_at, updated_at
	          FROM notify_config WHERE id = ?`

	var channel notify.Data
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&channel.ID,
		&channel.Enable,
		&channel.Name,
		&channel.Description,
		&channel.Type,
		&channel.Config,
		&channel.TestResult,
		&channel.LastTest,
		&channel.CreatedAt,
		&channel.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get notification channel by id: %w", err)
	}

	return &channel, nil
}

// Update 更新通知渠道
func (r *NotifyRepository) Update(ctx context.Context, channel *notify.Data) error {
	query := `UPDATE notify_config SET name = ?, type = ?, config = ?, enable = ?,
	          test_result = ?, last_test = ?, description = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		channel.Name,
		channel.Type,
		channel.Config,
		channel.Enable,
		channel.TestResult,
		channel.LastTest,
		channel.Description,
		local.Time(),
		channel.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update notification channel: %w", err)
	}

	return nil
}

// Delete 删除通知渠道
func (r *NotifyRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM notify_config WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notification channel: %w", err)
	}

	return nil
}

// List 获取通知渠道列表
func (r *NotifyRepository) List(ctx context.Context, offset, limit int) (*[]notify.Data, error) {
	query := `SELECT id, enable, name, description, type, config, test_result, last_test, created_at, updated_at
	          FROM notify_config ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list notification channels: %w", err)
	}
	defer rows.Close()

	var channels []notify.Data
	for rows.Next() {
		var channel notify.Data
		err := rows.Scan(
			&channel.ID,
			&channel.Enable,
			&channel.Name,
			&channel.Description,
			&channel.Type,
			&channel.Config,
			&channel.TestResult,
			&channel.LastTest,
			&channel.CreatedAt,
			&channel.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification channel: %w", err)
		}
		channels = append(channels, channel)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate notification channels: %w", err)
	}

	return &channels, nil
}

// Count 获取通知渠道总数
func (r *NotifyRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM notify_config`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count notification channels: %w", err)
	}

	return count, nil
}

// GetByTaskID 根据任务ID获取通知渠道列表
func (r *NotifyRepository) GetByTaskID(ctx context.Context, taskID int64) (*[]notify.Data, error) {
	query := `SELECT n.id, n.enable, n.name, n.description, n.type, n.config, n.test_result, n.last_test, n.created_at, n.updated_at
	          FROM notify_config n
	          INNER JOIN notify_task_relations ntr ON n.id = ntr.notify_id
	          WHERE ntr.task_id = ? ORDER BY n.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification channels by task id: %w", err)
	}
	defer rows.Close()

	var channels []notify.Data
	for rows.Next() {
		var channel notify.Data
		err := rows.Scan(
			&channel.ID,
			&channel.Enable,
			&channel.Name,
			&channel.Description,
			&channel.Type,
			&channel.Config,
			&channel.TestResult,
			&channel.LastTest,
			&channel.CreatedAt,
			&channel.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification channel: %w", err)
		}
		channels = append(channels, channel)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate notification channels: %w", err)
	}

	return &channels, nil
}

// AddTaskRelation 添加通知渠道与任务的关联
func (r *NotifyRepository) AddTaskRelation(ctx context.Context, notifyID, taskID int64) error {
	query := `INSERT OR IGNORE INTO notify_task_relations (task_id, notify_id) VALUES (?, ?)`

	_, err := r.db.ExecContext(ctx, query, taskID, notifyID)
	if err != nil {
		return fmt.Errorf("failed to add task relation: %w", err)
	}

	return nil
}

// NotifyTemplateRepository 通知模板数据访问实现
type NotifyTemplateRepository struct {
	db *database.Database
}

// newNotifyTemplateRepository 创建通知模板仓库
func newNotifyTemplateRepository(db *database.Database) interfaces.NotifyTemplateRepository {
	return &NotifyTemplateRepository{db: db}
}

// Create 创建通知模板
func (r *NotifyTemplateRepository) Create(ctx context.Context, template *notify.Template) error {
	query := `INSERT INTO notify_templates (name, description, templates, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?)`

	now := local.Time()
	result, err := r.db.ExecContext(ctx, query,
		template.Name,
		template.Description,
		template.Template,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create notify template: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get notify template id: %w", err)
	}

	template.ID = id
	template.CreatedAt = now
	template.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取通知模板
func (r *NotifyTemplateRepository) GetByID(ctx context.Context, id int64) (*notify.Template, error) {
	query := `SELECT id, name, description, templates, created_at, updated_at
	          FROM notify_templates WHERE id = ?`

	var template notify.Template
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&template.ID,
		&template.Name,
		&template.Description,
		&template.Template,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get notify template by id: %w", err)
	}

	return &template, nil
}

// Update 更新通知模板
func (r *NotifyTemplateRepository) Update(ctx context.Context, template *notify.Template) error {
	query := `UPDATE notify_templates SET name = ?, description = ?, templates = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		template.Name,
		template.Description,
		template.Template,
		local.Time(),
		template.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update notify template: %w", err)
	}

	return nil
}

// Delete 删除通知模板
func (r *NotifyTemplateRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM notify_templates WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notify template: %w", err)
	}

	return nil
}

// List 获取通知模板列表
func (r *NotifyTemplateRepository) List(ctx context.Context, offset, limit int) (*[]notify.Template, error) {
	query := `SELECT id, name, description, templates, created_at, updated_at
	          FROM notify_templates ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list notify templates: %w", err)
	}
	defer rows.Close()

	var templates []notify.Template
	for rows.Next() {
		var template notify.Template
		err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.Description,
			&template.Template,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notify template: %w", err)
		}
		templates = append(templates, template)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate notify templates: %w", err)
	}

	return &templates, nil
}

// Count 获取通知模板总数
func (r *NotifyTemplateRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM notify_templates`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count notify templates: %w", err)
	}

	return count, nil
}

// GetByTaskID 根据任务ID获取通知模板
func (r *NotifyTemplateRepository) GetByTaskID(ctx context.Context, taskID int64) (*notify.Template, error) {
	query := `SELECT nt.id, nt.name, nt.description, nt.templates, nt.created_at, nt.updated_at
	          FROM notify_templates nt
	          INNER JOIN task_notify_template_relations tntr ON nt.id = tntr.notify_template_id
	          WHERE tntr.task_id = ?`

	var template notify.Template
	err := r.db.QueryRowContext(ctx, query, taskID).Scan(
		&template.ID,
		&template.Name,
		&template.Description,
		&template.Template,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get notify template by task id: %w", err)
	}

	return &template, nil
}

// AddTaskRelation 添加通知模板与任务的关联
func (r *NotifyTemplateRepository) AddTaskRelation(ctx context.Context, templateID, taskID int64) error {
	query := `INSERT OR IGNORE INTO task_notify_template_relations (task_id, notify_template_id) VALUES (?, ?)`

	_, err := r.db.ExecContext(ctx, query, taskID, templateID)
	if err != nil {
		return fmt.Errorf("failed to add task relation: %w", err)
	}

	return nil
}
