package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
)

type ActivityLogService struct {
	logRepo repositories.ActivityLogRepository
	logger  *zap.Logger
}

func NewActivityLogService(logRepo repositories.ActivityLogRepository, logger *zap.Logger) *ActivityLogService {
	return &ActivityLogService{
		logRepo: logRepo,
		logger:  logger,
	}
}

func (s *ActivityLogService) Record(ctx context.Context, actorID *uuid.UUID, actorRole string, action string, targetType *string, targetID *string, metadata map[string]interface{}) {
	var metaStr *string
	if metadata != nil {
		if b, err := json.Marshal(metadata); err == nil {
			jsonStr := string(b)
			metaStr = &jsonStr
		}
	}

	var actorIDStr *string
	if actorID != nil {
		value := actorID.String()
		actorIDStr = &value
	}
	roleValue := actorRole
	log := &models.ActivityLog{
		ID:         uuid.New().String(),
		ActorID:    actorIDStr,
		ActorRole:  &roleValue,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Metadata:   metaStr,
		CreatedAt:  time.Now(),
	}

	if err := s.logRepo.Create(ctx, log); err != nil {
		s.logger.Warn("failed to record activity log", zap.Error(err))
	}
}

func (s *ActivityLogService) List(ctx context.Context, limit int) ([]*models.ActivityLog, error) {
	return s.logRepo.List(ctx, limit)
}
