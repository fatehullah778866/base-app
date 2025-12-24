package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
)

type settingsRepository struct {
	db *database.DB
}

func NewSettingsRepository(db *database.DB) SettingsRepository {
	return &settingsRepository{db: db}
}

func (r *settingsRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.ComprehensiveSettings, error) {
	var s models.ComprehensiveSettings
	var username, displayName, bio, dateOfBirth sql.NullString
	var twoFactorSecret, twoFactorBackupCodes, securityQuestions sql.NullString
	var passwordLastChanged, accountDeletionScheduledAt, accountDeactivatedAt sql.NullTime
	var connectedAccounts sql.NullString

	query := `
		SELECT user_id, username, display_name, bio, date_of_birth,
		       two_factor_enabled, two_factor_secret, two_factor_backup_codes,
		       security_questions, password_last_changed,
		       profile_visibility, email_visibility, phone_visibility,
		       allow_messaging, search_visibility, data_sharing_enabled,
		       email_notifications, sms_notifications, push_notifications,
		       notification_messages, notification_alerts, notification_promotions,
		       notification_security, language, timezone, theme, font_size,
		       high_contrast, reduced_motion, screen_reader,
		       connected_accounts, account_deletion_requested,
		       account_deletion_scheduled_at, account_deactivated,
		       account_deactivated_at, updated_at
		FROM user_settings_comprehensive
		WHERE user_id = ?
	`
	err := r.db.QueryRowContext(ctx, query, userID.String()).Scan(
		&s.UserID,
		&username, &displayName, &bio, &dateOfBirth,
		&s.TwoFactorEnabled, &twoFactorSecret, &twoFactorBackupCodes,
		&securityQuestions, &passwordLastChanged,
		&s.ProfileVisibility, &s.EmailVisibility, &s.PhoneVisibility,
		&s.AllowMessaging, &s.SearchVisibility, &s.DataSharingEnabled,
		&s.EmailNotifications, &s.SMSNotifications, &s.PushNotifications,
		&s.NotificationMessages, &s.NotificationAlerts, &s.NotificationPromotions,
		&s.NotificationSecurity, &s.Language, &s.Timezone, &s.Theme, &s.FontSize,
		&s.HighContrast, &s.ReducedMotion, &s.ScreenReader,
		&connectedAccounts, &s.AccountDeletionRequested,
		&accountDeletionScheduledAt, &s.AccountDeactivated,
		&accountDeactivatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if username.Valid {
		s.Username = &username.String
	}
	if displayName.Valid {
		s.DisplayName = &displayName.String
	}
	if bio.Valid {
		s.Bio = &bio.String
	}
	if dateOfBirth.Valid {
		s.DateOfBirth = &dateOfBirth.String
	}
	if twoFactorSecret.Valid {
		s.TwoFactorSecret = &twoFactorSecret.String
	}
	if twoFactorBackupCodes.Valid {
		s.TwoFactorBackupCodes = &twoFactorBackupCodes.String
	}
	if securityQuestions.Valid {
		s.SecurityQuestions = &securityQuestions.String
	}
	if passwordLastChanged.Valid {
		s.PasswordLastChanged = &passwordLastChanged.Time
	}
	if accountDeletionScheduledAt.Valid {
		s.AccountDeletionScheduledAt = &accountDeletionScheduledAt.Time
	}
	if accountDeactivatedAt.Valid {
		s.AccountDeactivatedAt = &accountDeactivatedAt.Time
	}
	if connectedAccounts.Valid {
		s.ConnectedAccounts = &connectedAccounts.String
	}

	return &s, nil
}

func (r *settingsRepository) Create(ctx context.Context, settings *models.ComprehensiveSettings) error {
	query := `
		INSERT INTO user_settings_comprehensive (
			user_id, username, display_name, bio, date_of_birth,
			two_factor_enabled, two_factor_secret, two_factor_backup_codes,
			security_questions, password_last_changed,
			profile_visibility, email_visibility, phone_visibility,
			allow_messaging, search_visibility, data_sharing_enabled,
			email_notifications, sms_notifications, push_notifications,
			notification_messages, notification_alerts, notification_promotions,
			notification_security, language, timezone, theme, font_size,
			high_contrast, reduced_motion, screen_reader,
			connected_accounts, account_deletion_requested,
			account_deletion_scheduled_at, account_deactivated,
			account_deactivated_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		settings.UserID.String(),
		settings.Username, settings.DisplayName, settings.Bio, settings.DateOfBirth,
		settings.TwoFactorEnabled, settings.TwoFactorSecret, settings.TwoFactorBackupCodes,
		settings.SecurityQuestions, settings.PasswordLastChanged,
		settings.ProfileVisibility, settings.EmailVisibility, settings.PhoneVisibility,
		settings.AllowMessaging, settings.SearchVisibility, settings.DataSharingEnabled,
		settings.EmailNotifications, settings.SMSNotifications, settings.PushNotifications,
		settings.NotificationMessages, settings.NotificationAlerts, settings.NotificationPromotions,
		settings.NotificationSecurity, settings.Language, settings.Timezone, settings.Theme, settings.FontSize,
		settings.HighContrast, settings.ReducedMotion, settings.ScreenReader,
		settings.ConnectedAccounts, settings.AccountDeletionRequested,
		settings.AccountDeletionScheduledAt, settings.AccountDeactivated,
		settings.AccountDeactivatedAt, time.Now(),
	)
	return err
}

func (r *settingsRepository) Update(ctx context.Context, settings *models.ComprehensiveSettings) error {
	query := `
		UPDATE user_settings_comprehensive SET
			username = ?, display_name = ?, bio = ?, date_of_birth = ?,
			two_factor_enabled = ?, two_factor_secret = ?, two_factor_backup_codes = ?,
			security_questions = ?, password_last_changed = ?,
			profile_visibility = ?, email_visibility = ?, phone_visibility = ?,
			allow_messaging = ?, search_visibility = ?, data_sharing_enabled = ?,
			email_notifications = ?, sms_notifications = ?, push_notifications = ?,
			notification_messages = ?, notification_alerts = ?, notification_promotions = ?,
			notification_security = ?, language = ?, timezone = ?, theme = ?, font_size = ?,
			high_contrast = ?, reduced_motion = ?, screen_reader = ?,
			connected_accounts = ?, account_deletion_requested = ?,
			account_deletion_scheduled_at = ?, account_deactivated = ?,
			account_deactivated_at = ?, updated_at = ?
		WHERE user_id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		settings.Username, settings.DisplayName, settings.Bio, settings.DateOfBirth,
		settings.TwoFactorEnabled, settings.TwoFactorSecret, settings.TwoFactorBackupCodes,
		settings.SecurityQuestions, settings.PasswordLastChanged,
		settings.ProfileVisibility, settings.EmailVisibility, settings.PhoneVisibility,
		settings.AllowMessaging, settings.SearchVisibility, settings.DataSharingEnabled,
		settings.EmailNotifications, settings.SMSNotifications, settings.PushNotifications,
		settings.NotificationMessages, settings.NotificationAlerts, settings.NotificationPromotions,
		settings.NotificationSecurity, settings.Language, settings.Timezone, settings.Theme, settings.FontSize,
		settings.HighContrast, settings.ReducedMotion, settings.ScreenReader,
		settings.ConnectedAccounts, settings.AccountDeletionRequested,
		settings.AccountDeletionScheduledAt, settings.AccountDeactivated,
		settings.AccountDeactivatedAt, time.Now(),
		settings.UserID.String(),
	)
	return err
}

