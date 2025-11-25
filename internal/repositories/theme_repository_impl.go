package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
)

type themeRepository struct {
	db *database.DB
}

func NewThemeRepository(db *database.DB) ThemeRepository {
	return &themeRepository{db: db}
}

func (r *themeRepository) GetGlobalTheme(ctx context.Context, userID uuid.UUID) (*models.ThemePreferences, error) {
	query := `
		SELECT user_id, kompassui_theme, kompassui_contrast, kompassui_text_direction,
			kompassui_brand, theme_synced_at, theme_sync_enabled
		FROM user_settings
		WHERE user_id = $1
	`

	theme := &models.ThemePreferences{}
	err := r.db.DB.QueryRowContext(ctx, query, userID).Scan(
		&theme.UserID, &theme.Theme, &theme.Contrast, &theme.TextDirection,
		&theme.Brand, &theme.SyncedAt, &theme.SyncEnabled,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("theme preferences not found")
	}
	if err != nil {
		return nil, err
	}

	return theme, nil
}

func (r *themeRepository) UpdateGlobalTheme(ctx context.Context, theme *models.ThemePreferences) error {
	// First check if user_settings exists
	var exists bool
	err := r.db.DB.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM user_settings WHERE user_id = $1)", theme.UserID,
	).Scan(&exists)

	if err != nil {
		return err
	}

	if !exists {
		// Insert new settings
		query := `
			INSERT INTO user_settings (
				user_id, kompassui_theme, kompassui_contrast, kompassui_text_direction,
				kompassui_brand, theme_synced_at, theme_sync_enabled, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		`
		_, err = r.db.DB.ExecContext(ctx, query,
			theme.UserID, theme.Theme, theme.Contrast, theme.TextDirection,
			theme.Brand, theme.SyncedAt, theme.SyncEnabled,
		)
	} else {
		// Update existing settings
		query := `
			UPDATE user_settings
			SET kompassui_theme = $2, kompassui_contrast = $3,
				kompassui_text_direction = $4, kompassui_brand = $5,
				theme_synced_at = $6, theme_sync_enabled = $7, updated_at = NOW()
			WHERE user_id = $1
		`
		_, err = r.db.DB.ExecContext(ctx, query,
			theme.UserID, theme.Theme, theme.Contrast, theme.TextDirection,
			theme.Brand, theme.SyncedAt, theme.SyncEnabled,
		)
	}

	return err
}

func (r *themeRepository) GetProductOverride(ctx context.Context, userID uuid.UUID, productName string) (*models.ProductThemeOverride, error) {
	query := `
		SELECT id, user_id, product_name, theme, contrast, text_direction, brand,
			created_at, updated_at
		FROM product_theme_preferences
		WHERE user_id = $1 AND product_name = $2
	`

	override := &models.ProductThemeOverride{}
	err := r.db.DB.QueryRowContext(ctx, query, userID, productName).Scan(
		&override.ID, &override.UserID, &override.ProductName,
		&override.Theme, &override.Contrast, &override.TextDirection, &override.Brand,
		&override.CreatedAt, &override.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Return nil instead of error for "not found"
	}
	if err != nil {
		return nil, err
	}

	return override, nil
}

func (r *themeRepository) UpsertProductOverride(ctx context.Context, override *models.ProductThemeOverride) error {
	query := `
		INSERT INTO product_theme_preferences (
			id, user_id, product_name, theme, contrast, text_direction, brand,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (user_id, product_name)
		DO UPDATE SET
			theme = EXCLUDED.theme,
			contrast = EXCLUDED.contrast,
			text_direction = EXCLUDED.text_direction,
			brand = EXCLUDED.brand,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.DB.ExecContext(ctx, query,
		override.ID, override.UserID, override.ProductName,
		override.Theme, override.Contrast, override.TextDirection, override.Brand,
		override.CreatedAt, override.UpdatedAt,
	)

	return err
}

func (r *themeRepository) DeleteProductOverride(ctx context.Context, userID uuid.UUID, productName string) error {
	query := `DELETE FROM product_theme_preferences WHERE user_id = $1 AND product_name = $2`
	_, err := r.db.DB.ExecContext(ctx, query, userID, productName)
	return err
}

