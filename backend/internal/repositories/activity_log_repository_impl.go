package repositories

import (
	"context"
	"database/sql"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
)

type activityLogRepository struct {
	db *database.DB
}

func NewActivityLogRepository(db *database.DB) ActivityLogRepository {
	return &activityLogRepository{db: db}
}

func (r *activityLogRepository) Create(ctx context.Context, log *models.ActivityLog) error {
	query := `
		INSERT INTO activity_logs (id, actor_id, actor_role, action, target_type, target_id, metadata, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		log.ID, log.ActorID, log.ActorRole, log.Action, log.TargetType, log.TargetID, log.Metadata, log.CreatedAt,
	)
	return err
}

func (r *activityLogRepository) List(ctx context.Context, limit int) ([]*models.ActivityLog, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}

	rows, err := r.db.DB.QueryContext(ctx, `
		SELECT id, actor_id, actor_role, action, target_type, target_id, metadata, created_at
		FROM activity_logs
		ORDER BY created_at DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.ActivityLog
	for rows.Next() {
		entry := &models.ActivityLog{}
		var actorID sql.NullString
		var actorRole sql.NullString
		var targetType sql.NullString
		var targetID sql.NullString
		var metadata sql.NullString

		if err := rows.Scan(
			&entry.ID, &actorID, &actorRole, &entry.Action, &targetType, &targetID, &metadata, &entry.CreatedAt,
		); err != nil {
			return nil, err
		}

		if actorID.Valid {
			entry.ActorID = &actorID.String
		}
		if actorRole.Valid {
			entry.ActorRole = &actorRole.String
		}
		if targetType.Valid {
			entry.TargetType = &targetType.String
		}
		if targetID.Valid {
			entry.TargetID = &targetID.String
		}
		if metadata.Valid {
			entry.Metadata = &metadata.String
		}

		logs = append(logs, entry)
	}

	return logs, nil
}
