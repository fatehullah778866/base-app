package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
)

type SettingsService struct {
	settingsRepo repositories.SettingsRepository
	userRepo     repositories.UserRepository
	logger       *zap.Logger
}

func NewSettingsService(
	settingsRepo repositories.SettingsRepository,
	userRepo repositories.UserRepository,
	logger *zap.Logger,
) *SettingsService {
	return &SettingsService{
		settingsRepo: settingsRepo,
		userRepo:     userRepo,
		logger:       logger,
	}
}

// GetSettings retrieves user settings, creating defaults if not exists
func (s *SettingsService) GetSettings(ctx context.Context, userID uuid.UUID) (*models.ComprehensiveSettings, error) {
	settings, err := s.settingsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create default settings if not exists
	if settings == nil {
		settings = &models.ComprehensiveSettings{
			UserID:            userID,
			ProfileVisibility: "public",
			EmailVisibility:   "private",
			PhoneVisibility:   "private",
			AllowMessaging:    "everyone",
			SearchVisibility:  true,
			EmailNotifications: true,
			PushNotifications:  true,
			NotificationMessages: true,
			NotificationAlerts:   true,
			NotificationSecurity: true,
			Language:          "en",
			Timezone:          "UTC",
			Theme:             "light",
			FontSize:          "medium",
			UpdatedAt:         time.Now(),
		}
		if err := s.settingsRepo.Create(ctx, settings); err != nil {
			return nil, err
		}
	}

	return settings, nil
}

// UpdateProfileSettings updates profile-related settings
func (s *SettingsService) UpdateProfileSettings(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	settings, err := s.GetSettings(ctx, userID)
	if err != nil {
		return err
	}

	// Update allowed profile fields
	if username, ok := updates["username"].(string); ok && username != "" {
		// Check if username is already taken
		existing, _ := s.settingsRepo.GetUsername(ctx, username)
		if existing != nil && existing.UserID != userID {
			return errors.New("username already taken")
		}
		settings.Username = &username
	}
	if displayName, ok := updates["display_name"].(string); ok {
		settings.DisplayName = &displayName
	}
	if bio, ok := updates["bio"].(string); ok {
		settings.Bio = &bio
	}
	if dateOfBirth, ok := updates["date_of_birth"].(string); ok {
		settings.DateOfBirth = &dateOfBirth
	}

	return s.settingsRepo.Update(ctx, settings)
}

// UpdateSecuritySettings updates security-related settings
func (s *SettingsService) UpdateSecuritySettings(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	settings, err := s.GetSettings(ctx, userID)
	if err != nil {
		return err
	}

	if twoFactorEnabled, ok := updates["two_factor_enabled"].(bool); ok {
		settings.TwoFactorEnabled = twoFactorEnabled
	}
	if securityQuestions, ok := updates["security_questions"].(string); ok {
		settings.SecurityQuestions = &securityQuestions
	}

	now := time.Now()
	settings.PasswordLastChanged = &now

	return s.settingsRepo.Update(ctx, settings)
}

// UpdatePrivacySettings updates privacy-related settings
func (s *SettingsService) UpdatePrivacySettings(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	settings, err := s.GetSettings(ctx, userID)
	if err != nil {
		return err
	}

	if profileVisibility, ok := updates["profile_visibility"].(string); ok {
		settings.ProfileVisibility = profileVisibility
	}
	if emailVisibility, ok := updates["email_visibility"].(string); ok {
		settings.EmailVisibility = emailVisibility
	}
	if phoneVisibility, ok := updates["phone_visibility"].(string); ok {
		settings.PhoneVisibility = phoneVisibility
	}
	if allowMessaging, ok := updates["allow_messaging"].(string); ok {
		settings.AllowMessaging = allowMessaging
	}
	if searchVisibility, ok := updates["search_visibility"].(bool); ok {
		settings.SearchVisibility = searchVisibility
	}
	if dataSharingEnabled, ok := updates["data_sharing_enabled"].(bool); ok {
		settings.DataSharingEnabled = dataSharingEnabled
	}

	return s.settingsRepo.Update(ctx, settings)
}

// UpdateNotificationSettings updates notification-related settings
func (s *SettingsService) UpdateNotificationSettings(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	settings, err := s.GetSettings(ctx, userID)
	if err != nil {
		return err
	}

	if emailNotifications, ok := updates["email_notifications"].(bool); ok {
		settings.EmailNotifications = emailNotifications
	}
	if smsNotifications, ok := updates["sms_notifications"].(bool); ok {
		settings.SMSNotifications = smsNotifications
	}
	if pushNotifications, ok := updates["push_notifications"].(bool); ok {
		settings.PushNotifications = pushNotifications
	}
	if notificationMessages, ok := updates["notification_messages"].(bool); ok {
		settings.NotificationMessages = notificationMessages
	}
	if notificationAlerts, ok := updates["notification_alerts"].(bool); ok {
		settings.NotificationAlerts = notificationAlerts
	}
	if notificationPromotions, ok := updates["notification_promotions"].(bool); ok {
		settings.NotificationPromotions = notificationPromotions
	}
	if notificationSecurity, ok := updates["notification_security"].(bool); ok {
		settings.NotificationSecurity = notificationSecurity
	}

	return s.settingsRepo.Update(ctx, settings)
}

