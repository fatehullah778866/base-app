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

type CustomCRUDService struct {
	crudRepo repositories.CustomCRUDRepository
	logger   *zap.Logger
}

func NewCustomCRUDService(crudRepo repositories.CustomCRUDRepository, logger *zap.Logger) *CustomCRUDService {
	return &CustomCRUDService{
		crudRepo: crudRepo,
		logger:   logger,
	}
}

func (s *CustomCRUDService) CreateEntity(ctx context.Context, createdBy uuid.UUID, entityName, displayName string, description *string, schema map[string]interface{}) (*models.CustomCRUDEntity, error) {
	// Validate schema is valid JSON
	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		return nil, errors.New("invalid schema format")
	}

	// Check if entity name already exists for this user (allow same name for different users)
	existing, _ := s.crudRepo.GetEntityByName(ctx, entityName)
	if existing != nil && existing.CreatedBy == createdBy {
		return nil, errors.New("you already have an entity with this name")
	}

	entity := &models.CustomCRUDEntity{
		ID:          uuid.New(),
		CreatedBy:   createdBy,
		EntityName:  entityName,
		DisplayName: displayName,
		Description: description,
		Schema:      string(schemaJSON),
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.crudRepo.CreateEntity(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *CustomCRUDService) GetEntity(ctx context.Context, id uuid.UUID) (*models.CustomCRUDEntity, error) {
	return s.crudRepo.GetEntityByID(ctx, id)
}

func (s *CustomCRUDService) ListEntities(ctx context.Context, createdBy *uuid.UUID, activeOnly bool) ([]*models.CustomCRUDEntity, error) {
	return s.crudRepo.ListEntities(ctx, createdBy, activeOnly)
}

func (s *CustomCRUDService) UpdateEntity(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.CustomCRUDEntity, error) {
	entity, err := s.crudRepo.GetEntityByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if displayName, ok := updates["display_name"].(string); ok {
		entity.DisplayName = displayName
	}
	if description, ok := updates["description"].(string); ok {
		entity.Description = &description
	}
	if schema, ok := updates["schema"].(map[string]interface{}); ok {
		schemaJSON, err := json.Marshal(schema)
		if err == nil {
			entity.Schema = string(schemaJSON)
		}
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		entity.IsActive = isActive
	}

	entity.UpdatedAt = time.Now()
	if err := s.crudRepo.UpdateEntity(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *CustomCRUDService) DeleteEntity(ctx context.Context, id uuid.UUID) error {
	return s.crudRepo.DeleteEntity(ctx, id)
}

func (s *CustomCRUDService) CreateData(ctx context.Context, entityID uuid.UUID, createdBy uuid.UUID, data map[string]interface{}) (*models.CustomCRUDData, error) {
	// Validate entity exists
	entity, err := s.crudRepo.GetEntityByID(ctx, entityID)
	if err != nil || entity == nil {
		return nil, errors.New("entity not found")
	}

	// Parse schema and validate
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(entity.Schema), &schema); err != nil {
		s.logger.Warn("Failed to parse entity schema", zap.Error(err))
	} else {
		// Validate data against schema
		if err := ValidateDataAgainstSchema(data, schema); err != nil {
			return nil, errors.New("data validation failed: " + err.Error())
		}
	}

	// Validate data against schema (simplified - in production, use JSON schema validator)
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("invalid data format")
	}

	crudData := &models.CustomCRUDData{
		ID:        uuid.New(),
		EntityID:  entityID,
		Data:      string(dataJSON),
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.crudRepo.CreateData(ctx, crudData); err != nil {
		return nil, err
	}

	return crudData, nil
}

func (s *CustomCRUDService) GetData(ctx context.Context, id uuid.UUID) (*models.CustomCRUDData, error) {
	return s.crudRepo.GetDataByID(ctx, id)
}

func (s *CustomCRUDService) ListData(ctx context.Context, entityID uuid.UUID, limit, offset int) ([]*models.CustomCRUDData, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.crudRepo.ListDataByEntity(ctx, entityID, limit, offset)
}

func (s *CustomCRUDService) UpdateData(ctx context.Context, id uuid.UUID, updatedBy uuid.UUID, data map[string]interface{}) (*models.CustomCRUDData, error) {
	crudData, err := s.crudRepo.GetDataByID(ctx, id)
	if err != nil {
		return nil, err
	}

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("invalid data format")
	}

	crudData.Data = string(dataJSON)
	crudData.UpdatedBy = &updatedBy
	crudData.UpdatedAt = time.Now()

	if err := s.crudRepo.UpdateData(ctx, crudData); err != nil {
		return nil, err
	}

	return crudData, nil
}

func (s *CustomCRUDService) DeleteData(ctx context.Context, id uuid.UUID) error {
	return s.crudRepo.DeleteData(ctx, id)
}

