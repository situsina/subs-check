package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bestruirui/bestsub/internal/database/interfaces"
	"github.com/bestruirui/bestsub/internal/database/sqlite/database"
	"github.com/bestruirui/bestsub/internal/models/auth"
	"github.com/bestruirui/bestsub/internal/utils/local"
	"github.com/bestruirui/bestsub/internal/utils/passwd"
)

// AuthRepository 认证数据访问实现
type AuthRepository struct {
	db *database.Database
}

// NewAuthRepository 创建认证仓库
func newAuthRepository(db *database.Database) interfaces.AuthRepository {
	return &AuthRepository{db: db}
}

// Get 获取认证信息
func (r *AuthRepository) Get(ctx context.Context) (*auth.Data, error) {
	query := `SELECT id, user_name, password, created_at, updated_at FROM auth LIMIT 1`

	var authData auth.Data
	err := r.db.QueryRowContext(ctx, query).Scan(
		&authData.ID,
		&authData.UserName,
		&authData.Password,
		&authData.CreatedAt,
		&authData.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get auth: %w", err)
	}

	return &authData, nil
}

// Update 更新认证信息
func (r *AuthRepository) Update(ctx context.Context, authData *auth.Data) error {

	hashedPassword, err := passwd.Hash(authData.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `UPDATE auth SET password = ?, updated_at = ?`

	result, err := r.db.ExecContext(ctx, query, hashedPassword, local.Time())
	if err != nil {
		return fmt.Errorf("failed to update auth: %w", err)
	}

	// 检查是否有记录被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("用户不存在")
	}

	return nil
}

// UpdateUsername 更新用户名
func (r *AuthRepository) UpdateUsername(ctx context.Context, username string) error {
	if username == "" {
		return fmt.Errorf("用户名不能为空")
	}

	// 单用户系统，直接更新唯一记录
	query := `UPDATE auth SET user_name = ?, updated_at = ?`

	result, err := r.db.ExecContext(ctx, query, username, local.Time())
	if err != nil {
		return fmt.Errorf("failed to update username: %w", err)
	}

	// 检查是否有记录被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("用户不存在")
	}

	return nil
}

// Initialize 初始化认证信息
func (r *AuthRepository) Initialize(ctx context.Context, authData *auth.Data) error {

	hashedPassword, err := passwd.Hash(authData.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `INSERT INTO auth (user_name, password, created_at, updated_at) VALUES (?, ?, ?, ?)`

	now := local.Time()
	_, err = r.db.ExecContext(ctx, query, authData.UserName, hashedPassword, now, now)
	if err != nil {
		return fmt.Errorf("failed to initialize auth: %w", err)
	}

	return nil
}

// IsInitialized 验证是否已初始化
func (r *AuthRepository) IsInitialized(ctx context.Context) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM auth LIMIT 1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check auth initialization: %w", err)
	}

	return exists, nil
}

// VerifyPassword 验证密码
func (r *AuthRepository) VerifyPassword(ctx context.Context, username, password string) error {
	if username == "" {
		return fmt.Errorf("用户名不能为空")
	}

	if password == "" {
		return fmt.Errorf("密码不能为空")
	}

	// 获取存储的密码哈希
	query := `SELECT password FROM auth WHERE user_name = ? LIMIT 1`

	var hashedPassword string
	err := r.db.QueryRowContext(ctx, query, username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("用户名或密码错误")
		}
		return fmt.Errorf("failed to get user password: %w", err)
	}

	// 验证密码
	return passwd.Verify(password, hashedPassword)
}

// SessionRepository 会话数据访问实现
type SessionRepository struct {
	db *database.Database
}

// NewSessionRepository 创建会话仓库
func newSessionRepository(db *database.Database) interfaces.SessionRepository {
	return &SessionRepository{db: db}
}

