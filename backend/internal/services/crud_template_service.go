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

type CRUDTemplateService struct {
	templateRepo repositories.CRUDTemplateRepository
	logger       *zap.Logger
}

func NewCRUDTemplateService(templateRepo repositories.CRUDTemplateRepository, logger *zap.Logger) *CRUDTemplateService {
	return &CRUDTemplateService{
		templateRepo: templateRepo,
		logger:       logger,
	}
}

// CreateTemplate creates a new CRUD template
func (s *CRUDTemplateService) CreateTemplate(ctx context.Context, createdBy uuid.UUID, name, displayName string, description *string, schema map[string]interface{}, icon, category *string) (*models.CRUDTemplate, error) {
	// Validate schema is valid JSON
	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		return nil, errors.New("invalid schema format")
	}

	// Check if template name already exists
	existing, _ := s.templateRepo.GetByName(ctx, name)
	if existing != nil {
		return nil, errors.New("template name already exists")
	}

	template := &models.CRUDTemplate{
		ID:          uuid.New(),
		Name:        name,
		DisplayName: displayName,
		Description: description,
		Schema:      string(schemaJSON),
		Icon:        icon,
		Category:    category,
		CreatedBy:   createdBy,
		IsActive:    true,
		IsSystem:    false, // User-created templates are not system templates
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.templateRepo.Create(ctx, template); err != nil {
		return nil, err
	}

	return template, nil
}

// GetTemplate retrieves a template by ID
func (s *CRUDTemplateService) GetTemplate(ctx context.Context, id uuid.UUID) (*models.CRUDTemplate, error) {
	return s.templateRepo.GetByID(ctx, id)
}

// GetTemplateByName retrieves a template by name
func (s *CRUDTemplateService) GetTemplateByName(ctx context.Context, name string) (*models.CRUDTemplate, error) {
	return s.templateRepo.GetByName(ctx, name)
}

// ListTemplates retrieves all templates, optionally filtered
func (s *CRUDTemplateService) ListTemplates(ctx context.Context, category *string, activeOnly bool) ([]*models.CRUDTemplate, error) {
	return s.templateRepo.List(ctx, category, activeOnly)
}

// ListTemplatesByCreator retrieves templates created by a specific admin
func (s *CRUDTemplateService) ListTemplatesByCreator(ctx context.Context, createdBy uuid.UUID) ([]*models.CRUDTemplate, error) {
	return s.templateRepo.ListByCreator(ctx, createdBy)
}

// UpdateTemplate updates an existing template
func (s *CRUDTemplateService) UpdateTemplate(ctx context.Context, id uuid.UUID, displayName *string, description *string, schema map[string]interface{}, icon, category *string) (*models.CRUDTemplate, error) {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if template == nil {
		return nil, errors.New("template not found")
	}

	if displayName != nil {
		template.DisplayName = *displayName
	}
	if description != nil {
		template.Description = description
	}
	if schema != nil {
		schemaJSON, err := json.Marshal(schema)
		if err != nil {
			return nil, errors.New("invalid schema format")
		}
		template.Schema = string(schemaJSON)
	}
	if icon != nil {
		template.Icon = icon
	}
	if category != nil {
		template.Category = category
	}
	template.UpdatedAt = time.Now()

	if err := s.templateRepo.Update(ctx, template); err != nil {
		return nil, err
	}

	return template, nil
}

// DeleteTemplate deletes a template (only if not system template)
func (s *CRUDTemplateService) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	template, err := s.templateRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if template == nil {
		return errors.New("template not found")
	}
	if template.IsSystem {
		return errors.New("cannot delete system template")
	}

	return s.templateRepo.Delete(ctx, id)
}

// ActivateTemplate activates a template
func (s *CRUDTemplateService) ActivateTemplate(ctx context.Context, id uuid.UUID) error {
	return s.templateRepo.Activate(ctx, id)
}

// DeactivateTemplate deactivates a template
func (s *CRUDTemplateService) DeactivateTemplate(ctx context.Context, id uuid.UUID) error {
	return s.templateRepo.Deactivate(ctx, id)
}

// GetTemplateSchemaAsMap converts template schema JSON to map
func (s *CRUDTemplateService) GetTemplateSchemaAsMap(template *models.CRUDTemplate) (map[string]interface{}, error) {
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(template.Schema), &schema); err != nil {
		return nil, err
	}
	return schema, nil
}


