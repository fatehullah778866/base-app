package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
)

type notificationRepository struct {
	db *database.DB
}

func NewNotificationRepository(db *database.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	query := `
		INSERT INTO notifications (id, user_id, type, title, message, link, is_read, metadata, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		notification.ID.String(),
		notification.UserID.String(),
		notification.Type,
		notification.Title,
		notification.Message,
		notification.Link,
		notification.IsRead,
		notification.Metadata,
		notification.CreatedAt,
	)
	return err
}

func (r *notificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Notification, error) {
	var n models.Notification
	var userIDStr, idStr string
	var link, metadata sql.NullString
	var readAt sql.NullTime

	query := `SELECT id, user_id, type, title, message, link, is_read, read_at, metadata, created_at
		FROM notifications WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &userIDStr, &n.Type, &n.Title, &n.Message,
		&link, &n.IsRead, &readAt, &metadata, &n.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	n.ID, _ = uuid.Parse(idStr)
	n.UserID, _ = uuid.Parse(userIDStr)
	if link.Valid {
		n.Link = &link.String
	}
	if readAt.Valid {
		n.ReadAt = &readAt.Time
	}
	if metadata.Valid {
		n.Metadata = &metadata.String
	}

	return &n, nil
}

func (r *notificationRepository) GetByUserID(ctx context.Context, userID uuid.UUID, unreadOnly bool, limit int) ([]*models.Notification, error) {
	var query string
	if unreadOnly {
		query = `SELECT id, user_id, type, title, message, link, is_read, read_at, metadata, created_at
			FROM notifications WHERE user_id = ? AND is_read = 0
			ORDER BY created_at DESC LIMIT ?`
	} else {
		query = `SELECT id, user_id, type, title, message, link, is_read, read_at, metadata, created_at
			FROM notifications WHERE user_id = ?
			ORDER BY created_at DESC LIMIT ?`
	}

	rows, err := r.db.QueryContext(ctx, query, userID.String(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		var n models.Notification
		var userIDStr, idStr string
		var link, metadata sql.NullString
		var readAt sql.NullTime

		err := rows.Scan(&idStr, &userIDStr, &n.Type, &n.Title, &n.Message,
			&link, &n.IsRead, &readAt, &metadata, &n.CreatedAt)
		if err != nil {
			return nil, err
		}

		n.ID, _ = uuid.Parse(idStr)
		n.UserID, _ = uuid.Parse(userIDStr)
		if link.Valid {
			n.Link = &link.String
		}
		if readAt.Valid {
			n.ReadAt = &readAt.Time
		}
		if metadata.Valid {
			n.Metadata = &metadata.String
		}

		notifications = append(notifications, &n)
	}

	return notifications, rows.Err()
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE notifications SET is_read = 1, read_at = datetime('now') WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE notifications SET is_read = 1, read_at = datetime('now') WHERE user_id = ? AND is_read = 0`
	_, err := r.db.ExecContext(ctx, query, userID.String())
	return err
}

func (r *notificationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM notifications WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

func (r *notificationRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = 0`
	err := r.db.QueryRowContext(ctx, query, userID.String()).Scan(&count)
	return count, err
}

