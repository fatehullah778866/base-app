package repositories

import (
	"context"
	"database/sql"
	"fmt"

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
			signup_source, status, password_changed_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Name,
		user.FirstName, user.LastName, user.Phone,
		user.SignupSource, user.Status,
		user.PasswordChangedAt, user.CreatedAt, user.UpdatedAt,
	)

	return err
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, email_verified, password_hash, name, first_name, last_name,
			photo_url, phone, phone_verified, status, signup_source,
			password_changed_at, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.EmailVerified, &user.PasswordHash,
		&user.Name, &user.FirstName, &user.LastName, &user.PhotoURL,
		&user.Phone, &user.PhoneVerified, &user.Status, &user.SignupSource,
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
			photo_url, phone, phone_verified, status, signup_source,
			password_changed_at, created_at, updated_at, last_login_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.EmailVerified, &user.PasswordHash,
		&user.Name, &user.FirstName, &user.LastName, &user.PhotoURL,
		&user.Phone, &user.PhoneVerified, &user.Status, &user.SignupSource,
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
		SET email = $2, name = $3, first_name = $4, last_name = $5,
			photo_url = $6, phone = $7, status = $8, last_login_at = $9,
			updated_at = $10
		WHERE id = $1
	`

	_, err := r.db.DB.ExecContext(ctx, query,
		user.ID, user.Email, user.Name, user.FirstName, user.LastName,
		user.PhotoURL, user.Phone, user.Status, user.LastLoginAt, user.UpdatedAt,
	)

	return err
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET status = 'deleted', updated_at = NOW() WHERE id = $1`
	_, err := r.db.DB.ExecContext(ctx, query, id)
	return err
}

