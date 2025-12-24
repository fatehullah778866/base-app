package repositories

import (
	"context"

	"github.com/google/uuid"

	"base-app-service/internal/models"
)

type AdminSettingsRepository interface {
	GetByAdminID(ctx context.Context, adminID uuid.UUID) (*models.AdminSettings, error)
	Create(ctx context.Context, settings *models.AdminSettings) error
	Update(ctx context.Context, settings *models.AdminSettings) error
}

type CustomCRUDRepository interface {
	CreateEntity(ctx context.Context, entity *models.CustomCRUDEntity) error
	GetEntityByID(ctx context.Context, id uuid.UUID) (*models.CustomCRUDEntity, error)
	GetEntityByName(ctx context.Context, name string) (*models.CustomCRUDEntity, error)
	ListEntities(ctx context.Context, createdBy *uuid.UUID, activeOnly bool) ([]*models.CustomCRUDEntity, error)
	UpdateEntity(ctx context.Context, entity *models.CustomCRUDEntity) error
	DeleteEntity(ctx context.Context, id uuid.UUID) error
	
	CreateData(ctx context.Context, data *models.CustomCRUDData) error
	GetDataByID(ctx context.Context, id uuid.UUID) (*models.CustomCRUDData, error)
	ListDataByEntity(ctx context.Context, entityID uuid.UUID, limit int, offset int) ([]*models.CustomCRUDData, error)
	UpdateData(ctx context.Context, data *models.CustomCRUDData) error
	DeleteData(ctx context.Context, id uuid.UUID) error
}

type AdminActivityLogRepository interface {
	Create(ctx context.Context, log *models.AdminActivityLog) error
	GetByAdminID(ctx context.Context, adminID uuid.UUID, limit int) ([]*models.AdminActivityLog, error)
	GetByEntityType(ctx context.Context, entityType string, limit int) ([]*models.AdminActivityLog, error)
}

type UserManagementActionRepository interface {
	Create(ctx context.Context, action *models.UserManagementAction) error
	GetByAdminID(ctx context.Context, adminID uuid.UUID, limit int) ([]*models.UserManagementAction, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*models.UserManagementAction, error)
}

