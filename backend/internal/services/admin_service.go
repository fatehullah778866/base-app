package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
	"base-app-service/pkg/auth"
)

type AdminService struct {
	userRepo    repositories.UserRepository
	authService *AuthService
	logService  *ActivityLogService
	requestRepo repositories.AccessRequestRepository
	logger      *zap.Logger
}

type AdminLoginRequest struct {
	Email      string
	Password   string
	IPAddress  *string
	UserAgent  *string
	DeviceID   *string
	DeviceName *string
}

type CreateAdminRequest struct {
	Email    string
	Name     string
	Password string
}

const defaultAdminEmail = "admin@gmail.com"
const defaultAdminPassword = "admin123"
const defaultAdminName = "Admin"

func NewAdminService(
	userRepo repositories.UserRepository,
	authService *AuthService,
	logService *ActivityLogService,
	requestRepo repositories.AccessRequestRepository,
	logger *zap.Logger,
) *AdminService {
	return &AdminService{
		userRepo:    userRepo,
		authService: authService,
		logService:  logService,
		requestRepo: requestRepo,
		logger:      logger,
	}
}

func (s *AdminService) Login(ctx context.Context, req AdminLoginRequest) (*models.User, *models.Session, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if strings.EqualFold(req.Email, defaultAdminEmail) && req.Password == defaultAdminPassword {
			created, createErr := s.ensureDefaultAdmin(ctx)
			if createErr != nil {
				s.logger.Warn("failed to auto-create default admin on login", zap.Error(createErr))
				return nil, nil, errors.New("invalid credentials")
			}
			user = created
		} else {
			return nil, nil, errors.New("invalid credentials")
		}
	}

	// Auto-repair default admin account if credentials match defaults but stored state is off.
	if strings.EqualFold(user.Email, defaultAdminEmail) && req.Password == defaultAdminPassword {
		if repaired, err := s.repairDefaultAdmin(ctx, user); err == nil {
			user = repaired
		} else {
			s.logger.Warn("failed to repair default admin", zap.Error(err))
		}
	}

	if strings.ToLower(user.Role) != "admin" {
		return nil, nil, errors.New("invalid credentials")
	}

	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, nil, errors.New("invalid credentials")
	}

	if user.Status != "active" && user.Status != "pending" {
		return nil, nil, errors.New("account is not active")
	}

	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(ctx, user)

	session, err := s.authService.createSession(ctx, user, req.IPAddress, req.UserAgent, req.DeviceID, req.DeviceName)
	if err != nil {
		return nil, nil, err
	}

	s.logService.Record(ctx, &user.ID, "admin", "admin_login", strPtr("admin"), strPtr(user.ID.String()), nil)

	return user, session, nil
}

