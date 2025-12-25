package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
)

type AdminSettingsService struct {
	settingsRepo repositories.AdminSettingsRepository
	logger       *zap.Logger
}

func NewAdminSettingsService(settingsRepo repositories.AdminSettingsRepository, logger *zap.Logger) *AdminSettingsService {
	return &AdminSettingsService{
		settingsRepo: settingsRepo,
		logger:       logger,
	}
}

func (s *AdminSettingsService) GetSettings(ctx context.Context, adminID uuid.UUID) (*models.AdminSettings, error) {
	settings, err := s.settingsRepo.GetByAdminID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	if settings == nil {
		// Create default settings
		defaultCode := "Kompasstech2025@"
		settings = &models.AdminSettings{
			AdminID:              adminID,
			AdminVerificationCode: &defaultCode,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		if err := s.settingsRepo.Create(ctx, settings); err != nil {
			return nil, err
		}
	}
	return settings, nil
}

// GetSystemVerificationCode gets the system-wide verification code
// Gets the code from the first admin's settings, or returns default
func (s *AdminSettingsService) GetSystemVerificationCode(ctx context.Context) (string, error) {
	// Try to get from first admin's settings
	code, err := s.settingsRepo.GetFirstAdminVerificationCode(ctx)
	if err != nil {
		s.logger.Warn("Failed to get verification code from database, using default", zap.Error(err))
		return "Kompasstech2025@", nil
	}
	// Ensure code is trimmed
	return strings.TrimSpace(code), nil
}

func (s *AdminSettingsService) UpdateSettings(ctx context.Context, adminID uuid.UUID, updates map[string]interface{}) error {
	settings, err := s.GetSettings(ctx, adminID)
	if err != nil {
		return err
	}

	if layout, ok := updates["dashboard_layout"].(string); ok {
		settings.DashboardLayout = &layout
	}
	if permissions, ok := updates["default_permissions"].(string); ok {
		settings.DefaultPermissions = &permissions
	}
	if notifications, ok := updates["notification_preferences"].(string); ok {
		settings.NotificationPreferences = &notifications
	}
	if theme, ok := updates["theme_preferences"].(string); ok {
		settings.ThemePreferences = &theme
	}
	if code, ok := updates["admin_verification_code"].(string); ok {
		trimmedCode := strings.TrimSpace(code)
		settings.AdminVerificationCode = &trimmedCode
	}

	settings.UpdatedAt = time.Now()
	return s.settingsRepo.Update(ctx, settings)
}

