package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
)

type dashboardRepository struct {
	db *database.DB
}

func NewDashboardRepository(db *database.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) Create(ctx context.Context, item *models.DashboardItem) error {
	query := `
		INSERT INTO dashboard_items (id, user_id, title, description, category, status, priority, metadata, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		item.ID.String(),
		item.UserID.String(),
		item.Title,
		item.Description,
		item.Category,
		item.Status,
		item.Priority,
		item.Metadata,
		now,
		now,
	)
	return err
}

func (r *dashboardRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.DashboardItem, error) {
	var item models.DashboardItem
	var userIDStr, idStr string
	var description, category, metadata sql.NullString
	var deletedAt sql.NullTime

	query := `
		SELECT id, user_id, title, description, category, status, priority, metadata, created_at, updated_at, deleted_at
		FROM dashboard_items
		WHERE id = ? AND deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr,
		&userIDStr,
		&item.Title,
		&description,
		&category,
		&item.Status,
		&item.Priority,
		&metadata,
		&item.CreatedAt,
		&item.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, err
	}

	item.ID, _ = uuid.Parse(idStr)
	item.UserID, _ = uuid.Parse(userIDStr)
	if description.Valid {
		item.Description = &description.String
	}
	if category.Valid {
		item.Category = &category.String
	}
	if metadata.Valid {
		item.Metadata = &metadata.String
	}
	if deletedAt.Valid {
		item.DeletedAt = &deletedAt.Time
	}

	return &item, nil
}

func (r *dashboardRepository) GetByUserID(ctx context.Context, userID uuid.UUID, status string) ([]*models.DashboardItem, error) {
	var items []*models.DashboardItem
	var query string
	var args []interface{}

	if status == "" {
		query = `
			SELECT id, user_id, title, description, category, status, priority, metadata, created_at, updated_at, deleted_at
			FROM dashboard_items
			WHERE user_id = ? AND deleted_at IS NULL
			ORDER BY priority DESC, created_at DESC
		`
		args = []interface{}{userID.String()}
	} else {
		query = `
			SELECT id, user_id, title, description, category, status, priority, metadata, created_at, updated_at, deleted_at
			FROM dashboard_items
			WHERE user_id = ? AND status = ? AND deleted_at IS NULL
			ORDER BY priority DESC, created_at DESC
		`
		args = []interface{}{userID.String(), status}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.DashboardItem
		var userIDStr, idStr string
		var description, category, metadata sql.NullString
		var deletedAt sql.NullTime

		err := rows.Scan(
			&idStr,
			&userIDStr,
			&item.Title,
			&description,
			&category,
			&item.Status,
			&item.Priority,
			&metadata,
			&item.CreatedAt,
			&item.UpdatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, err
		}

		item.ID, _ = uuid.Parse(idStr)
		item.UserID, _ = uuid.Parse(userIDStr)
		if description.Valid {
			item.Description = &description.String
		}
		if category.Valid {
			item.Category = &category.String
		}
		if metadata.Valid {
			item.Metadata = &metadata.String
		}
		if deletedAt.Valid {
			item.DeletedAt = &deletedAt.Time
		}

		items = append(items, &item)
	}

	return items, rows.Err()
}

func (r *dashboardRepository) Update(ctx context.Context, item *models.DashboardItem) error {
	query := `
		UPDATE dashboard_items SET
			title = ?, description = ?, category = ?, status = ?, priority = ?, metadata = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query,
		item.Title,
		item.Description,
		item.Category,
		item.Status,
		item.Priority,
		item.Metadata,
		time.Now(),
		item.ID.String(),
	)
	return err
}

func (r *dashboardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM dashboard_items WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

func (r *dashboardRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE dashboard_items SET deleted_at = datetime('now'), status = 'deleted' WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

