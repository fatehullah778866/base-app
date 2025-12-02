package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
)

type ThemeService struct {
	themeRepo repositories.ThemeRepository
	logger    *zap.Logger
}

func NewThemeService(themeRepo repositories.ThemeRepository, logger *zap.Logger) *ThemeService {
	return &ThemeService{
		themeRepo: themeRepo,
		logger:    logger,
	}
}

type ThemeUpdateRequest struct {
	Theme         *string
	Contrast      *string
	TextDirection *string
	Brand         *string
}

func (s *ThemeService) GetTheme(ctx context.Context, userID uuid.UUID, productName *string) (*models.ThemePreferences, error) {
	if productName != nil {
		override, err := s.themeRepo.GetProductOverride(ctx, userID, *productName)
		if err == nil && override != nil {
			return &models.ThemePreferences{
				UserID:        userID,
				Theme:         override.Theme,
				Contrast:      override.Contrast,
				TextDirection: override.TextDirection,
				Brand:         override.Brand,
				SyncedAt:      override.UpdatedAt,
				SyncEnabled:   true,
			}, nil
		}
	}

	return s.themeRepo.GetGlobalTheme(ctx, userID)
}

func (s *ThemeService) UpdateTheme(ctx context.Context, userID uuid.UUID, req ThemeUpdateRequest) (*models.ThemePreferences, error) {
	theme, err := s.themeRepo.GetGlobalTheme(ctx, userID)
	if err != nil {
		// Create if doesn't exist
		theme = &models.ThemePreferences{
			UserID:        userID,
			Theme:         "auto",
			Contrast:      "standard",
			TextDirection: "auto",
			SyncEnabled:   true,
		}
	}

	if req.Theme != nil {
		theme.Theme = *req.Theme
	}
	if req.Contrast != nil {
		theme.Contrast = *req.Contrast
	}
	if req.TextDirection != nil {
		theme.TextDirection = *req.TextDirection
	}
	if req.Brand != nil {
		theme.Brand = req.Brand
	}

	theme.SyncedAt = time.Now()

	if err := s.themeRepo.UpdateGlobalTheme(ctx, theme); err != nil {
		return nil, err
	}

	s.logger.Info("Theme updated", zap.String("user_id", userID.String()))

	return theme, nil
}

func (s *ThemeService) SyncTheme(ctx context.Context, userID uuid.UUID, clientTheme *models.ThemePreferences) (*models.ThemePreferences, []string, error) {
	serverTheme, err := s.themeRepo.GetGlobalTheme(ctx, userID)
	if err != nil {
		// Create if doesn't exist
		serverTheme = clientTheme
		serverTheme.SyncedAt = time.Now()
		if err := s.themeRepo.UpdateGlobalTheme(ctx, serverTheme); err != nil {
			return nil, nil, err
		}
		return serverTheme, []string{}, nil
	}

	// Detect conflicts
	conflicts := []string{}
	if clientTheme.Theme != serverTheme.Theme && serverTheme.SyncedAt.After(clientTheme.SyncedAt) {
		conflicts = append(conflicts, "theme")
	}
	if clientTheme.Contrast != serverTheme.Contrast && serverTheme.SyncedAt.After(clientTheme.SyncedAt) {
		conflicts = append(conflicts, "contrast")
	}
	if clientTheme.TextDirection != serverTheme.TextDirection && serverTheme.SyncedAt.After(clientTheme.SyncedAt) {
		conflicts = append(conflicts, "text_direction")
	}

	// If no conflicts, update server with client data
	if len(conflicts) == 0 {
		serverTheme.Theme = clientTheme.Theme
		serverTheme.Contrast = clientTheme.Contrast
		serverTheme.TextDirection = clientTheme.TextDirection
		serverTheme.Brand = clientTheme.Brand
		serverTheme.SyncedAt = time.Now()

		if err := s.themeRepo.UpdateGlobalTheme(ctx, serverTheme); err != nil {
			return nil, nil, err
		}
	}

	return serverTheme, conflicts, nil
}

func (s *ThemeService) SetProductOverride(ctx context.Context, userID uuid.UUID, productName string, req ThemeUpdateRequest) (*models.ProductThemeOverride, error) {
	override, err := s.themeRepo.GetProductOverride(ctx, userID, productName)
	if err != nil {
		// Create new override
		override = &models.ProductThemeOverride{
			ID:            uuid.New(),
			UserID:        userID,
			ProductName:   productName,
			Theme:         "auto",
			Contrast:      "standard",
			TextDirection: "auto",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
	}

	if req.Theme != nil {
		override.Theme = *req.Theme
	}
	if req.Contrast != nil {
		override.Contrast = *req.Contrast
	}
	if req.TextDirection != nil {
		override.TextDirection = *req.TextDirection
	}
	if req.Brand != nil {
		override.Brand = req.Brand
	}

	override.UpdatedAt = time.Now()

	if err := s.themeRepo.UpsertProductOverride(ctx, override); err != nil {
		return nil, err
	}

	return override, nil
}

func (s *ThemeService) GetProductOverride(ctx context.Context, userID uuid.UUID, productName string) (*models.ProductThemeOverride, error) {
	return s.themeRepo.GetProductOverride(ctx, userID, productName)
}

func (s *ThemeService) RemoveProductOverride(ctx context.Context, userID uuid.UUID, productName string) error {
	return s.themeRepo.DeleteProductOverride(ctx, userID, productName)
}
