package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
)

type crudTemplateRepository struct {
	db *database.DB
}

func NewCRUDTemplateRepository(db *database.DB) CRUDTemplateRepository {
	return &crudTemplateRepository{db: db}
}

func (r *crudTemplateRepository) Create(ctx context.Context, template *models.CRUDTemplate) error {
	query := `INSERT INTO crud_templates (id, name, display_name, description, schema, icon, category, created_by, is_active, is_system, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	isActive := 0
	if template.IsActive {
		isActive = 1
	}
	isSystem := 0
	if template.IsSystem {
		isSystem = 1
	}

	_, err := r.db.ExecContext(ctx, query,
		template.ID.String(),
		template.Name,
		template.DisplayName,
		template.Description,
		template.Schema,
		template.Icon,
		template.Category,
		template.CreatedBy.String(),
		isActive,
		isSystem,
		time.Now(),
		time.Now(),
	)
	return err
}

func (r *crudTemplateRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.CRUDTemplate, error) {
	var t models.CRUDTemplate
	var description, icon, category sql.NullString
	var isActive, isSystem int

	query := `SELECT id, name, display_name, description, schema, icon, category, created_by, is_active, is_system, created_at, updated_at
		FROM crud_templates WHERE id = ?`
	
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&t.ID, &t.Name, &t.DisplayName, &description, &t.Schema, &icon, &category,
		&t.CreatedBy, &isActive, &isSystem, &t.CreatedAt, &t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if description.Valid {
		t.Description = &description.String
	}
	if icon.Valid {
		t.Icon = &icon.String
	}
	if category.Valid {
		t.Category = &category.String
	}
	t.IsActive = isActive == 1
	t.IsSystem = isSystem == 1

	return &t, nil
}

func (r *crudTemplateRepository) GetByName(ctx context.Context, name string) (*models.CRUDTemplate, error) {
	var t models.CRUDTemplate
	var description, icon, category sql.NullString
	var isActive, isSystem int

	query := `SELECT id, name, display_name, description, schema, icon, category, created_by, is_active, is_system, created_at, updated_at
		FROM crud_templates WHERE name = ?`
	
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&t.ID, &t.Name, &t.DisplayName, &description, &t.Schema, &icon, &category,
		&t.CreatedBy, &isActive, &isSystem, &t.CreatedAt, &t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if description.Valid {
		t.Description = &description.String
	}
	if icon.Valid {
		t.Icon = &icon.String
	}
	if category.Valid {
		t.Category = &category.String
	}
	t.IsActive = isActive == 1
	t.IsSystem = isSystem == 1

	return &t, nil
}

func (r *crudTemplateRepository) List(ctx context.Context, category *string, activeOnly bool) ([]*models.CRUDTemplate, error) {
	query := `SELECT id, name, display_name, description, schema, icon, category, created_by, is_active, is_system, created_at, updated_at
		FROM crud_templates WHERE 1=1`
	args := []interface{}{}

	if category != nil {
		query += ` AND category = ?`
		args = append(args, *category)
	}

	if activeOnly {
		query += ` AND is_active = 1`
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*models.CRUDTemplate
	for rows.Next() {
		var t models.CRUDTemplate
		var description, icon, category sql.NullString
		var isActive, isSystem int

		err := rows.Scan(
			&t.ID, &t.Name, &t.DisplayName, &description, &t.Schema, &icon, &category,
			&t.CreatedBy, &isActive, &isSystem, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			t.Description = &description.String
		}
		if icon.Valid {
			t.Icon = &icon.String
		}
		if category.Valid {
			t.Category = &category.String
		}
		t.IsActive = isActive == 1
		t.IsSystem = isSystem == 1

		templates = append(templates, &t)
	}

	return templates, rows.Err()
}

func (r *crudTemplateRepository) ListByCreator(ctx context.Context, createdBy uuid.UUID) ([]*models.CRUDTemplate, error) {
	query := `SELECT id, name, display_name, description, schema, icon, category, created_by, is_active, is_system, created_at, updated_at
		FROM crud_templates WHERE created_by = ? ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, createdBy.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*models.CRUDTemplate
	for rows.Next() {
		var t models.CRUDTemplate
		var description, icon, category sql.NullString
		var isActive, isSystem int

		err := rows.Scan(
			&t.ID, &t.Name, &t.DisplayName, &description, &t.Schema, &icon, &category,
			&t.CreatedBy, &isActive, &isSystem, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if description.Valid {
			t.Description = &description.String
		}
		if icon.Valid {
			t.Icon = &icon.String
		}
		if category.Valid {
			t.Category = &category.String
		}
		t.IsActive = isActive == 1
		t.IsSystem = isSystem == 1

		templates = append(templates, &t)
	}

	return templates, rows.Err()
}

func (r *crudTemplateRepository) Update(ctx context.Context, template *models.CRUDTemplate) error {
	isActive := 0
	if template.IsActive {
		isActive = 1
	}

	query := `UPDATE crud_templates SET display_name = ?, description = ?, schema = ?, icon = ?, category = ?, is_active = ?, updated_at = ?
		WHERE id = ?`
	
	_, err := r.db.ExecContext(ctx, query,
		template.DisplayName,
		template.Description,
		template.Schema,
		template.Icon,
		template.Category,
		isActive,
		time.Now(),
		template.ID.String(),
	)
	return err
}

func (r *crudTemplateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Only delete if not a system template
	query := `DELETE FROM crud_templates WHERE id = ? AND is_system = 0`
	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *crudTemplateRepository) Activate(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE crud_templates SET is_active = 1, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id.String())
	return err
}

func (r *crudTemplateRepository) Deactivate(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE crud_templates SET is_active = 0, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id.String())
	return err
}

