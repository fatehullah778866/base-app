package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"base-app-service/internal/config"
	"base-app-service/internal/database"
	"base-app-service/internal/handlers"
	"base-app-service/internal/middleware"
	"base-app-service/internal/repositories"
	"base-app-service/internal/services"
	"go.uber.org/zap"
)

func setupTestServer(t *testing.T) (*httptest.Server, func()) {
	logger, _ := zap.NewDevelopment()
	
	// Use in-memory SQLite for testing
	dbConfig := database.DatabaseConfig{
		Driver:     "sqlite",
		SQLitePath: "file::memory:?cache=shared",
	}
	db, err := database.NewConnection(dbConfig, logger)
	require.NoError(t, err)

	// Run migrations
	err = db.RunMigrations("../../migrations")
	require.NoError(t, err)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	deviceRepo := repositories.NewDeviceRepository(db)

	// Initialize services
	cfg, _ := config.Load()
	authService := services.NewAuthService(
		userRepo, sessionRepo, deviceRepo,
		cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry,
		logger,
	)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, nil, nil, logger)

	// Setup router
	router := http.NewServeMux()
	router.HandleFunc("/v1/auth/signup", authHandler.Signup)
	router.HandleFunc("/v1/auth/login", authHandler.Login)

	server := httptest.NewServer(router)

	return server, func() {
		server.Close()
		db.Close()
	}
}

func TestSignupIntegration(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	t.Run("successful signup", func(t *testing.T) {
		payload := map[string]interface{}{
			"email":         "test@example.com",
			"password":      "Test123!@#",
			"name":          "Test User",
			"terms_accepted": true,
			"terms_version":  "1.0",
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", server.URL+"/v1/auth/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.True(t, result["success"].(bool))
	})

	t.Run("invalid email", func(t *testing.T) {
		payload := map[string]interface{}{
			"email":         "invalid-email",
			"password":      "Test123!@#",
			"name":          "Test User",
			"terms_accepted": true,
			"terms_version":  "1.0",
		}

		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", server.URL+"/v1/auth/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})
}

