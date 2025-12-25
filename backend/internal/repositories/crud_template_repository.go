package repositories

import (
	"context"

	"github.com/google/uuid"

	"base-app-service/internal/models"
)

// CRUDTemplateRepository defines the interface for CRUD template operations
type CRUDTemplateRepository interface {
	// Create creates a new CRUD template
	Create(ctx context.Context, template *models.CRUDTemplate) error

	// GetByID retrieves a template by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.CRUDTemplate, error)

	// GetByName retrieves a template by name
	GetByName(ctx context.Context, name string) (*models.CRUDTemplate, error)

	// List retrieves all active templates, optionally filtered by category
	List(ctx context.Context, category *string, activeOnly bool) ([]*models.CRUDTemplate, error)

	// ListByCreator retrieves templates created by a specific admin
	ListByCreator(ctx context.Context, createdBy uuid.UUID) ([]*models.CRUDTemplate, error)

	// Update updates an existing template
	Update(ctx context.Context, template *models.CRUDTemplate) error

	// Delete deletes a template (only if not system template)
	Delete(ctx context.Context, id uuid.UUID) error

	// Activate activates a template
	Activate(ctx context.Context, id uuid.UUID) error

	// Deactivate deactivates a template
	Deactivate(ctx context.Context, id uuid.UUID) error
}


