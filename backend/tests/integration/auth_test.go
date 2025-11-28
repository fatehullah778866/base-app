package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"base-app-service/internal/handlers"
	"base-app-service/internal/models"
	"base-app-service/internal/services"
)

func TestHealthCheck(t *testing.T) {
	router := setupTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "healthy", resp["status"])
}

func TestSignup(t *testing.T) {
	router := setupTestRouter(t)

	payload := map[string]interface{}{
		"email":             "test@example.com",
		"password":          "TestPassword123!",
		"name":              "Test User",
		"terms_accepted":    true,
		"terms_version":     "1.0",
		"marketing_consent": true,
	}

	body, err := json.Marshal(payload)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/v1/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Product-Name", "test-product")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			User struct {
				ID            string `json:"id"`
				Email         string `json:"email"`
				Name          string `json:"name"`
				Status        string `json:"status"`
				EmailVerified bool   `json:"email_verified"`
			} `json:"user"`
			Session struct {
				Token        string `json:"token"`
				RefreshToken string `json:"refresh_token"`
				ExpiresAt    string `json:"expires_at"`
			} `json:"session"`
		} `json:"data"`
	}

	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.True(t, response.Success)
	assert.Equal(t, "test@example.com", response.Data.User.Email)
	assert.Equal(t, "Test User", response.Data.User.Name)
	assert.Equal(t, "pending", response.Data.User.Status)
	assert.False(t, response.Data.User.EmailVerified)
	assert.NotEmpty(t, response.Data.User.ID)
	assert.NotEmpty(t, response.Data.Session.Token)
	assert.NotEmpty(t, response.Data.Session.RefreshToken)
	assert.NotEmpty(t, response.Data.Session.ExpiresAt)
}

