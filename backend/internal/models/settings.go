package models

import (
	"time"

	"github.com/google/uuid"
)

// ComprehensiveSettings represents all user settings in one model
type ComprehensiveSettings struct {
	UserID uuid.UUID `db:"user_id" json:"user_id"`
	
	// Profile Settings
	Username    *string `db:"username" json:"username"`
	DisplayName *string `db:"display_name" json:"display_name"`
	Bio         *string `db:"bio" json:"bio"`
	DateOfBirth *string `db:"date_of_birth" json:"date_of_birth"`
	
	// Security Settings
	TwoFactorEnabled   bool     `db:"two_factor_enabled" json:"two_factor_enabled"`
	TwoFactorSecret    *string  `db:"two_factor_secret" json:"-"` // Hidden from JSON
	TwoFactorBackupCodes *string `db:"two_factor_backup_codes" json:"-"` // Hidden from JSON
	SecurityQuestions  *string  `db:"security_questions" json:"security_questions"` // JSON string
	PasswordLastChanged *time.Time `db:"password_last_changed" json:"password_last_changed"`
	
	// Privacy Settings
	ProfileVisibility string `db:"profile_visibility" json:"profile_visibility"` // public, private, friends
	EmailVisibility   string `db:"email_visibility" json:"email_visibility"`
	PhoneVisibility   string `db:"phone_visibility" json:"phone_visibility"`
	AllowMessaging    string `db:"allow_messaging" json:"allow_messaging"` // everyone, friends, none
	SearchVisibility  bool   `db:"search_visibility" json:"search_visibility"`
	DataSharingEnabled bool   `db:"data_sharing_enabled" json:"data_sharing_enabled"`
	
	// Notification Settings
	EmailNotifications  bool `db:"email_notifications" json:"email_notifications"`
	SMSNotifications    bool `db:"sms_notifications" json:"sms_notifications"`
	PushNotifications   bool `db:"push_notifications" json:"push_notifications"`
	NotificationMessages bool `db:"notification_messages" json:"notification_messages"`
	NotificationAlerts  bool `db:"notification_alerts" json:"notification_alerts"`
	NotificationPromotions bool `db:"notification_promotions" json:"notification_promotions"`
	NotificationSecurity bool `db:"notification_security" json:"notification_security"`
	
	// Account Preferences
	Language     string `db:"language" json:"language"`
	Timezone     string `db:"timezone" json:"timezone"`
	Theme        string `db:"theme" json:"theme"` // light, dark, auto
	FontSize     string `db:"font_size" json:"font_size"` // small, medium, large
	HighContrast bool   `db:"high_contrast" json:"high_contrast"`
	ReducedMotion bool  `db:"reduced_motion" json:"reduced_motion"`
	ScreenReader  bool  `db:"screen_reader" json:"screen_reader"`
	
	// Connected Accounts (JSON string)
	ConnectedAccounts *string `db:"connected_accounts" json:"connected_accounts"` // JSON array
	
	// Data & Account Control
	AccountDeletionRequested bool       `db:"account_deletion_requested" json:"account_deletion_requested"`
	AccountDeletionScheduledAt *time.Time `db:"account_deletion_scheduled_at" json:"account_deletion_scheduled_at"`
	AccountDeactivated        bool       `db:"account_deactivated" json:"account_deactivated"`
	AccountDeactivatedAt      *time.Time `db:"account_deactivated_at" json:"account_deactivated_at"`
	
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// ConnectedAccount represents a connected third-party account
type ConnectedAccount struct {
	Provider    string    `json:"provider"` // google, facebook, apple, etc.
	Email       string    `json:"email"`
	ConnectedAt time.Time `json:"connected_at"`
}

