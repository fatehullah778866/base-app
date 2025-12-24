package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
)

type DashboardService struct {
	dashboardRepo repositories.DashboardRepository
	logger         *zap.Logger
}

func NewDashboardService(
	dashboardRepo repositories.DashboardRepository,
	logger *zap.Logger,
) *DashboardService {
	return &DashboardService{
		dashboardRepo: dashboardRepo,
		logger:        logger,
	}
}

// CreateItem creates a new dashboard item
func (s *DashboardService) CreateItem(ctx context.Context, userID uuid.UUID, title string, description *string, category *string, metadata *string) (*models.DashboardItem, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	item := &models.DashboardItem{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       title,
		Description: description,
		Category:    category,
		Status:      "active",
		Priority:    0,
		Metadata:    metadata,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.dashboardRepo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

// GetItem retrieves a dashboard item by ID
func (s *DashboardService) GetItem(ctx context.Context, id uuid.UUID) (*models.DashboardItem, error) {
	return s.dashboardRepo.GetByID(ctx, id)
}

// ListItems retrieves all dashboard items for a user
func (s *DashboardService) ListItems(ctx context.Context, userID uuid.UUID, status string) ([]*models.DashboardItem, error) {
	return s.dashboardRepo.GetByUserID(ctx, userID, status)
}

// UpdateItem updates a dashboard item
func (s *DashboardService) UpdateItem(ctx context.Context, id uuid.UUID, userID uuid.UUID, updates map[string]interface{}) (*models.DashboardItem, error) {
	item, err := s.dashboardRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if item.UserID != userID {
		return nil, errors.New("unauthorized: item does not belong to user")
	}

	// Update allowed fields
	if title, ok := updates["title"].(string); ok && title != "" {
		item.Title = title
	}
	if description, ok := updates["description"].(string); ok {
		item.Description = &description
	}
	if category, ok := updates["category"].(string); ok {
		item.Category = &category
	}
	if status, ok := updates["status"].(string); ok {
		item.Status = status
	}
	if priority, ok := updates["priority"].(float64); ok {
		item.Priority = int(priority)
	}
	if metadata, ok := updates["metadata"].(string); ok {
		item.Metadata = &metadata
	}

	item.UpdatedAt = time.Now()

	if err := s.dashboardRepo.Update(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

// DeleteItem permanently deletes a dashboard item
func (s *DashboardService) DeleteItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	item, err := s.dashboardRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify ownership
	if item.UserID != userID {
		return errors.New("unauthorized: item does not belong to user")
	}

	return s.dashboardRepo.Delete(ctx, id)
}

// SoftDeleteItem soft deletes a dashboard item
func (s *DashboardService) SoftDeleteItem(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	item, err := s.dashboardRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify ownership
	if item.UserID != userID {
		return errors.New("unauthorized: item does not belong to user")
	}

	return s.dashboardRepo.SoftDelete(ctx, id)
}

