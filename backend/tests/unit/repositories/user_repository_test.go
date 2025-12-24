package repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
	"go.uber.org/zap"
)

func setupTestDB(t *testing.T) (*database.DB, func()) {
	logger, _ := zap.NewDevelopment()
	dbConfig := database.DatabaseConfig{
		Driver:     "sqlite",
		SQLitePath: "file::memory:?cache=shared",
	}
	db, err := database.NewConnection(dbConfig, logger)
	require.NoError(t, err)

	err = db.RunMigrations("../../migrations")
	require.NoError(t, err)

	return db, func() {
		db.Close()
	}
}

func TestUserRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)
	ctx := context.Background()

	user := &models.User{
		ID:                uuid.New(),
		Email:             "test@example.com",
		PasswordHash:      "hashed_password",
		Name:              "Test User",
		Status:            "active",
		Role:              "user",
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	// Verify user was created
	retrieved, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, retrieved.Email)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)
	ctx := context.Background()

	user := &models.User{
		ID:                uuid.New(),
		Email:             "test@example.com",
		PasswordHash:      "hashed_password",
		Name:              "Test User",
		Status:            "active",
		Role:              "user",
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	repo.Create(ctx, user)

	retrieved, err := repo.GetByEmail(ctx, "test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.Email, retrieved.Email)
}

