package repositories

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
)

type adminSettingsRepository struct {
	db *database.DB
}

func NewAdminSettingsRepository(db *database.DB) AdminSettingsRepository {
	return &adminSettingsRepository{db: db}
}

func (r *adminSettingsRepository) GetByAdminID(ctx context.Context, adminID uuid.UUID) (*models.AdminSettings, error) {
	var s models.AdminSettings
	var dashboardLayout, defaultPermissions, notificationPreferences, themePreferences, adminVerificationCode sql.NullString

	query := `SELECT admin_id, dashboard_layout, default_permissions, notification_preferences, theme_preferences, admin_verification_code, created_at, updated_at
		FROM admin_settings WHERE admin_id = ?`
	err := r.db.QueryRowContext(ctx, query, adminID.String()).Scan(
		&s.AdminID, &dashboardLayout, &defaultPermissions, &notificationPreferences, &themePreferences, &adminVerificationCode,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if dashboardLayout.Valid {
		s.DashboardLayout = &dashboardLayout.String
	}
	if defaultPermissions.Valid {
		s.DefaultPermissions = &defaultPermissions.String
	}
	if notificationPreferences.Valid {
		s.NotificationPreferences = &notificationPreferences.String
	}
	if themePreferences.Valid {
		s.ThemePreferences = &themePreferences.String
	}
	if adminVerificationCode.Valid {
		s.AdminVerificationCode = &adminVerificationCode.String
	}

	return &s, nil
}

func (r *adminSettingsRepository) Create(ctx context.Context, settings *models.AdminSettings) error {
	defaultCode := "Kompasstech2025@"
	if settings.AdminVerificationCode == nil {
		settings.AdminVerificationCode = &defaultCode
	}
	query := `INSERT INTO admin_settings (admin_id, dashboard_layout, default_permissions, notification_preferences, theme_preferences, admin_verification_code, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		settings.AdminID.String(),
		settings.DashboardLayout,
		settings.DefaultPermissions,
		settings.NotificationPreferences,
		settings.ThemePreferences,
		settings.AdminVerificationCode,
		time.Now(),
		time.Now(),
	)
	return err
}

func (r *adminSettingsRepository) Update(ctx context.Context, settings *models.AdminSettings) error {
	query := `UPDATE admin_settings SET dashboard_layout = ?, default_permissions = ?, notification_preferences = ?, theme_preferences = ?, admin_verification_code = ?, updated_at = ?
		WHERE admin_id = ?`
	_, err := r.db.ExecContext(ctx, query,
		settings.DashboardLayout,
		settings.DefaultPermissions,
		settings.NotificationPreferences,
		settings.ThemePreferences,
		settings.AdminVerificationCode,
		time.Now(),
		settings.AdminID.String(),
	)
	return err
}

func (r *adminSettingsRepository) GetFirstAdminVerificationCode(ctx context.Context) (string, error) {
	var adminVerificationCode sql.NullString
	
	// Get the verification code from the first admin's settings
	query := `SELECT admin_verification_code FROM admin_settings 
		WHERE admin_verification_code IS NOT NULL AND admin_verification_code != ''
		ORDER BY created_at ASC LIMIT 1`
	
	err := r.db.QueryRowContext(ctx, query).Scan(&adminVerificationCode)
	if err == sql.ErrNoRows {
		// No admin settings found, return default
		return "Kompasstech2025@", nil
	}
	if err != nil {
		return "Kompasstech2025@", err
	}
	
	if adminVerificationCode.Valid && adminVerificationCode.String != "" {
		// Trim the code before returning
		return strings.TrimSpace(adminVerificationCode.String), nil
	}
	
	// Fallback to default
	return "Kompasstech2025@", nil
}

type customCRUDRepository struct {
	db *database.DB
}

func NewCustomCRUDRepository(db *database.DB) CustomCRUDRepository {
	return &customCRUDRepository{db: db}
}

func (r *customCRUDRepository) CreateEntity(ctx context.Context, entity *models.CustomCRUDEntity) error {
	query := `INSERT INTO custom_crud_entities (id, created_by, entity_name, display_name, description, schema, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		entity.ID.String(),
		entity.CreatedBy.String(),
		entity.EntityName,
		entity.DisplayName,
		entity.Description,
		entity.Schema,
		entity.IsActive,
		entity.CreatedAt,
		entity.UpdatedAt,
	)
	return err
}

func (r *customCRUDRepository) GetEntityByID(ctx context.Context, id uuid.UUID) (*models.CustomCRUDEntity, error) {
	var e models.CustomCRUDEntity
	var createdByStr, idStr string
	var description sql.NullString

	query := `SELECT id, created_by, entity_name, display_name, description, schema, is_active, created_at, updated_at
		FROM custom_crud_entities WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &createdByStr, &e.EntityName, &e.DisplayName, &description,
		&e.Schema, &e.IsActive, &e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	e.ID, _ = uuid.Parse(idStr)
	e.CreatedBy, _ = uuid.Parse(createdByStr)
	if description.Valid {
		e.Description = &description.String
	}

	return &e, nil
}

func (r *customCRUDRepository) GetEntityByName(ctx context.Context, name string) (*models.CustomCRUDEntity, error) {
	var e models.CustomCRUDEntity
	var createdByStr, idStr string
	var description sql.NullString

	query := `SELECT id, created_by, entity_name, display_name, description, schema, is_active, created_at, updated_at
		FROM custom_crud_entities WHERE entity_name = ?`
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&idStr, &createdByStr, &e.EntityName, &e.DisplayName, &description,
		&e.Schema, &e.IsActive, &e.CreatedAt, &e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	e.ID, _ = uuid.Parse(idStr)
	e.CreatedBy, _ = uuid.Parse(createdByStr)
	if description.Valid {
		e.Description = &description.String
	}

	return &e, nil
}

func (r *customCRUDRepository) ListEntities(ctx context.Context, createdBy *uuid.UUID, activeOnly bool) ([]*models.CustomCRUDEntity, error) {
	var query string
	var args []interface{}

	if createdBy != nil && activeOnly {
		query = `SELECT id, created_by, entity_name, display_name, description, schema, is_active, created_at, updated_at
			FROM custom_crud_entities WHERE created_by = ? AND is_active = 1 ORDER BY created_at DESC`
		args = []interface{}{createdBy.String()}
	} else if createdBy != nil {
		query = `SELECT id, created_by, entity_name, display_name, description, schema, is_active, created_at, updated_at
			FROM custom_crud_entities WHERE created_by = ? ORDER BY created_at DESC`
		args = []interface{}{createdBy.String()}
	} else if activeOnly {
		query = `SELECT id, created_by, entity_name, display_name, description, schema, is_active, created_at, updated_at
			FROM custom_crud_entities WHERE is_active = 1 ORDER BY created_at DESC`
		args = []interface{}{}
	} else {
		query = `SELECT id, created_by, entity_name, display_name, description, schema, is_active, created_at, updated_at
			FROM custom_crud_entities ORDER BY created_at DESC`
		args = []interface{}{}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*models.CustomCRUDEntity
	for rows.Next() {
		var e models.CustomCRUDEntity
		var createdByStr, idStr string
		var description sql.NullString

		err := rows.Scan(&idStr, &createdByStr, &e.EntityName, &e.DisplayName, &description,
			&e.Schema, &e.IsActive, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}

		e.ID, _ = uuid.Parse(idStr)
		e.CreatedBy, _ = uuid.Parse(createdByStr)
		if description.Valid {
			e.Description = &description.String
		}

		entities = append(entities, &e)
	}

	return entities, rows.Err()
}

func (r *customCRUDRepository) UpdateEntity(ctx context.Context, entity *models.CustomCRUDEntity) error {
	query := `UPDATE custom_crud_entities SET display_name = ?, description = ?, schema = ?, is_active = ?, updated_at = ?
		WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		entity.DisplayName,
		entity.Description,
		entity.Schema,
		entity.IsActive,
		time.Now(),
		entity.ID.String(),
	)
	return err
}