func (r *settingsRepository) GetUsername(ctx context.Context, username string) (*models.ComprehensiveSettings, error) {
	var s models.ComprehensiveSettings
	var userIDStr string
	var usernameVal, displayName, bio, dateOfBirth sql.NullString
	var twoFactorSecret, twoFactorBackupCodes, securityQuestions sql.NullString
	var passwordLastChanged, accountDeletionScheduledAt, accountDeactivatedAt sql.NullTime
	var connectedAccounts sql.NullString

	query := `
		SELECT user_id, username, display_name, bio, date_of_birth,
		       two_factor_enabled, two_factor_secret, two_factor_backup_codes,
		       security_questions, password_last_changed,
		       profile_visibility, email_visibility, phone_visibility,
		       allow_messaging, search_visibility, data_sharing_enabled,
		       email_notifications, sms_notifications, push_notifications,
		       notification_messages, notification_alerts, notification_promotions,
		       notification_security, language, timezone, theme, font_size,
		       high_contrast, reduced_motion, screen_reader,
		       connected_accounts, account_deletion_requested,
		       account_deletion_scheduled_at, account_deactivated,
		       account_deactivated_at, updated_at
		FROM user_settings_comprehensive
		WHERE username = ?
	`
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&userIDStr,
		&usernameVal, &displayName, &bio, &dateOfBirth,
		&s.TwoFactorEnabled, &twoFactorSecret, &twoFactorBackupCodes,
		&securityQuestions, &passwordLastChanged,
		&s.ProfileVisibility, &s.EmailVisibility, &s.PhoneVisibility,
		&s.AllowMessaging, &s.SearchVisibility, &s.DataSharingEnabled,
		&s.EmailNotifications, &s.SMSNotifications, &s.PushNotifications,
		&s.NotificationMessages, &s.NotificationAlerts, &s.NotificationPromotions,
		&s.NotificationSecurity, &s.Language, &s.Timezone, &s.Theme, &s.FontSize,
		&s.HighContrast, &s.ReducedMotion, &s.ScreenReader,
		&connectedAccounts, &s.AccountDeletionRequested,
		&accountDeletionScheduledAt, &s.AccountDeactivated,
		&accountDeactivatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	userID, _ := uuid.Parse(userIDStr)
	s.UserID = userID
	if usernameVal.Valid {
		s.Username = &usernameVal.String
	}
	if displayName.Valid {
		s.DisplayName = &displayName.String
	}
	if bio.Valid {
		s.Bio = &bio.String
	}
	if dateOfBirth.Valid {
		s.DateOfBirth = &dateOfBirth.String
	}
	if twoFactorSecret.Valid {
		s.TwoFactorSecret = &twoFactorSecret.String
	}
	if twoFactorBackupCodes.Valid {
		s.TwoFactorBackupCodes = &twoFactorBackupCodes.String
	}
	if securityQuestions.Valid {
		s.SecurityQuestions = &securityQuestions.String
	}
	if passwordLastChanged.Valid {
		s.PasswordLastChanged = &passwordLastChanged.Time
	}
	if accountDeletionScheduledAt.Valid {
		s.AccountDeletionScheduledAt = &accountDeletionScheduledAt.Time
	}
	if accountDeactivatedAt.Valid {
		s.AccountDeactivatedAt = &accountDeactivatedAt.Time
	}
	if connectedAccounts.Valid {
		s.ConnectedAccounts = &connectedAccounts.String
	}

	return &s, nil
}