func (s *AdminService) AddAdmin(ctx context.Context, actorID uuid.UUID, req CreateAdminRequest) (*models.User, error) {
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("admin already exists")
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	admin := &models.User{
		ID:                uuid.New(),
		Email:             req.Email,
		PasswordHash:      passwordHash,
		Name:              req.Name,
		Status:            "active",
		Role:              "admin",
		PasswordChangedAt: now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := s.userRepo.Create(ctx, admin); err != nil {
		return nil, err
	}

	s.logService.Record(ctx, &actorID, "admin", "admin_created", strPtr("admin"), strPtr(admin.ID.String()), nil)

	return admin, nil
}

type CreateUserRequest struct {
	Email    string
	Name     string
	Password string
	Role     string
	Status   string
}

func (s *AdminService) CreateUser(ctx context.Context, actorID uuid.UUID, req CreateUserRequest) (*models.User, error) {
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	if req.Status == "" {
		req.Status = "active"
	}
	if req.Role == "" {
		req.Role = "user"
	}

	now := time.Now()
	user := &models.User{
		ID:                uuid.New(),
		Email:             req.Email,
		PasswordHash:      passwordHash,
		Name:              req.Name,
		Status:            req.Status,
		Role:              req.Role,
		PasswordChangedAt: now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Only log if actorID is valid (not nil)
	if actorID != uuid.Nil {
		s.logService.Record(ctx, &actorID, "admin", "user_created", strPtr("user"), strPtr(user.ID.String()), nil)
	} else {
		// Log as system action for public admin creation
		s.logService.Record(ctx, nil, "system", "admin_created_public", strPtr("admin"), strPtr(user.ID.String()), nil)
	}

	return user, nil
}

func (s *AdminService) ListUsers(ctx context.Context, search string) ([]*models.User, error) {
	return s.userRepo.List(ctx, search)
}

func (s *AdminService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *AdminService) ListAdmins(ctx context.Context, search string) ([]*models.User, error) {
	users, err := s.userRepo.List(ctx, search)
	if err != nil {
		return nil, err
	}
	var admins []*models.User
	for _, u := range users {
		if strings.ToLower(u.Role) == "admin" {
			admins = append(admins, u)
		}
	}
	return admins, nil
}

func (s *AdminService) SetUserStatus(ctx context.Context, actorID uuid.UUID, userID uuid.UUID, status string) error {
	if status != "active" && status != "disabled" {
		return errors.New("invalid status")
	}
	if err := s.userRepo.SetStatus(ctx, userID, status); err != nil {
		return err
	}

	s.logService.Record(ctx, &actorID, "admin", "user_status_changed", strPtr("user"), strPtr(userID.String()), map[string]interface{}{
		"status": status,
	})
	return nil
}

func (s *AdminService) ListLogs(ctx context.Context, limit int) ([]*models.ActivityLog, error) {
	return s.logService.List(ctx, limit)
}

func (s *AdminService) ListRequests(ctx context.Context, status *string) ([]*models.AccessRequest, error) {
	return s.requestRepo.List(ctx, status)
}

func (s *AdminService) UpdateRequestStatus(ctx context.Context, actorID uuid.UUID, id string, status string, feedback *string) (*models.AccessRequest, error) {
	status = strings.ToLower(status)
	if status != "approved" && status != "rejected" && status != "pending" {
		return nil, errors.New("invalid status")
	}

	req, err := s.requestRepo.UpdateStatus(ctx, id, status, feedback)
	if err != nil {
		return nil, err
	}

	s.logService.Record(ctx, &actorID, "admin", "request_status_changed", strPtr("request"), strPtr(id), map[string]interface{}{
		"status":   status,
		"feedback": feedback,
	})
	return req, nil
}

func strPtr(value string) *string {
	return &value
}

func (s *AdminService) ensureDefaultAdmin(ctx context.Context) (*models.User, error) {
	now := time.Now()
	passwordHash, err := auth.HashPassword(defaultAdminPassword)
	if err != nil {
		return nil, err
	}
	admin := &models.User{
		ID:                uuid.New(),
		Email:             defaultAdminEmail,
		PasswordHash:      passwordHash,
		Name:              defaultAdminName,
		Status:            "active",
		Role:              "admin",
		PasswordChangedAt: now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := s.userRepo.Create(ctx, admin); err != nil {
		return nil, err
	}
	s.logger.Info("Auto-created default admin during login")
	return admin, nil
}

func (s *AdminService) repairDefaultAdmin(ctx context.Context, user *models.User) (*models.User, error) {
	now := time.Now()
	changed := false

	if strings.ToLower(user.Role) != "admin" {
		user.Role = "admin"
		changed = true
	}
	if user.Status != "active" && user.Status != "pending" {
		user.Status = "active"
		changed = true
	}
	if user.Name == "" {
		user.Name = defaultAdminName
		changed = true
	}
	if !auth.CheckPasswordHash(defaultAdminPassword, user.PasswordHash) {
		hash, err := auth.HashPassword(defaultAdminPassword)
		if err != nil {
			return user, err
		}
		user.PasswordHash = hash
		user.PasswordChangedAt = now
		changed = true
	}
	if changed {
		user.UpdatedAt = now
		if err := s.userRepo.Update(ctx, user); err != nil {
			return user, err
		}
		s.logger.Info("Repaired default admin account on login")
	}
	return user, nil
}
