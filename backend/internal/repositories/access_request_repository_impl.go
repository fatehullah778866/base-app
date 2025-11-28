package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
)

type accessRequestRepository struct {
	db *database.DB
}

func NewAccessRequestRepository(db *database.DB) AccessRequestRepository {
	return &accessRequestRepository{db: db}
}

func (r *accessRequestRepository) Create(ctx context.Context, request *models.AccessRequest) error {
	query := `
		INSERT INTO access_requests (id, user_id, title, details, status, feedback, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		request.ID, request.UserID, request.Title, request.Details, request.Status, request.Feedback, request.CreatedAt, request.UpdatedAt,
	)
	return err
}

func (r *accessRequestRepository) List(ctx context.Context, status *string) ([]*models.AccessRequest, error) {
	var rows *sql.Rows
	var err error
	if status != nil {
		rows, err = r.db.DB.QueryContext(ctx, `
			SELECT id, user_id, title, details, status, feedback, created_at, updated_at
			FROM access_requests
			WHERE status = ?
			ORDER BY created_at DESC
			LIMIT 200
		`, *status)
	} else {
		rows, err = r.db.DB.QueryContext(ctx, `
			SELECT id, user_id, title, details, status, feedback, created_at, updated_at
			FROM access_requests
			ORDER BY created_at DESC
			LIMIT 200
		`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*models.AccessRequest
	for rows.Next() {
		req := &models.AccessRequest{}
		if err := rows.Scan(&req.ID, &req.UserID, &req.Title, &req.Details, &req.Status, &req.Feedback, &req.CreatedAt, &req.UpdatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (r *accessRequestRepository) ListByUser(ctx context.Context, userID string) ([]*models.AccessRequest, error) {
	rows, err := r.db.DB.QueryContext(ctx, `
		SELECT id, user_id, title, details, status, feedback, created_at, updated_at
		FROM access_requests
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT 200
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*models.AccessRequest
	for rows.Next() {
		req := &models.AccessRequest{}
		if err := rows.Scan(&req.ID, &req.UserID, &req.Title, &req.Details, &req.Status, &req.Feedback, &req.CreatedAt, &req.UpdatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	return requests, nil
}

func (r *accessRequestRepository) UpdateStatus(ctx context.Context, id string, status string, feedback *string) (*models.AccessRequest, error) {
	_, err := r.db.DB.ExecContext(ctx, `
		UPDATE access_requests
		SET status = ?, feedback = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, status, feedback, id)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *accessRequestRepository) GetByID(ctx context.Context, id string) (*models.AccessRequest, error) {
	req := &models.AccessRequest{}
	err := r.db.DB.QueryRowContext(ctx, `
		SELECT id, user_id, title, details, status, feedback, created_at, updated_at
		FROM access_requests
		WHERE id = ?
	`, id).Scan(&req.ID, &req.UserID, &req.Title, &req.Details, &req.Status, &req.Feedback, &req.CreatedAt, &req.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("access request not found")
	}
	if err != nil {
		return nil, err
	}
	return req, nil
}
