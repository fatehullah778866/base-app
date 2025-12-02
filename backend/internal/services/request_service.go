package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
)

type RequestService struct {
	repo   repositories.AccessRequestRepository
	logger *zap.Logger
}

type CreateRequestInput struct {
	Title   *string
	Details *string
}

func NewRequestService(repo repositories.AccessRequestRepository, logger *zap.Logger) *RequestService {
	return &RequestService{
		repo:   repo,
		logger: logger,
	}
}

func (s *RequestService) Create(ctx context.Context, userID uuid.UUID, input CreateRequestInput) (*models.AccessRequest, error) {
	now := time.Now()
	req := &models.AccessRequest{
		ID:        uuid.New().String(),
		UserID:    userID.String(),
		Title:     input.Title,
		Details:   input.Details,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *RequestService) ListForUser(ctx context.Context, userID uuid.UUID) ([]*models.AccessRequest, error) {
	return s.repo.ListByUser(ctx, userID.String())
}

func (s *RequestService) Get(ctx context.Context, id string) (*models.AccessRequest, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *RequestService) UpdateStatus(ctx context.Context, id string, status string, feedback *string) (*models.AccessRequest, error) {
	status = strings.ToLower(status)
	if status != "approved" && status != "rejected" && status != "pending" {
		return nil, errors.New("invalid status")
	}
	return s.repo.UpdateStatus(ctx, id, status, feedback)
}
