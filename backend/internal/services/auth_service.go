package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
	"base-app-service/pkg/auth"
)

type AuthService struct {
	userRepo      repositories.UserRepository
	sessionRepo   repositories.SessionRepository
	deviceRepo    repositories.DeviceRepository
	jwtSecret     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	logger        *zap.Logger
}

func NewAuthService(
	userRepo repositories.UserRepository,
	sessionRepo repositories.SessionRepository,
	deviceRepo repositories.DeviceRepository,
	jwtSecret string,
	accessExpiry, refreshExpiry time.Duration,
	logger *zap.Logger,
) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		sessionRepo:   sessionRepo,
		deviceRepo:    deviceRepo,
		jwtSecret:     jwtSecret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		logger:        logger,
	}
}

type SignupRequest struct {
	Email            string
	Password         string
	Name             string
	FirstName        *string
	LastName         *string
	Phone            *string
	SignupSource     *string
	MarketingConsent bool
	TermsAccepted    bool
	TermsVersion     string
	IPAddress        *string
	UserAgent        *string
	DeviceID         *string
	DeviceName       *string
}

type LoginRequest struct {
	Email      string
	Password   string
	RememberMe bool
	DeviceID   *string
	DeviceName *string
	IPAddress  *string
	UserAgent  *string
}

func (s *AuthService) Signup(ctx context.Context, req SignupRequest) (*models.User, *models.Session, error) {
	// Check if user exists
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, nil, errors.New("email already exists")
	}

	// Hash password
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, nil, err
	}

	// Create user
	now := time.Now()
	user := &models.User{
		ID:                uuid.New(),
		Email:             req.Email,
		PasswordHash:      passwordHash,
		Name:              req.Name,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Phone:             req.Phone,
		Status:            "pending",
		Role:              "user",
		SignupSource:      req.SignupSource,
		PasswordChangedAt: now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	// Create session
	session, err := s.createSession(ctx, user, req.IPAddress, req.UserAgent, req.DeviceID, req.DeviceName)
	if err != nil {
		return nil, nil, err
	}

	s.logger.Info("User signed up", zap.String("user_id", user.ID.String()))

	return user, session, nil
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*models.User, *models.Session, bool, error) {
	// Get user
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, nil, false, errors.New("invalid credentials")
	}

	// Check password
	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, nil, false, errors.New("invalid credentials")
	}

	// Check status
	if user.Status != "active" && user.Status != "pending" {
		return nil, nil, false, errors.New("account is not active")
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(ctx, user)

	// Check if device exists
	var device *models.Device
	var isNewDevice bool
	if req.DeviceID != nil {
		device, _ = s.deviceRepo.GetByDeviceID(ctx, user.ID, req.DeviceID)
		isNewDevice = device == nil

		// Create or update device
		if device == nil {
			device = &models.Device{
				ID:         uuid.New(),
				UserID:     user.ID,
				DeviceID:   *req.DeviceID,
				DeviceName: req.DeviceName,
				IPAddress:  req.IPAddress,
				CreatedAt:  now,
				LastUsedAt: now,
			}
			s.deviceRepo.Create(ctx, device)
		} else {
			device.LastUsedAt = now
			s.deviceRepo.Update(ctx, device)
		}
	}

	// Create session
	session, err := s.createSession(ctx, user, req.IPAddress, req.UserAgent, req.DeviceID, req.DeviceName)
	if err != nil {
		return nil, nil, false, err
	}

	s.logger.Info("User logged in", zap.String("user_id", user.ID.String()))

	return user, session, isNewDevice, nil
}

func (s *AuthService) createSession(ctx context.Context, user *models.User, ipAddress, userAgent, deviceID, deviceName *string) (*models.Session, error) {
	sessionID := uuid.New()

	// Generate tokens
	tokenPair, err := auth.GenerateTokenPair(
		user.ID.String(),
		sessionID.String(),
		user.Role,
		s.jwtSecret,
		s.accessExpiry,
		s.refreshExpiry,
	)
	if err != nil {
		return nil, err
	}

	refreshExpiresAt := time.Now().Add(s.refreshExpiry)
	now := time.Now()

	// Create session
	session := &models.Session{
		ID:                    sessionID,
		UserID:                user.ID,
		Token:                 tokenPair.AccessToken,
		RefreshToken:          &tokenPair.RefreshToken,
		RefreshTokenExpiresAt: &refreshExpiresAt,
		DeviceID:              deviceID,
		DeviceName:            deviceName,
		IPAddress:             ipAddress,
		IsActive:              true,
		ExpiresAt:             tokenPair.ExpiresAt,
		CreatedAt:             now,
		LastUsedAt:            now,
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	// Validate refresh token
	claims, err := auth.ValidateToken(refreshToken, s.jwtSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get session
	session, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("session not found")
	}

	// Check if expired
	if session.RefreshTokenExpiresAt != nil && session.RefreshTokenExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	// Generate new tokens
	tokenPair, err := auth.GenerateTokenPair(
		claims.UserID,
		session.ID.String(),
		claims.Role,
		s.jwtSecret,
		s.accessExpiry,
		s.refreshExpiry,
	)
	if err != nil {
		return nil, err
	}

	// Update session
	session.Token = tokenPair.AccessToken
	newRefreshToken := tokenPair.RefreshToken
	session.RefreshToken = &newRefreshToken
	refreshExpiresAt := time.Now().Add(s.refreshExpiry)
	session.RefreshTokenExpiresAt = &refreshExpiresAt
	session.ExpiresAt = tokenPair.ExpiresAt
	session.LastUsedAt = time.Now()

	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID uuid.UUID, revokeAll bool) error {
	if revokeAll {
		// Get user ID from session first
		session, err := s.sessionRepo.GetByID(ctx, sessionID)
		if err != nil {
			return err
		}
		return s.sessionRepo.RevokeAllForUser(ctx, session.UserID)
	}
	return s.sessionRepo.Revoke(ctx, sessionID)
}