func (r *customCRUDRepository) DeleteEntity(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM custom_crud_entities WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

func (r *customCRUDRepository) CreateData(ctx context.Context, data *models.CustomCRUDData) error {
	query := `INSERT INTO custom_crud_data (id, entity_id, data, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		data.ID.String(),
		data.EntityID.String(),
		data.Data,
		data.CreatedBy.String(),
		data.CreatedAt,
		data.UpdatedAt,
	)
	return err
}

func (r *customCRUDRepository) GetDataByID(ctx context.Context, id uuid.UUID) (*models.CustomCRUDData, error) {
	var d models.CustomCRUDData
	var entityIDStr, createdByStr, idStr string
	var updatedBy sql.NullString
	var deletedAt sql.NullTime

	query := `SELECT id, entity_id, data, created_by, updated_by, created_at, updated_at, deleted_at
		FROM custom_crud_data WHERE id = ? AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&idStr, &entityIDStr, &d.Data, &createdByStr, &updatedBy,
		&d.CreatedAt, &d.UpdatedAt, &deletedAt,
	)
	if err != nil {
		return nil, err
	}

	d.ID, _ = uuid.Parse(idStr)
	d.EntityID, _ = uuid.Parse(entityIDStr)
	d.CreatedBy, _ = uuid.Parse(createdByStr)
	if updatedBy.Valid {
		id, _ := uuid.Parse(updatedBy.String)
		d.UpdatedBy = &id
	}
	if deletedAt.Valid {
		d.DeletedAt = &deletedAt.Time
	}

	return &d, nil
}

