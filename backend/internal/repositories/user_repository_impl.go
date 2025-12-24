package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
)

type userRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (
			id, email, password_hash, name, first_name, last_name, phone,
			role, signup_source, status, password_changed_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Name,
		user.FirstName, user.LastName, user.Phone, user.Role,
		user.SignupSource, user.Status,
		user.PasswordChangedAt, user.CreatedAt, user.UpdatedAt,
	)

	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, email_verified, password_hash, name, first_name, last_name,
			photo_url, phone, role, phone_verified, status, signup_source,
			password_changed_at, created_at, updated_at, last_login_at
		FROM users
		WHERE id = ?
	`

	user := &models.User{}
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.EmailVerified, &user.PasswordHash,
		&user.Name, &user.FirstName, &user.LastName, &user.PhotoURL,
		&user.Phone, &user.Role, &user.PhoneVerified, &user.Status, &user.SignupSource,
		&user.PasswordChangedAt, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, email_verified, password_hash, name, first_name, last_name,
			photo_url, phone, role, phone_verified, status, signup_source,
			password_changed_at, created_at, updated_at, last_login_at
		FROM users
		WHERE email = ?
	`

	user := &models.User{}
	err := r.db.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.EmailVerified, &user.PasswordHash,
		&user.Name, &user.FirstName, &user.LastName, &user.PhotoURL,
		&user.Phone, &user.Role, &user.PhoneVerified, &user.Status, &user.SignupSource,
		&user.PasswordChangedAt, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET email = ?, name = ?, first_name = ?, last_name = ?,
			photo_url = ?, phone = ?, status = ?, last_login_at = ?,
			updated_at = ?, role = ?
		WHERE id = ?
	`

	_, err := r.db.DB.ExecContext(ctx, query,
		user.Email, user.Name, user.FirstName, user.LastName,
		user.PhotoURL, user.Phone, user.Status, user.LastLoginAt, user.UpdatedAt, user.Role, user.ID,
	)

	return err
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string, changedAt time.Time) error {
	query := `
		UPDATE users
		SET password_hash = ?, password_changed_at = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.DB.ExecContext(ctx, query, passwordHash, changedAt, time.Now(), userID)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET status = 'deleted', updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.DB.ExecContext(ctx, query, id)
	return err
}

func (r *userRepository) MarkDeleted(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET status = 'deleted',
		    status_changed_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := r.db.DB.ExecContext(ctx, query, id)
	return err
}

func (r *userRepository) PurgeDeletedBefore(ctx context.Context, cutoff time.Time) error {
	query := `
		DELETE FROM users
		WHERE status = 'deleted'
		  AND status_changed_at IS NOT NULL
		  AND status_changed_at <= ?
	`
	_, err := r.db.DB.ExecContext(ctx, query, cutoff)
	return err
}

func (r *userRepository) List(ctx context.Context, search string) ([]*models.User, error) {
	var rows *sql.Rows
	var err error

	if search != "" {
		pattern := "%" + search + "%"
		query := `
			SELECT id, email, email_verified, password_hash, name, first_name, last_name,
				photo_url, phone, role, phone_verified, status, signup_source,
				password_changed_at, created_at, updated_at, last_login_at
			FROM users
			WHERE email LIKE ? OR name LIKE ?
			ORDER BY created_at DESC
			LIMIT 200
		`
		rows, err = r.db.DB.QueryContext(ctx, query, pattern, pattern)
	} else {
		query := `
			SELECT id, email, email_verified, password_hash, name, first_name, last_name,
				photo_url, phone, role, phone_verified, status, signup_source,
				password_changed_at, created_at, updated_at, last_login_at
			FROM users
			ORDER BY created_at DESC
			LIMIT 200
		`
		rows, err = r.db.DB.QueryContext(ctx, query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(
			&user.ID, &user.Email, &user.EmailVerified, &user.PasswordHash,
			&user.Name, &user.FirstName, &user.LastName, &user.PhotoURL,
			&user.Phone, &user.Role, &user.PhoneVerified, &user.Status, &user.SignupSource,
			&user.PasswordChangedAt, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) SetStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `
		UPDATE users
		SET status = ?, status_changed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := r.db.DB.ExecContext(ctx, query, status, id)
	return err
}
