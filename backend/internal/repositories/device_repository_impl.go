package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"base-app-service/internal/database"
	"base-app-service/internal/models"
)

type deviceRepository struct {
	db *database.DB
}

func NewDeviceRepository(db *database.DB) DeviceRepository {
	return &deviceRepository{db: db}
}

func (r *deviceRepository) Create(ctx context.Context, device *models.Device) error {
	query := `
		INSERT INTO user_devices (
			id, user_id, device_id, device_name, device_type, os, browser,
			ip_address, location_country, location_city, is_trusted,
			last_used_at, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.DB.ExecContext(ctx, query,
		device.ID, device.UserID, device.DeviceID, device.DeviceName,
		device.DeviceType, device.OS, device.Browser, device.IPAddress,
		device.LocationCountry, device.LocationCity, device.IsTrusted,
		device.LastUsedAt, device.CreatedAt,
	)

	return err
}

func (r *deviceRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Device, error) {
	query := `
		SELECT id, user_id, device_id, device_name, device_type, os, browser,
			ip_address, location_country, location_city, is_trusted,
			trusted_at, last_used_at, created_at
		FROM user_devices
		WHERE id = ?
	`

	device := &models.Device{}
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&device.ID, &device.UserID, &device.DeviceID, &device.DeviceName,
		&device.DeviceType, &device.OS, &device.Browser, &device.IPAddress,
		&device.LocationCountry, &device.LocationCity, &device.IsTrusted,
		&device.TrustedAt, &device.LastUsedAt, &device.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("device not found")
	}
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (r *deviceRepository) GetByDeviceID(ctx context.Context, userID uuid.UUID, deviceID *string) (*models.Device, error) {
	if deviceID == nil {
		return nil, fmt.Errorf("device_id is required")
	}

	query := `
		SELECT id, user_id, device_id, device_name, device_type, os, browser,
			ip_address, location_country, location_city, is_trusted,
			trusted_at, last_used_at, created_at
		FROM user_devices
		WHERE user_id = ? AND device_id = ?
	`

	device := &models.Device{}
	err := r.db.DB.QueryRowContext(ctx, query, userID, *deviceID).Scan(
		&device.ID, &device.UserID, &device.DeviceID, &device.DeviceName,
		&device.DeviceType, &device.OS, &device.Browser, &device.IPAddress,
		&device.LocationCountry, &device.LocationCity, &device.IsTrusted,
		&device.TrustedAt, &device.LastUsedAt, &device.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not an error, just doesn't exist
	}
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (r *deviceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Device, error) {
	query := `
		SELECT id, user_id, device_id, device_name, device_type, os, browser,
			ip_address, location_country, location_city, is_trusted,
			trusted_at, last_used_at, created_at
		FROM user_devices
		WHERE user_id = ?
		ORDER BY last_used_at DESC
	`

	rows, err := r.db.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*models.Device
	for rows.Next() {
		device := &models.Device{}
		err := rows.Scan(
			&device.ID, &device.UserID, &device.DeviceID, &device.DeviceName,
			&device.DeviceType, &device.OS, &device.Browser, &device.IPAddress,
			&device.LocationCountry, &device.LocationCity, &device.IsTrusted,
			&device.TrustedAt, &device.LastUsedAt, &device.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func (r *deviceRepository) Update(ctx context.Context, device *models.Device) error {
	query := `
		UPDATE user_devices
		SET device_name = ?, last_used_at = ?, is_trusted = ?, trusted_at = ?
		WHERE id = ?
	`

	_, err := r.db.DB.ExecContext(ctx, query,
		device.DeviceName, device.LastUsedAt,
		device.IsTrusted, device.TrustedAt, device.ID,
	)

	return err
}

func (r *deviceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM user_devices WHERE id = ?`
	_, err := r.db.DB.ExecContext(ctx, query, id)
	return err
}
