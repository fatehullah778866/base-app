package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
	"base-app-service/internal/services"
)

// Mock repositories
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, search string) ([]*models.User, error) {
	args := m.Called(ctx, search)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) SetStatus(ctx context.Context, id uuid.UUID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockUserRepository) MarkDeleted(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) PurgeDeletedBefore(ctx context.Context, cutoff time.Time) error {
	args := m.Called(ctx, cutoff)
	return args.Error(0)
}

func TestAuthService_Signup(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	userRepo := new(MockUserRepository)
	sessionRepo := new(MockSessionRepository)
	deviceRepo := new(MockDeviceRepository)

	authService := services.NewAuthService(
		userRepo,
		sessionRepo,
		deviceRepo,
		"test-secret",
		15*time.Minute,
		30*24*time.Hour,
		logger,
	)

	t.Run("successful signup", func(t *testing.T) {
		ctx := context.Background()
		req := services.SignupRequest{
			Email:    "test@example.com",
			Password: "Test123!@#",
			Name:     "Test User",
		}

		userRepo.On("GetByEmail", ctx, req.Email).Return(nil, nil)
		userRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Return(nil)
		sessionRepo.On("Create", ctx, mock.AnythingOfType("*models.Session")).Return(nil)
		deviceRepo.On("CreateOrUpdate", ctx, mock.AnythingOfType("*models.Device")).Return(nil)

		user, session, err := authService.Signup(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotNil(t, session)
		assert.Equal(t, req.Email, user.Email)
		userRepo.AssertExpectations(t)
	})

	t.Run("email already exists", func(t *testing.T) {
		ctx := context.Background()
		req := services.SignupRequest{
			Email:    "existing@example.com",
			Password: "Test123!@#",
			Name:     "Test User",
		}

		existingUser := &models.User{Email: req.Email}
		userRepo.On("GetByEmail", ctx, req.Email).Return(existingUser, nil)

		user, session, err := authService.Signup(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Nil(t, session)
		assert.Contains(t, err.Error(), "already exists")
	})
}

// Mock repositories (simplified - you'd need to implement all methods)
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Create(ctx context.Context, session *models.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

type MockDeviceRepository struct {
	mock.Mock
}

func (m *MockDeviceRepository) CreateOrUpdate(ctx context.Context, device *models.Device) error {
	args := m.Called(ctx, device)
	return args.Error(0)
}