// UpdateAccountPreferences updates account preference settings
func (s *SettingsService) UpdateAccountPreferences(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	settings, err := s.GetSettings(ctx, userID)
	if err != nil {
		return err
	}

	if language, ok := updates["language"].(string); ok {
		settings.Language = language
	}
	if timezone, ok := updates["timezone"].(string); ok {
		settings.Timezone = timezone
	}
	if theme, ok := updates["theme"].(string); ok {
		settings.Theme = theme
	}
	if fontSize, ok := updates["font_size"].(string); ok {
		settings.FontSize = fontSize
	}
	if highContrast, ok := updates["high_contrast"].(bool); ok {
		settings.HighContrast = highContrast
	}
	if reducedMotion, ok := updates["reduced_motion"].(bool); ok {
		settings.ReducedMotion = reducedMotion
	}
	if screenReader, ok := updates["screen_reader"].(bool); ok {
		settings.ScreenReader = screenReader
	}

	return s.settingsRepo.Update(ctx, settings)
}

// AddConnectedAccount adds a connected third-party account
func (s *SettingsService) AddConnectedAccount(ctx context.Context, userID uuid.UUID, account models.ConnectedAccount) error {
	settings, err := s.GetSettings(ctx, userID)
	if err != nil {
		return err
	}

	var accounts []models.ConnectedAccount
	if settings.ConnectedAccounts != nil {
		if err := json.Unmarshal([]byte(*settings.ConnectedAccounts), &accounts); err != nil {
			accounts = []models.ConnectedAccount{}
		}
	}

	// Check if account already exists
	for i, acc := range accounts {
		if acc.Provider == account.Provider {
			accounts[i] = account // Update existing
			accountsJSON, _ := json.Marshal(accounts)
			accountsStr := string(accountsJSON)
			settings.ConnectedAccounts = &accountsStr
			return s.settingsRepo.Update(ctx, settings)
		}
	}

	// Add new account
	accounts = append(accounts, account)
	accountsJSON, _ := json.Marshal(accounts)
	accountsStr := string(accountsJSON)
	settings.ConnectedAccounts = &accountsStr

	return s.settingsRepo.Update(ctx, settings)
}

// RemoveConnectedAccount removes a connected account
func (s *SettingsService) RemoveConnectedAccount(ctx context.Context, userID uuid.UUID, provider string) error {
	settings, err := s.GetSettings(ctx, userID)
	if err != nil {
		return err
	}

	if settings.ConnectedAccounts == nil {
		return nil
	}

	var accounts []models.ConnectedAccount
	if err := json.Unmarshal([]byte(*settings.ConnectedAccounts), &accounts); err != nil {
		return err
	}

	// Remove account
	var filtered []models.ConnectedAccount
	for _, acc := range accounts {
		if acc.Provider != provider {
			filtered = append(filtered, acc)
		}
	}

	accountsJSON, _ := json.Marshal(filtered)
	accountsStr := string(accountsJSON)
	settings.ConnectedAccounts = &accountsStr

	return s.settingsRepo.Update(ctx, settings)
}

// RequestAccountDeletion schedules account deletion
func (s *SettingsService) RequestAccountDeletion(ctx context.Context, userID uuid.UUID, daysUntilDeletion int) error {
	settings, err := s.GetSettings(ctx, userID)
	if err != nil {
		return err
	}

	settings.AccountDeletionRequested = true
	scheduledAt := time.Now().AddDate(0, 0, daysUntilDeletion)
	settings.AccountDeletionScheduledAt = &scheduledAt

	return s.settingsRepo.Update(ctx, settings)
}

// DeactivateAccount temporarily deactivates account
func (s *SettingsService) DeactivateAccount(ctx context.Context, userID uuid.UUID) error {
	settings, err := s.GetSettings(ctx, userID)
	if err != nil {
		return err
	}

	settings.AccountDeactivated = true
	now := time.Now()
	settings.AccountDeactivatedAt = &now

	return s.settingsRepo.Update(ctx, settings)
}

// ReactivateAccount reactivates a deactivated account
func (s *SettingsService) ReactivateAccount(ctx context.Context, userID uuid.UUID) error {
	settings, err := s.GetSettings(ctx, userID)
	if err != nil {
		return err
	}

	settings.AccountDeactivated = false
	settings.AccountDeactivatedAt = nil

	return s.settingsRepo.Update(ctx, settings)
}

