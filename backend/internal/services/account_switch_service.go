package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
)

type AccountSwitchService struct {
	accountSwitchRepo repositories.AccountSwitchRepository
	userRepo          repositories.UserRepository
	logger            *zap.Logger
}

func NewAccountSwitchService(
	accountSwitchRepo repositories.AccountSwitchRepository,
	userRepo repositories.UserRepository,
	logger *zap.Logger,
) *AccountSwitchService {
	return &AccountSwitchService{
		accountSwitchRepo: accountSwitchRepo,
		userRepo:          userRepo,
		logger:            logger,
	}
}

// SwitchAccount switches user context (for multi-account or role switching)
func (s *AccountSwitchService) SwitchAccount(ctx context.Context, userID uuid.UUID, switchedToUserID *uuid.UUID, switchedToRole *string, reason *string, ipAddress *string, userAgent *string) (*models.AccountSwitch, error) {
	// Get current user to determine current role
	currentUser, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	var switchedFromRole *string
	if currentUser != nil {
		switchedFromRole = &currentUser.Role
	}

	// If switching to another user, verify it exists
	if switchedToUserID != nil {
		targetUser, err := s.userRepo.GetByID(ctx, *switchedToUserID)
		if err != nil || targetUser == nil {
			return nil, errors.New("target user not found")
		}
	}

	switchRecord := &models.AccountSwitch{
		ID:               uuid.New(),
		UserID:           userID,
		SwitchedToUserID: switchedToUserID,
		SwitchedToRole:   switchedToRole,
		SwitchedFromRole: switchedFromRole,
		Reason:           reason,
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
		CreatedAt:        time.Now(),
	}

	if err := s.accountSwitchRepo.Create(ctx, switchRecord); err != nil {
		return nil, err
	}

	return switchRecord, nil
}

// GetSwitchHistory retrieves account switch history for a user
func (s *AccountSwitchService) GetSwitchHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*models.AccountSwitch, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.accountSwitchRepo.GetByUserID(ctx, userID, limit)
}

