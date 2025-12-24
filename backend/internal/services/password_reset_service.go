package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
	"base-app-service/pkg/auth"
)

type PasswordResetService struct {
	userRepo           repositories.UserRepository
	passwordResetRepo  repositories.PasswordResetRepository
	logger             *zap.Logger
	tokenExpiry        time.Duration
}

func NewPasswordResetService(
	userRepo repositories.UserRepository,
	passwordResetRepo repositories.PasswordResetRepository,
	logger *zap.Logger,
) *PasswordResetService {
	return &PasswordResetService{
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		logger:            logger,
		tokenExpiry:       1 * time.Hour, // 1 hour expiry
	}
}

// RequestPasswordReset generates a reset token and stores it, returns the token
func (s *PasswordResetService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// Don't reveal if user exists or not (security best practice)
		s.logger.Info("Password reset requested for email", zap.String("email", email))
		return "", nil
	}

	// Generate secure random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Delete any existing tokens for this user
	_ = s.passwordResetRepo.DeleteByUserID(ctx, user.ID)

	// Create new reset token
	resetToken := &models.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(s.tokenExpiry),
		CreatedAt: time.Now(),
	}

	if err := s.passwordResetRepo.Create(ctx, resetToken); err != nil {
		return "", fmt.Errorf("failed to create reset token: %w", err)
	}

	s.logger.Info("Password reset token generated",
		zap.String("user_id", user.ID.String()),
	)

	return token, nil
}

// ResetPassword validates token and resets password
func (s *PasswordResetService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Validate password strength using enhanced validation
	validation := auth.ValidatePassword(newPassword)
	if !validation.Valid {
		return fmt.Errorf("password validation failed: %s", validation.Errors[0])
	}

	// Get reset token
	resetToken, err := s.passwordResetRepo.GetByToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	if resetToken.UsedAt != nil {
		return errors.New("reset token has already been used")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, resetToken.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Hash new password
	passwordHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	user.PasswordHash = passwordHash
	user.PasswordChangedAt = time.Now()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	if err := s.passwordResetRepo.MarkAsUsed(ctx, resetToken.ID); err != nil {
		s.logger.Warn("Failed to mark reset token as used", zap.Error(err))
	}

	// Update password last changed in settings if exists
	// This will be handled in settings service

	s.logger.Info("Password reset successful", zap.String("user_id", user.ID.String()))
	return nil
}