func setupTestRouter(t *testing.T) *mux.Router {
	t.Helper()

	logger := zap.NewNop()

	userRepo := newMockUserRepo()
	sessionRepo := newMockSessionRepo()
	deviceRepo := newMockDeviceRepo()

	authService := services.NewAuthService(
		userRepo,
		sessionRepo,
		deviceRepo,
		"test-secret",
		15*time.Minute,
		24*time.Hour,
		logger,
	)

	authHandler := handlers.NewAuthHandler(authService, logger)

	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"healthy"}`)
	}).Methods(http.MethodGet)

	v1 := router.PathPrefix("/v1").Subrouter()
	public := v1.PathPrefix("").Subrouter()
	public.HandleFunc("/auth/signup", authHandler.Signup).Methods(http.MethodPost)

	return router
}

type mockUserRepo struct {
	mu       sync.Mutex
	users    map[uuid.UUID]*models.User
	emailMap map[string]uuid.UUID
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:    make(map[uuid.UUID]*models.User),
		emailMap: make(map[string]uuid.UUID),
	}
}

func (m *mockUserRepo) Create(ctx context.Context, user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.emailMap[user.Email]; exists {
		return fmt.Errorf("email already exists")
	}

	copied := *user
	m.users[user.ID] = &copied
	m.emailMap[user.Email] = user.ID
	return nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	copied := *user
	return &copied, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id, ok := m.emailMap[email]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	user := m.users[id]
	copied := *user
	return &copied, nil
}

func (m *mockUserRepo) Update(ctx context.Context, user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.users[user.ID]; !ok {
		return fmt.Errorf("user not found")
	}
	copied := *user
	m.users[user.ID] = &copied
	m.emailMap[user.Email] = user.ID
	return nil
}

func (m *mockUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[id]
	if !ok {
		return fmt.Errorf("user not found")
	}

	copied := *user
	copied.Status = "deleted"
	m.users[id] = &copied
	return nil
}

type mockSessionRepo struct {
	mu           sync.Mutex
	sessions     map[uuid.UUID]*models.Session
	tokenIndex   map[string]uuid.UUID
	refreshIndex map[string]uuid.UUID
}

func newMockSessionRepo() *mockSessionRepo {
	return &mockSessionRepo{
		sessions:     make(map[uuid.UUID]*models.Session),
		tokenIndex:   make(map[string]uuid.UUID),
		refreshIndex: make(map[string]uuid.UUID),
	}
}

func (m *mockSessionRepo) Create(ctx context.Context, session *models.Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	copied := *session
	m.sessions[session.ID] = &copied
	m.tokenIndex[session.Token] = session.ID
	if session.RefreshToken != nil {
		m.refreshIndex[*session.RefreshToken] = session.ID
	}
	return nil
}

func (m *mockSessionRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[id]
	if !ok || !session.IsActive {
		return nil, fmt.Errorf("session not found")
	}
	copied := *session
	return &copied, nil
}

func (m *mockSessionRepo) GetByToken(ctx context.Context, token string) (*models.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id, ok := m.tokenIndex[token]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	return m.GetByID(ctx, id)
}

func (m *mockSessionRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id, ok := m.refreshIndex[refreshToken]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	return m.GetByID(ctx, id)
}

func (m *mockSessionRepo) Update(ctx context.Context, session *models.Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	existing, ok := m.sessions[session.ID]
	if !ok {
		return fmt.Errorf("session not found")
	}

	if existing.RefreshToken != nil {
		delete(m.refreshIndex, *existing.RefreshToken)
	}
	delete(m.tokenIndex, existing.Token)

	copied := *session
	m.sessions[session.ID] = &copied
	m.tokenIndex[session.Token] = session.ID
	if session.RefreshToken != nil {
		m.refreshIndex[*session.RefreshToken] = session.ID
	}
	return nil
}

func (m *mockSessionRepo) Revoke(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[id]
	if !ok {
		return fmt.Errorf("session not found")
	}
	copied := *session
	copied.IsActive = false
	m.sessions[id] = &copied
	return nil
}

func (m *mockSessionRepo) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, session := range m.sessions {
		if session.UserID == userID {
			copied := *session
			copied.IsActive = false
			m.sessions[id] = &copied
		}
	}
	return nil
}

type mockDeviceRepo struct {
	mu          sync.Mutex
	devices     map[uuid.UUID]*models.Device
	userDevices map[uuid.UUID]map[string]uuid.UUID
}

func newMockDeviceRepo() *mockDeviceRepo {
	return &mockDeviceRepo{
		devices:     make(map[uuid.UUID]*models.Device),
		userDevices: make(map[uuid.UUID]map[string]uuid.UUID),
	}
}

func (m *mockDeviceRepo) Create(ctx context.Context, device *models.Device) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	copied := *device
	m.devices[device.ID] = &copied
	if device.DeviceID != "" {
		if m.userDevices[device.UserID] == nil {
			m.userDevices[device.UserID] = make(map[string]uuid.UUID)
		}
		m.userDevices[device.UserID][device.DeviceID] = device.ID
	}
	return nil
}

func (m *mockDeviceRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Device, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	device, ok := m.devices[id]
	if !ok {
		return nil, fmt.Errorf("device not found")
	}
	copied := *device
	return &copied, nil
}

func (m *mockDeviceRepo) GetByDeviceID(ctx context.Context, userID uuid.UUID, deviceID *string) (*models.Device, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if deviceID == nil {
		return nil, fmt.Errorf("device_id is required")
	}

	ids := m.userDevices[userID]
	if ids == nil {
		return nil, nil
	}

	id, ok := ids[*deviceID]
	if !ok {
		return nil, nil
	}
	device := m.devices[id]
	copied := *device
	return &copied, nil
}

func (m *mockDeviceRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Device, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ids := m.userDevices[userID]
	if ids == nil {
		return []*models.Device{}, nil
	}

	devices := make([]*models.Device, 0, len(ids))
	for _, id := range ids {
		device := m.devices[id]
		copied := *device
		devices = append(devices, &copied)
	}
	return devices, nil
}

func (m *mockDeviceRepo) Update(ctx context.Context, device *models.Device) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.devices[device.ID]; !ok {
		return fmt.Errorf("device not found")
	}
	copied := *device
	m.devices[device.ID] = &copied
	return nil
}

func (m *mockDeviceRepo) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	device, ok := m.devices[id]
	if !ok {
		return fmt.Errorf("device not found")
	}

	delete(m.devices, id)
	if device.DeviceID != "" {
		if ids := m.userDevices[device.UserID]; ids != nil {
			delete(ids, device.DeviceID)
		}
	}
	return nil
}