func (r *customCRUDRepository) ListDataByEntity(ctx context.Context, entityID uuid.UUID, limit int, offset int) ([]*models.CustomCRUDData, error) {
	query := `SELECT id, entity_id, data, created_by, updated_by, created_at, updated_at, deleted_at
		FROM custom_crud_data WHERE entity_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, entityID.String(), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dataList []*models.CustomCRUDData
	for rows.Next() {
		var d models.CustomCRUDData
		var entityIDStr, createdByStr, idStr string
		var updatedBy sql.NullString
		var deletedAt sql.NullTime

		err := rows.Scan(&idStr, &entityIDStr, &d.Data, &createdByStr, &updatedBy,
			&d.CreatedAt, &d.UpdatedAt, &deletedAt)
		if err != nil {
			return nil, err
		}

		d.ID, _ = uuid.Parse(idStr)
		d.EntityID, _ = uuid.Parse(entityIDStr)
		d.CreatedBy, _ = uuid.Parse(createdByStr)
		if updatedBy.Valid {
			id, _ := uuid.Parse(updatedBy.String)
			d.UpdatedBy = &id
		}
		if deletedAt.Valid {
			d.DeletedAt = &deletedAt.Time
		}

		dataList = append(dataList, &d)
	}

	return dataList, rows.Err()
}

func (r *customCRUDRepository) UpdateData(ctx context.Context, data *models.CustomCRUDData) error {
	var updatedByStr *string
	if data.UpdatedBy != nil {
		str := data.UpdatedBy.String()
		updatedByStr = &str
	}

	query := `UPDATE custom_crud_data SET data = ?, updated_by = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query,
		data.Data,
		updatedByStr,
		time.Now(),
		data.ID.String(),
	)
	return err
}

func (r *customCRUDRepository) DeleteData(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE custom_crud_data SET deleted_at = datetime('now') WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id.String())
	return err
}

type adminActivityLogRepository struct {
	db *database.DB
}

func NewAdminActivityLogRepository(db *database.DB) AdminActivityLogRepository {
	return &adminActivityLogRepository{db: db}
}

