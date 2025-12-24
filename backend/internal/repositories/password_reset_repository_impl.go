package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
)

type passwordResetRepository struct {
	db *database.DB
}

func NewPasswordResetRepository(db *database.DB) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) Create(ctx context.Context, token *models.PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (id, user_id, token, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		token.ID.String(),
		token.UserID.String(),
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
	)
	return err
}

func (r *passwordResetRepository) GetByToken(ctx context.Context, token string) (*models.PasswordResetToken, error) {
	var t models.PasswordResetToken
	var userIDStr, idStr string
	var usedAt sql.NullTime

	query := `
		SELECT id, user_id, token, expires_at, used_at, created_at
		FROM password_reset_tokens
		WHERE token = ? AND expires_at > datetime('now')
	`
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&idStr,
		&userIDStr,
		&t.Token,
		&t.ExpiresAt,
		&usedAt,
		&t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	id, _ := uuid.Parse(idStr)
	userID, _ := uuid.Parse(userIDStr)
	t.ID = id
	t.UserID = userID
	if usedAt.Valid {
		t.UsedAt = &usedAt.Time
	}

	return &t, nil
}

func (r *passwordResetRepository) MarkAsUsed(ctx context.Context, tokenID uuid.UUID) error {
	query := `UPDATE password_reset_tokens SET used_at = datetime('now') WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, tokenID.String())
	return err
}

func (r *passwordResetRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	query := `DELETE FROM password_reset_tokens WHERE expires_at < ?`
	_, err := r.db.ExecContext(ctx, query, before)
	return err
}

func (r *passwordResetRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM password_reset_tokens WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID.String())
	return err
}

