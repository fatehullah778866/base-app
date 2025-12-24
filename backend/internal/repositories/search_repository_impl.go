package repositories

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
)

type searchRepository struct {
	db *database.DB
}

func NewSearchRepository(db *database.DB) SearchRepository {
	return &searchRepository{db: db}
}

func (r *searchRepository) SaveSearchHistory(ctx context.Context, history *models.SearchHistory) error {
	query := `INSERT INTO search_history (id, user_id, query, search_type, results_count, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query,
		history.ID.String(),
		history.UserID.String(),
		history.Query,
		history.SearchType,
		history.ResultsCount,
		history.CreatedAt,
	)
	return err
}

func (r *searchRepository) GetSearchHistory(ctx context.Context, userID uuid.UUID, limit int) ([]*models.SearchHistory, error) {
	if limit <= 0 {
		limit = 50
	}
	
	query := `SELECT id, user_id, query, search_type, results_count, created_at
		FROM search_history
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ?`
	
	rows, err := r.db.QueryContext(ctx, query, userID.String(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*models.SearchHistory
	for rows.Next() {
		var h models.SearchHistory
		var userIDStr, idStr string
		var searchType sql.NullString

		err := rows.Scan(&idStr, &userIDStr, &h.Query, &searchType, &h.ResultsCount, &h.CreatedAt)
		if err != nil {
			return nil, err
		}

		h.ID, _ = uuid.Parse(idStr)
		h.UserID, _ = uuid.Parse(userIDStr)
		if searchType.Valid {
			h.SearchType = &searchType.String
		}

		history = append(history, &h)
	}

	return history, rows.Err()
}

func (r *searchRepository) ClearSearchHistory(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM search_history WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID.String())
	return err
}

func (r *searchRepository) SearchDashboardItems(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.DashboardItem, error) {
	// Use FTS5 for full-text search
	ftsQuery := `SELECT id, user_id, title, description, category, status, priority, metadata, created_at, updated_at, deleted_at
		FROM dashboard_items_fts
		JOIN dashboard_items ON dashboard_items_fts.id = dashboard_items.id
		WHERE dashboard_items.user_id = ? AND dashboard_items_fts MATCH ?
		ORDER BY rank LIMIT ?`
	
	// Format query for FTS5 (space-separated terms)
	searchTerms := strings.Join(strings.Fields(query), " OR ")
	
	rows, err := r.db.QueryContext(ctx, ftsQuery, userID.String(), searchTerms, limit)
	if err != nil {
		// Fallback to simple LIKE search if FTS5 fails
		fallbackQuery := `SELECT id, user_id, title, description, category, status, priority, metadata, created_at, updated_at, deleted_at
			FROM dashboard_items
			WHERE user_id = ? AND (title LIKE ? OR description LIKE ?)
			ORDER BY created_at DESC LIMIT ?`
		searchPattern := "%" + query + "%"
		rows, err = r.db.QueryContext(ctx, fallbackQuery, userID.String(), searchPattern, searchPattern, limit)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	var items []*models.DashboardItem
	for rows.Next() {
		var item models.DashboardItem
		var userIDStr, idStr string
		var description, category, metadata sql.NullString
		var deletedAt sql.NullTime

		err := rows.Scan(&idStr, &userIDStr, &item.Title, &description, &category,
			&item.Status, &item.Priority, &metadata, &item.CreatedAt, &item.UpdatedAt, &deletedAt)
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

func (r *searchRepository) SearchMessages(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.Message, error) {
	// Use FTS5 for full-text search
	ftsQuery := `SELECT id, sender_id, recipient_id, subject, content, is_read, read_at, is_archived, archived_at, metadata, created_at, updated_at
		FROM messages_fts
		JOIN messages ON messages_fts.id = messages.id
		WHERE (messages.sender_id = ? OR messages.recipient_id = ?) AND messages_fts MATCH ?
		ORDER BY rank LIMIT ?`
	
	searchTerms := strings.Join(strings.Fields(query), " OR ")
	
	rows, err := r.db.QueryContext(ctx, ftsQuery, userID.String(), userID.String(), searchTerms, limit)
	if err != nil {
		// Fallback to simple LIKE search
		fallbackQuery := `SELECT id, sender_id, recipient_id, subject, content, is_read, read_at, is_archived, archived_at, metadata, created_at, updated_at
			FROM messages
			WHERE (sender_id = ? OR recipient_id = ?) AND (subject LIKE ? OR content LIKE ?)
			ORDER BY created_at DESC LIMIT ?`
		searchPattern := "%" + query + "%"
		rows, err = r.db.QueryContext(ctx, fallbackQuery, userID.String(), userID.String(), searchPattern, searchPattern, limit)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var m models.Message
		var senderIDStr, recipientIDStr, idStr string
		var subject, metadata sql.NullString
		var readAt, archivedAt sql.NullTime

		err := rows.Scan(&idStr, &senderIDStr, &recipientIDStr, &subject, &m.Content,
			&m.IsRead, &readAt, &m.IsArchived, &archivedAt, &metadata,
			&m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}

		m.ID, _ = uuid.Parse(idStr)
		m.SenderID, _ = uuid.Parse(senderIDStr)
		m.RecipientID, _ = uuid.Parse(recipientIDStr)
		if subject.Valid {
			m.Subject = &subject.String
		}
		if readAt.Valid {
			m.ReadAt = &readAt.Time
		}
		if archivedAt.Valid {
			m.ArchivedAt = &archivedAt.Time
		}
		if metadata.Valid {
			m.Metadata = &metadata.String
		}

		messages = append(messages, &m)
	}

	return messages, rows.Err()
}

func (r *searchRepository) SearchUsers(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.User, error) {
	// Simple LIKE search for users (name, email)
	searchPattern := "%" + query + "%"
	searchQuery := `SELECT id, email, email_verified, email_verification_token, password_hash, password_changed_at,
		name, first_name, last_name, photo_url, role, phone, phone_verified, status, signup_source,
		created_at, updated_at, last_login_at
		FROM users
		WHERE (name LIKE ? OR email LIKE ?) AND id != ? AND status IN ('active', 'pending')
		ORDER BY name LIMIT ?`
	
	rows, err := r.db.QueryContext(ctx, searchQuery, searchPattern, searchPattern, userID.String(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		var idStr string
		var firstName, lastName, photoURL, phone, signupSource sql.NullString
		var lastLoginAt sql.NullTime

		err := rows.Scan(&idStr, &u.Email, &u.EmailVerified, &u.EmailVerificationToken,
			&u.PasswordHash, &u.PasswordChangedAt, &u.Name, &firstName, &lastName,
			&photoURL, &u.Role, &phone, &u.PhoneVerified, &u.Status, &signupSource,
			&u.CreatedAt, &u.UpdatedAt, &lastLoginAt)
		if err != nil {
			return nil, err
		}

		u.ID, _ = uuid.Parse(idStr)
		if firstName.Valid {
			u.FirstName = &firstName.String
		}
		if lastName.Valid {
			u.LastName = &lastName.String
		}
		if photoURL.Valid {
			u.PhotoURL = &photoURL.String
		}
		if phone.Valid {
			u.Phone = &phone.String
		}
		if signupSource.Valid {
			u.SignupSource = &signupSource.String
		}
		if lastLoginAt.Valid {
			u.LastLoginAt = &lastLoginAt.Time
		}

		users = append(users, &u)
	}

	return users, rows.Err()
}

func (r *searchRepository) SearchUsersByLocation(ctx context.Context, userID uuid.UUID, country, city *string, limit int) ([]*models.User, error) {
	// Search users by location (simplified - in production, you'd have a location table)
	// For now, we'll search in user settings or sessions for location data
	query := `SELECT DISTINCT u.id, u.email, u.email_verified, u.email_verification_token, u.password_hash, u.password_changed_at,
		u.name, u.first_name, u.last_name, u.photo_url, u.role, u.phone, u.phone_verified, u.status, u.signup_source,
		u.created_at, u.updated_at, u.last_login_at
		FROM users u
		LEFT JOIN sessions s ON s.user_id = u.id
		WHERE u.id != ? AND u.status IN ('active', 'pending')`
	
	args := []interface{}{userID.String()}
	
	if country != nil {
		query += ` AND s.location_country LIKE ?`
		args = append(args, "%"+*country+"%")
	}
	if city != nil {
		query += ` AND s.location_city LIKE ?`
		args = append(args, "%"+*city+"%")
	}
	
	query += ` ORDER BY u.name LIMIT ?`
	args = append(args, limit)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var u models.User
		var idStr string
		var firstName, lastName, photoURL, phone, signupSource sql.NullString
		var lastLoginAt sql.NullTime

		err := rows.Scan(&idStr, &u.Email, &u.EmailVerified, &u.EmailVerificationToken,
			&u.PasswordHash, &u.PasswordChangedAt, &u.Name, &firstName, &lastName,
			&photoURL, &u.Role, &phone, &u.PhoneVerified, &u.Status, &signupSource,
			&u.CreatedAt, &u.UpdatedAt, &lastLoginAt)
		if err != nil {
			return nil, err
		}

		u.ID, _ = uuid.Parse(idStr)
		if firstName.Valid {
			u.FirstName = &firstName.String
		}
		if lastName.Valid {
			u.LastName = &lastName.String
		}
		if photoURL.Valid {
			u.PhotoURL = &photoURL.String
		}
		if phone.Valid {
			u.Phone = &phone.String
		}
		if signupSource.Valid {
			u.SignupSource = &signupSource.String
		}
		if lastLoginAt.Valid {
			u.LastLoginAt = &lastLoginAt.Time
		}

		users = append(users, &u)
	}

	return users, rows.Err()
}

func (r *searchRepository) SearchNotifications(ctx context.Context, userID uuid.UUID, query string, limit int) ([]*models.Notification, error) {
	searchPattern := "%" + query + "%"
	searchQuery := `SELECT id, user_id, title, message, type, is_read, read_at, metadata, created_at
		FROM notifications
		WHERE user_id = ? AND (title LIKE ? OR message LIKE ?)
		ORDER BY created_at DESC
		LIMIT ?`
	
	rows, err := r.db.QueryContext(ctx, searchQuery, userID.String(), searchPattern, searchPattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		var n models.Notification
		var userIDStr, idStr string
		var metadata sql.NullString
		var readAt sql.NullTime

		err := rows.Scan(&idStr, &userIDStr, &n.Title, &n.Message, &n.Type, &n.IsRead, &readAt, &metadata, &n.CreatedAt)
		if err != nil {
			return nil, err
		}

		n.ID, _ = uuid.Parse(idStr)
		n.UserID, _ = uuid.Parse(userIDStr)
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