// Create 创建会话
func (r *SessionRepository) Create(ctx context.Context, session *auth.Session) error {
	query := `INSERT INTO sessions (token_hash, expires_at, refresh_token, ip_address, user_agent, is_active, created_at, updated_at)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	now := local.Time()
	result, err := r.db.ExecContext(ctx, query,
		session.TokenHash,
		session.ExpiresAt,
		session.RefreshToken,
		session.IPAddress,
		session.UserAgent,
		session.IsActive,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get session id: %w", err)
	}

	session.ID = id
	session.CreatedAt = now
	session.UpdatedAt = now

	return nil
}

// GetByID 根据ID获取会话
func (r *SessionRepository) GetByID(ctx context.Context, id int64) (*auth.Session, error) {
	query := `SELECT id, token_hash, expires_at, refresh_token, ip_address, user_agent, is_active, created_at, updated_at
	          FROM sessions WHERE id = ?`

	var session auth.Session
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.TokenHash,
		&session.ExpiresAt,
		&session.RefreshToken,
		&session.IPAddress,
		&session.UserAgent,
		&session.IsActive,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get session by id: %w", err)
	}

	return &session, nil
}

// GetByTokenHash 根据Token哈希获取会话
func (r *SessionRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*auth.Session, error) {
	query := `SELECT id, token_hash, expires_at, refresh_token, ip_address, user_agent, is_active, created_at, updated_at
	          FROM sessions WHERE token_hash = ? AND is_active = true`

	var session auth.Session
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&session.ID,
		&session.TokenHash,
		&session.ExpiresAt,
		&session.RefreshToken,
		&session.IPAddress,
		&session.UserAgent,
		&session.IsActive,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get session by token hash: %w", err)
	}

	return &session, nil
}

// GetByRefreshToken 根据刷新Token获取会话
func (r *SessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*auth.Session, error) {
	query := `SELECT id, token_hash, expires_at, refresh_token, ip_address, user_agent, is_active, created_at, updated_at
	          FROM sessions WHERE refresh_token = ? AND is_active = true`

	var session auth.Session
	err := r.db.QueryRowContext(ctx, query, refreshToken).Scan(
		&session.ID,
		&session.TokenHash,
		&session.ExpiresAt,
		&session.RefreshToken,
		&session.IPAddress,
		&session.UserAgent,
		&session.IsActive,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get session by refresh token: %w", err)
	}

	return &session, nil
}

// Update 更新会话
func (r *SessionRepository) Update(ctx context.Context, session *auth.Session) error {
	query := `UPDATE sessions SET token_hash = ?, expires_at = ?, refresh_token = ?, ip_address = ?,
	          user_agent = ?, is_active = ?, updated_at = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		session.TokenHash,
		session.ExpiresAt,
		session.RefreshToken,
		session.IPAddress,
		session.UserAgent,
		session.IsActive,
		local.Time(),
		session.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

// Delete 删除会话
func (r *SessionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM sessions WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

// DeleteAll 删除所有会话
func (r *SessionRepository) DeleteAll(ctx context.Context) error {
	query := `DELETE FROM sessions`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete all sessions: %w", err)
	}

	return nil
}

// DeleteExpired 删除过期会话
func (r *SessionRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < ?`

	_, err := r.db.ExecContext(ctx, query, local.Time())
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}

// GetAllActive 获取所有活跃会话
func (r *SessionRepository) GetAllActive(ctx context.Context) ([]*auth.Session, error) {
	query := `SELECT id, token_hash, expires_at, refresh_token, ip_address, user_agent, is_active, created_at, updated_at
	          FROM sessions WHERE is_active = true ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}
	defer rows.Close()

	var sessions []*auth.Session
	for rows.Next() {
		var session auth.Session
		err := rows.Scan(
			&session.ID,
			&session.TokenHash,
			&session.ExpiresAt,
			&session.RefreshToken,
			&session.IPAddress,
			&session.UserAgent,
			&session.IsActive,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, &session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate sessions: %w", err)
	}

	return sessions, nil
}

// DeactivateAll 停用所有会话
func (r *SessionRepository) DeactivateAll(ctx context.Context) error {
	query := `UPDATE sessions SET is_active = false, updated_at = ?`

	_, err := r.db.ExecContext(ctx, query, local.Time())
	if err != nil {
		return fmt.Errorf("failed to deactivate all sessions: %w", err)
	}

	return nil
}
