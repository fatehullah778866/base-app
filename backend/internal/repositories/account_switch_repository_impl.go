package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
)

type accountSwitchRepository struct {
	db *database.DB
}

func NewAccountSwitchRepository(db *database.DB) AccountSwitchRepository {
	return &accountSwitchRepository{db: db}
}

func (r *accountSwitchRepository) Create(ctx context.Context, switchRecord *models.AccountSwitch) error {
	query := `
		INSERT INTO account_switches (id, user_id, switched_to_user_id, switched_to_role, switched_from_role, reason, ip_address, user_agent, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	var switchedToUserID *string
	if switchRecord.SwitchedToUserID != nil {
		switchedToUserIDStr := switchRecord.SwitchedToUserID.String()
		switchedToUserID = &switchedToUserIDStr
	}

	_, err := r.db.ExecContext(ctx, query,
		switchRecord.ID.String(),
		switchRecord.UserID.String(),
		switchedToUserID,
		switchRecord.SwitchedToRole,
		switchRecord.SwitchedFromRole,
		switchRecord.Reason,
		switchRecord.IPAddress,
		switchRecord.UserAgent,
		switchRecord.CreatedAt,
	)
	return err
}

func (r *accountSwitchRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]*models.AccountSwitch, error) {
	query := `SELECT id, user_id, switched_to_user_id, switched_to_role, switched_from_role, reason, ip_address, user_agent, created_at
		FROM account_switches WHERE user_id = ? ORDER BY created_at DESC LIMIT ?`
	rows, err := r.db.QueryContext(ctx, query, userID.String(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var switches []*models.AccountSwitch
	for rows.Next() {
		var s models.AccountSwitch
		var userIDStr, idStr string
		var switchedToUserID, switchedToRole, switchedFromRole, reason, ipAddress, userAgent sql.NullString

		err := rows.Scan(&idStr, &userIDStr, &switchedToUserID, &switchedToRole, &switchedFromRole,
			&reason, &ipAddress, &userAgent, &s.CreatedAt)
		if err != nil {
			return nil, err
		}

		s.ID, _ = uuid.Parse(idStr)
		s.UserID, _ = uuid.Parse(userIDStr)
		if switchedToUserID.Valid {
			id, _ := uuid.Parse(switchedToUserID.String)
			s.SwitchedToUserID = &id
		}
		if switchedToRole.Valid {
			s.SwitchedToRole = &switchedToRole.String
		}
		if switchedFromRole.Valid {
			s.SwitchedFromRole = &switchedFromRole.String
		}
		if reason.Valid {
			s.Reason = &reason.String
		}
		if ipAddress.Valid {
			s.IPAddress = &ipAddress.String
		}
		if userAgent.Valid {
			s.UserAgent = &userAgent.String
		}

		switches = append(switches, &s)
	}

	return switches, rows.Err()
}