func (r *adminActivityLogRepository) Create(ctx context.Context, log *models.AdminActivityLog) error {
	query := `INSERT INTO admin_activity_logs (id, admin_id, action, entity_type, entity_id, details, ip_address, user_agent, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		log.ID.String(),
		log.AdminID.String(),
		log.Action,
		log.EntityType,
		log.EntityID,
		log.Details,
		log.IPAddress,
		log.UserAgent,
		log.CreatedAt,
	)
	return err
}

func (r *adminActivityLogRepository) GetByAdminID(ctx context.Context, adminID uuid.UUID, limit int) ([]*models.AdminActivityLog, error) {
	query := `SELECT id, admin_id, action, entity_type, entity_id, details, ip_address, user_agent, created_at
		FROM admin_activity_logs WHERE admin_id = ? ORDER BY created_at DESC LIMIT ?`
	rows, err := r.db.QueryContext(ctx, query, adminID.String(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.AdminActivityLog
	for rows.Next() {
		var l models.AdminActivityLog
		var adminIDStr, idStr string
		var entityID, details, ipAddress, userAgent sql.NullString

		err := rows.Scan(&idStr, &adminIDStr, &l.Action, &l.EntityType, &entityID,
			&details, &ipAddress, &userAgent, &l.CreatedAt)
		if err != nil {
			return nil, err
		}

		l.ID, _ = uuid.Parse(idStr)
		l.AdminID, _ = uuid.Parse(adminIDStr)
		if entityID.Valid {
			l.EntityID = &entityID.String
		}
		if details.Valid {
			l.Details = &details.String
		}
		if ipAddress.Valid {
			l.IPAddress = &ipAddress.String
		}
		if userAgent.Valid {
			l.UserAgent = &userAgent.String
		}

		logs = append(logs, &l)
	}

	return logs, rows.Err()
}

func (r *adminActivityLogRepository) GetByEntityType(ctx context.Context, entityType string, limit int) ([]*models.AdminActivityLog, error) {
	query := `SELECT id, admin_id, action, entity_type, entity_id, details, ip_address, user_agent, created_at
		FROM admin_activity_logs WHERE entity_type = ? ORDER BY created_at DESC LIMIT ?`
	rows, err := r.db.QueryContext(ctx, query, entityType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.AdminActivityLog
	for rows.Next() {
		var l models.AdminActivityLog
		var adminIDStr, idStr string
		var entityID, details, ipAddress, userAgent sql.NullString

		err := rows.Scan(&idStr, &adminIDStr, &l.Action, &l.EntityType, &entityID,
			&details, &ipAddress, &userAgent, &l.CreatedAt)
		if err != nil {
			return nil, err
		}

		l.ID, _ = uuid.Parse(idStr)
		l.AdminID, _ = uuid.Parse(adminIDStr)
		if entityID.Valid {
			l.EntityID = &entityID.String
		}
		if details.Valid {
			l.Details = &details.String
		}
		if ipAddress.Valid {
			l.IPAddress = &ipAddress.String
		}
		if userAgent.Valid {
			l.UserAgent = &userAgent.String
		}

		logs = append(logs, &l)
	}

	return logs, rows.Err()
}

type userManagementActionRepository struct {
	db *database.DB
}

func NewUserManagementActionRepository(db *database.DB) UserManagementActionRepository {
	return &userManagementActionRepository{db: db}
}

func (r *userManagementActionRepository) Create(ctx context.Context, action *models.UserManagementAction) error {
	query := `INSERT INTO user_management_actions (id, admin_id, user_id, action_type, changes, reason, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		action.ID.String(),
		action.AdminID.String(),
		action.UserID.String(),
		action.ActionType,
		action.Changes,
		action.Reason,
		action.CreatedAt,
	)
	return err
}

func (r *userManagementActionRepository) GetByAdminID(ctx context.Context, adminID uuid.UUID, limit int) ([]*models.UserManagementAction, error) {
	query := `SELECT id, admin_id, user_id, action_type, changes, reason, created_at
		FROM user_management_actions WHERE admin_id = ? ORDER BY created_at DESC LIMIT ?`
	rows, err := r.db.QueryContext(ctx, query, adminID.String(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []*models.UserManagementAction
	for rows.Next() {
		var a models.UserManagementAction
		var adminIDStr, userIDStr, idStr string
		var changes, reason sql.NullString

		err := rows.Scan(&idStr, &adminIDStr, &userIDStr, &a.ActionType, &changes, &reason, &a.CreatedAt)
		if err != nil {
			return nil, err
		}

		a.ID, _ = uuid.Parse(idStr)
		a.AdminID, _ = uuid.Parse(adminIDStr)
		a.UserID, _ = uuid.Parse(userIDStr)
		if changes.Valid {
			a.Changes = &changes.String
		}
		if reason.Valid {
			a.Reason = &reason.String
		}

		actions = append(actions, &a)
	}

	return actions, rows.Err()
}

func (r *userManagementActionRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*models.UserManagementAction, error) {
	query := `SELECT id, admin_id, user_id, action_type, changes, reason, created_at
		FROM user_management_actions WHERE user_id = ? ORDER BY created_at DESC LIMIT ?`
	rows, err := r.db.QueryContext(ctx, query, userID.String(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []*models.UserManagementAction
	for rows.Next() {
		var a models.UserManagementAction
		var adminIDStr, userIDStr, idStr string
		var changes, reason sql.NullString

		err := rows.Scan(&idStr, &adminIDStr, &userIDStr, &a.ActionType, &changes, &reason, &a.CreatedAt)
		if err != nil {
			return nil, err
		}

		a.ID, _ = uuid.Parse(idStr)
		a.AdminID, _ = uuid.Parse(adminIDStr)
		a.UserID, _ = uuid.Parse(userIDStr)
		if changes.Valid {
			a.Changes = &changes.String
		}
		if reason.Valid {
			a.Reason = &reason.String
		}

		actions = append(actions, &a)
	}

	return actions, rows.Err()
}

