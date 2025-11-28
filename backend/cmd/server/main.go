package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"base-app-service/internal/config"
	"base-app-service/internal/database"
	"base-app-service/internal/handlers"
	"base-app-service/internal/middleware"
	"base-app-service/internal/models"
	"base-app-service/internal/repositories"
	"base-app-service/internal/services"
	"base-app-service/pkg/auth"
)

func main() {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initialize logger
	var logger *zap.Logger
	if cfg.Logging.Format == "json" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	// Connect to database
	dbConfig := database.DatabaseConfig{
		Driver:                cfg.Database.Driver,
		Host:                  cfg.Database.Host,
		Port:                  cfg.Database.Port,
		User:                  cfg.Database.User,
		Password:              cfg.Database.Password,
		Name:                  cfg.Database.Name,
		SSLMode:               cfg.Database.SSLMode,
		SQLitePath:            cfg.Database.SQLitePath,
		MaxConnections:        cfg.Database.MaxConnections,
		MaxIdleConnections:    cfg.Database.MaxIdleConnections,
		ConnectionMaxLifetime: cfg.Database.ConnectionMaxLifetime,
	}
	db, err := database.NewConnection(dbConfig, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := db.RunMigrations("migrations"); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	deviceRepo := repositories.NewDeviceRepository(db)
	themeRepo := repositories.NewThemeRepository(db)
	logRepo := repositories.NewActivityLogRepository(db)
	requestRepo := repositories.NewAccessRequestRepository(db)

	// Seed default admin account
	if err := seedDefaultAdmin(context.Background(), userRepo, logger); err != nil {
		logger.Warn("Failed to seed default admin", zap.Error(err))
	}

	// Initialize services
	authService := services.NewAuthService(
		userRepo, sessionRepo, deviceRepo,
		cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry,
		logger,
	)
	activityLogService := services.NewActivityLogService(logRepo, logger)
	themeService := services.NewThemeService(themeRepo, logger)
	requestService := services.NewRequestService(requestRepo, logger)
	adminService := services.NewAdminService(userRepo, authService, activityLogService, requestRepo, logger)

	// Background purge for soft-deleted users (5-day retention)
	startDeletedUserSweeper(ctx, userRepo, logger)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, logger)
	userHandler := handlers.NewUserHandler(userRepo, requestRepo, themeRepo, sessionRepo, logger)
	themeHandler := handlers.NewThemeHandler(themeService, logger)
	adminHandler := handlers.NewAdminHandler(adminService, logger)
	requestHandler := handlers.NewRequestHandler(requestService, logger)

	// Setup router
	router := mux.NewRouter()

	// Middleware
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.ErrorRecovery(logger))

	// Health check
	router.HandleFunc("/health", healthCheck).Methods("GET")

	// API v1 routes
	v1 := router.PathPrefix("/v1").Subrouter()

	// Public routes
	public := v1.PathPrefix("").Subrouter()
	public.HandleFunc("/auth/signup", authHandler.Signup).Methods("POST")
	public.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	public.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST")
	public.HandleFunc("/admin/login", adminHandler.Login).Methods("POST")

	// Protected routes
	protected := v1.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret, logger))
	protected.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")
	protected.HandleFunc("/users/me", userHandler.GetCurrentUser).Methods("GET")
	protected.HandleFunc("/users/me", userHandler.UpdateProfile).Methods("PUT")
	protected.HandleFunc("/users/me/export", userHandler.ExportData).Methods("GET")
	protected.HandleFunc("/users/me/delete", userHandler.RequestDeletion).Methods("POST")
	protected.HandleFunc("/users/me/settings/theme", themeHandler.GetTheme).Methods("GET")
	protected.HandleFunc("/users/me/settings/theme", themeHandler.UpdateTheme).Methods("PUT")
	protected.HandleFunc("/users/me/settings/theme/sync", themeHandler.SyncTheme).Methods("POST")
	protected.HandleFunc("/requests", requestHandler.Create).Methods("POST")
	protected.HandleFunc("/requests", requestHandler.ListMine).Methods("GET")

	// Admin protected routes
	adminProtected := v1.PathPrefix("/admin").Subrouter()
	adminProtected.Use(middleware.AuthMiddleware(cfg.JWT.Secret, logger))
	adminProtected.Use(middleware.RequireRole("admin", logger))
	adminProtected.HandleFunc("/users", adminHandler.ListUsers).Methods("GET")
	adminProtected.HandleFunc("/users/{id}", adminHandler.GetUser).Methods("GET")
	adminProtected.HandleFunc("/users/{id}/status", adminHandler.UpdateUserStatus).Methods("POST")
	adminProtected.HandleFunc("/logs", adminHandler.ListLogs).Methods("GET")
	adminProtected.HandleFunc("/admins", adminHandler.AddAdmin).Methods("POST")
	adminProtected.HandleFunc("/admins", adminHandler.ListAdmins).Methods("GET")
	adminProtected.HandleFunc("/requests", adminHandler.ListRequests).Methods("GET")
	adminProtected.HandleFunc("/requests/{id}/status", adminHandler.UpdateRequestStatus).Methods("POST")

	// Static frontend (served from ../frontend relative to backend/)
	frontendDir := os.Getenv("FRONTEND_DIR")
	if frontendDir == "" {
		defaultDir := filepath.Clean(filepath.Join("..", "frontend"))
		if _, err := os.Stat(defaultDir); err == nil {
			frontendDir = defaultDir
		} else {
			frontendDir = "frontend"
		}
	}
	staticServer := http.FileServer(http.Dir(frontendDir))
	router.PathPrefix("/").Handler(staticServer)

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		logger.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"healthy"}`)
}

func seedDefaultAdmin(ctx context.Context, userRepo repositories.UserRepository, logger *zap.Logger) error {
	const adminEmail = "admin@gmail.com"
	const adminPassword = "admin123"
	const adminName = "Admin"

	existing, _ := userRepo.GetByEmail(ctx, adminEmail)
	if existing != nil {
		return nil
	}

	hash, err := auth.HashPassword(adminPassword)
	if err != nil {
		return err
	}

	now := time.Now()
	admin := &models.User{
		ID:                uuid.New(),
		Email:             adminEmail,
		PasswordHash:      hash,
		Name:              adminName,
		Status:            "active",
		Role:              "admin",
		PasswordChangedAt: now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := userRepo.Create(ctx, admin); err != nil {
		return err
	}

	logger.Info("Seeded default admin account", zap.String("email", adminEmail))
	return nil
}

// startDeletedUserSweeper removes users that have been soft-deleted for more than 5 days.
func startDeletedUserSweeper(ctx context.Context, userRepo repositories.UserRepository, logger *zap.Logger) {
	run := func() {
		cutoff := time.Now().Add(-5 * 24 * time.Hour)
		if err := userRepo.PurgeDeletedBefore(context.Background(), cutoff); err != nil {
			logger.Warn("Failed to purge soft-deleted users", zap.Error(err))
		}
	}

	run()

	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				run()
			}
		}
	}()
}
