package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
)

type sessionRepository struct {
	db *database.DB
}

func NewSessionRepository(db *database.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (
			id, user_id, token, refresh_token, refresh_token_expires_at,
			device_id, device_name, ip_address, is_active, expires_at,
			created_at, last_used_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.DB.ExecContext(ctx, query,
		session.ID, session.UserID, session.Token, session.RefreshToken,
		session.RefreshTokenExpiresAt, session.DeviceID, session.DeviceName,
		session.IPAddress, session.IsActive, session.ExpiresAt,
		session.CreatedAt, session.LastUsedAt,
	)

	return err
}

func (r *sessionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	query := `
		SELECT id, user_id, token, refresh_token, refresh_token_expires_at,
			device_id, device_type, device_name, os, browser, ip_address,
			location_country, location_city, is_active, expires_at,
			created_at, last_used_at
		FROM sessions
		WHERE id = ? AND is_active = 1
	`

	session := &models.Session{}
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&session.ID, &session.UserID, &session.Token, &session.RefreshToken,
		&session.RefreshTokenExpiresAt, &session.DeviceID, &session.DeviceType,
		&session.DeviceName, &session.OS, &session.Browser, &session.IPAddress,
		&session.LocationCountry, &session.LocationCity, &session.IsActive,
		&session.ExpiresAt, &session.CreatedAt, &session.LastUsedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *sessionRepository) GetByToken(ctx context.Context, token string) (*models.Session, error) {
	query := `
		SELECT id, user_id, token, refresh_token, refresh_token_expires_at,
			device_id, device_type, device_name, os, browser, ip_address,
			location_country, location_city, is_active, expires_at,
			created_at, last_used_at
		FROM sessions
		WHERE token = ? AND is_active = 1
	`

	session := &models.Session{}
	err := r.db.DB.QueryRowContext(ctx, query, token).Scan(
		&session.ID, &session.UserID, &session.Token, &session.RefreshToken,
		&session.RefreshTokenExpiresAt, &session.DeviceID, &session.DeviceType,
		&session.DeviceName, &session.OS, &session.Browser, &session.IPAddress,
		&session.LocationCountry, &session.LocationCity, &session.IsActive,
		&session.ExpiresAt, &session.CreatedAt, &session.LastUsedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *sessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	query := `
		SELECT id, user_id, token, refresh_token, refresh_token_expires_at,
			device_id, device_type, device_name, os, browser, ip_address,
			location_country, location_city, is_active, expires_at,
			created_at, last_used_at
		FROM sessions
		WHERE refresh_token = ? AND is_active = 1
	`

	session := &models.Session{}
	err := r.db.DB.QueryRowContext(ctx, query, refreshToken).Scan(
		&session.ID, &session.UserID, &session.Token, &session.RefreshToken,
		&session.RefreshTokenExpiresAt, &session.DeviceID, &session.DeviceType,
		&session.DeviceName, &session.OS, &session.Browser, &session.IPAddress,
		&session.LocationCountry, &session.LocationCity, &session.IsActive,
		&session.ExpiresAt, &session.CreatedAt, &session.LastUsedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *sessionRepository) Update(ctx context.Context, session *models.Session) error {
	query := `
		UPDATE sessions
		SET token = ?, refresh_token = ?, refresh_token_expires_at = ?,
			expires_at = ?, last_used_at = ?
		WHERE id = ?
	`

	_, err := r.db.DB.ExecContext(ctx, query,
		session.Token, session.RefreshToken,
		session.RefreshTokenExpiresAt, session.ExpiresAt, session.LastUsedAt, session.ID,
	)

	return err
}

func (r *sessionRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE sessions SET is_active = 0, revoked_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.DB.ExecContext(ctx, query, id)
	return err
}

func (r *sessionRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE sessions SET is_active = 0, revoked_at = CURRENT_TIMESTAMP WHERE user_id = ? AND is_active = 1`
	_, err := r.db.DB.ExecContext(ctx, query, userID)
	return err
}
