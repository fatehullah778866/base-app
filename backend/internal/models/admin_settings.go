package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminSettings struct {
	AdminID              uuid.UUID `db:"admin_id" json:"admin_id"`
	DashboardLayout      *string   `db:"dashboard_layout" json:"dashboard_layout"` // JSON
	DefaultPermissions   *string   `db:"default_permissions" json:"default_permissions"` // JSON array
	NotificationPreferences *string `db:"notification_preferences" json:"notification_preferences"` // JSON
	ThemePreferences      *string   `db:"theme_preferences" json:"theme_preferences"` // JSON
	AdminVerificationCode *string   `db:"admin_verification_code" json:"admin_verification_code"` // Verification code for admin creation
	CreatedAt            time.Time `db:"created_at" json:"created_at"`
	UpdatedAt            time.Time `db:"updated_at" json:"updated_at"`
}

type CustomCRUDEntity struct {
	ID          uuid.UUID `db:"id" json:"id"`
	CreatedBy   uuid.UUID `db:"created_by" json:"created_by"`
	EntityName  string    `db:"entity_name" json:"entity_name"` // e.g., "products"
	DisplayName string    `db:"display_name" json:"display_name"`
	Description *string   `db:"description" json:"description"`
	Schema       string    `db:"schema" json:"schema"` // JSON schema
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type CustomCRUDData struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	EntityID  uuid.UUID  `db:"entity_id" json:"entity_id"`
	Data      string     `db:"data" json:"data"` // JSON
	CreatedBy uuid.UUID  `db:"created_by" json:"created_by"`
	UpdatedBy *uuid.UUID `db:"updated_by" json:"updated_by"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type AdminActivityLog struct {
	ID         uuid.UUID `db:"id" json:"id"`
	AdminID    uuid.UUID `db:"admin_id" json:"admin_id"`
	Action     string    `db:"action" json:"action"`
	EntityType string    `db:"entity_type" json:"entity_type"`
	EntityID   *string   `db:"entity_id" json:"entity_id"`
	Details    *string   `db:"details" json:"details"` // JSON
	IPAddress  *string   `db:"ip_address" json:"ip_address"`
	UserAgent  *string   `db:"user_agent" json:"user_agent"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type UserManagementAction struct {
	ID         uuid.UUID `db:"id" json:"id"`
	AdminID    uuid.UUID `db:"admin_id" json:"admin_id"`
	UserID     uuid.UUID `db:"user_id" json:"user_id"`
	ActionType string    `db:"action_type" json:"action_type"`
	Changes    *string   `db:"changes" json:"changes"` // JSON
	Reason     *string   `db:"reason" json:"reason"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type AdminPermission struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	AdminID       uuid.UUID  `db:"admin_id" json:"admin_id"`
	PermissionName string    `db:"permission_name" json:"permission_name"`
	GrantedAt     time.Time  `db:"granted_at" json:"granted_at"`
	GrantedBy     *uuid.UUID `db:"granted_by" json:"granted_by"`
}

// CRUDTemplate represents a CRUD template stored in the database
type CRUDTemplate struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"` // e.g., "portfolio", "visa"
	DisplayName string    `db:"display_name" json:"display_name"`
	Description *string   `db:"description" json:"description"`
	Schema      string    `db:"schema" json:"schema"` // JSON schema
	Icon        *string   `db:"icon" json:"icon"`
	Category    *string   `db:"category" json:"category"`
	CreatedBy   uuid.UUID `db:"created_by" json:"created_by"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	IsSystem    bool      `db:"is_system" json:"is_system"` // System templates cannot be deleted
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

