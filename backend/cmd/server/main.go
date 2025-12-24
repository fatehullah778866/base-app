package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"base-app-service/internal/cache"
	"base-app-service/internal/config"
	"base-app-service/internal/database"
	"base-app-service/internal/handlers"
	"base-app-service/internal/middleware"
	"base-app-service/internal/models"
	"base-app-service/internal/monitoring"
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
	
	// New repositories
	passwordResetRepo := repositories.NewPasswordResetRepository(db)
	settingsRepo := repositories.NewSettingsRepository(db)
	dashboardRepo := repositories.NewDashboardRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	messageRepo := repositories.NewMessageRepository(db)
	accountSwitchRepo := repositories.NewAccountSwitchRepository(db)
	searchRepo := repositories.NewSearchRepository(db)
	adminSettingsRepo := repositories.NewAdminSettingsRepository(db)
	customCRUDRepo := repositories.NewCustomCRUDRepository(db)
	// adminActivityLogRepo := repositories.NewAdminActivityLogRepository(db) // Reserved for future use
	// userManagementActionRepo := repositories.NewUserManagementActionRepository(db) // Reserved for future use

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
	passwordResetService := services.NewPasswordResetService(userRepo, passwordResetRepo, logger)
	activityLogService := services.NewActivityLogService(logRepo, logger)
	themeService := services.NewThemeService(themeRepo, logger)
	requestService := services.NewRequestService(requestRepo, logger)
	adminService := services.NewAdminService(userRepo, authService, activityLogService, requestRepo, logger)
	
	// New services
	settingsService := services.NewSettingsService(settingsRepo, userRepo, logger)
	dashboardService := services.NewDashboardService(dashboardRepo, logger)
	notificationService := services.NewNotificationService(notificationRepo, logger)
	messagingService := services.NewMessagingService(messageRepo, userRepo, logger)
	accountSwitchService := services.NewAccountSwitchService(accountSwitchRepo, userRepo, logger)
	searchService := services.NewSearchService(
		searchRepo,
		dashboardRepo,
		messageRepo,
		notificationRepo,
		customCRUDRepo,
		userRepo,
		logger,
	)
	adminSettingsService := services.NewAdminSettingsService(adminSettingsRepo, logger)
	customCRUDService := services.NewCustomCRUDService(customCRUDRepo, logger)
	
	// Email service
	emailConfig := services.GetEmailConfigFromEnv()
	emailService := services.NewEmailService(emailConfig, logger)
	
	// File service
	uploadDir := getEnv("UPLOAD_DIR", "uploads")
	fileService := services.NewFileService(services.FileUploadConfig{
		UploadDir: uploadDir,
		MaxSize:   10 * 1024 * 1024, // 10MB default
	}, logger)
	
	// Cache
	_ = cache.NewInMemoryCache(logger) // Reserved for future use
	
	// Monitoring
	metrics := monitoring.NewMetrics(logger)
	healthChecker := monitoring.NewHealthChecker(db, logger)

	// Background purge for soft-deleted users (5-day retention)
	startDeletedUserSweeper(ctx, userRepo, logger)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, passwordResetService, emailService, logger)
	userHandler := handlers.NewUserHandler(userRepo, requestRepo, themeRepo, sessionRepo, logger)
	themeHandler := handlers.NewThemeHandler(themeService, logger)
	adminHandler := handlers.NewAdminHandler(adminService, adminSettingsService, customCRUDService, logger)
	requestHandler := handlers.NewRequestHandler(requestService, logger)
	
	// New handlers
	settingsHandler := handlers.NewSettingsHandler(settingsService, sessionRepo, logger)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService, logger)
	notificationHandler := handlers.NewNotificationHandler(notificationService, logger)
	messagingHandler := handlers.NewMessagingHandler(messagingService, logger)
	accountSwitchHandler := handlers.NewAccountSwitchHandler(accountSwitchService, logger)
	searchHandler := handlers.NewSearchHandler(searchService, logger)
	fileUploadHandler := handlers.NewFileUploadHandler(fileService, logger)

	// Setup router
	router := mux.NewRouter()

	// Initialize rate limiter
	rateLimiter := middleware.NewInMemoryRateLimiter()

	// Middleware (order matters!)
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.RequestSizeLimitMiddleware(10 * 1024 * 1024)) // 10MB
	router.Use(middleware.CORSMiddleware())
	// CSRF middleware - skip for API endpoints
	csrfMiddleware := middleware.CSRFMiddleware()
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip CSRF for API endpoints
			if strings.HasPrefix(r.URL.Path, "/v1/") {
				next.ServeHTTP(w, r)
				return
			}
			csrfMiddleware(next).ServeHTTP(w, r)
		})
	})
	router.Use(monitoring.MetricsMiddleware(metrics))
	router.Use(middleware.RateLimitMiddleware(rateLimiter, 100, 1*time.Minute, logger)) // 100 req/min per IP/user
	router.Use(middleware.ErrorRecovery(logger))

	// Health check endpoints
	router.HandleFunc("/health", healthChecker.HealthCheck).Methods("GET")
	router.HandleFunc("/health/ready", healthChecker.ReadinessCheck).Methods("GET")
	router.HandleFunc("/health/live", healthChecker.LivenessCheck).Methods("GET")
	router.HandleFunc("/metrics", metrics.MetricsHandler).Methods("GET")

	// API v1 routes
	v1 := router.PathPrefix("/v1").Subrouter()

	// Public routes
	public := v1.PathPrefix("").Subrouter()
	public.HandleFunc("/auth/signup", authHandler.Signup).Methods("POST")
	public.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	public.HandleFunc("/auth/refresh", authHandler.RefreshToken).Methods("POST")
	public.HandleFunc("/auth/forgot-password", authHandler.ForgotPassword).Methods("POST")
	public.HandleFunc("/auth/reset-password", authHandler.ResetPassword).Methods("POST")
	public.HandleFunc("/admin/login", adminHandler.Login).Methods("POST")
	public.HandleFunc("/admin/verify-code", adminHandler.VerifyAdminCode).Methods("POST") // Verify admin code
	public.HandleFunc("/admin/create", adminHandler.CreateAdminPublic).Methods("POST") // Public admin creation with verification

	// Protected routes
	protected := v1.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret, logger))
	protected.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")
	protected.HandleFunc("/users/me", userHandler.GetCurrentUser).Methods("GET")
	protected.HandleFunc("/users/me", userHandler.UpdateProfile).Methods("PUT")
	protected.HandleFunc("/users/me/password", userHandler.ChangePassword).Methods("PUT")
	protected.HandleFunc("/users/me/export", userHandler.ExportData).Methods("GET")
	protected.HandleFunc("/users/me/delete", userHandler.RequestDeletion).Methods("POST")
	protected.HandleFunc("/users/me/settings/theme", themeHandler.GetTheme).Methods("GET")
	protected.HandleFunc("/users/me/settings/theme", themeHandler.UpdateTheme).Methods("PUT")
	protected.HandleFunc("/users/me/settings/theme/sync", themeHandler.SyncTheme).Methods("POST")
	protected.HandleFunc("/requests", requestHandler.Create).Methods("POST")
	protected.HandleFunc("/requests", requestHandler.ListMine).Methods("GET")
	
	// Settings routes
	protected.HandleFunc("/users/me/settings", settingsHandler.GetSettings).Methods("GET")
	protected.HandleFunc("/users/me/settings/sessions", settingsHandler.GetActiveSessions).Methods("GET")
	protected.HandleFunc("/users/me/settings/sessions/logout-all", settingsHandler.LogoutAllDevices).Methods("POST")
	protected.HandleFunc("/users/me/settings/profile", settingsHandler.UpdateProfileSettings).Methods("PUT")
	protected.HandleFunc("/users/me/settings/security", settingsHandler.UpdateSecuritySettings).Methods("PUT")
	protected.HandleFunc("/users/me/settings/privacy", settingsHandler.UpdatePrivacySettings).Methods("PUT")
	protected.HandleFunc("/users/me/settings/notifications", settingsHandler.UpdateNotificationSettings).Methods("PUT")
	protected.HandleFunc("/users/me/settings/preferences", settingsHandler.UpdateAccountPreferences).Methods("PUT")
	protected.HandleFunc("/users/me/settings/connected-accounts", settingsHandler.AddConnectedAccount).Methods("POST")
	protected.HandleFunc("/users/me/settings/connected-accounts", settingsHandler.RemoveConnectedAccount).Methods("DELETE")
	protected.HandleFunc("/users/me/settings/account/deactivate", settingsHandler.DeactivateAccount).Methods("POST")
	protected.HandleFunc("/users/me/settings/account/reactivate", settingsHandler.ReactivateAccount).Methods("POST")
	protected.HandleFunc("/users/me/settings/account/delete", settingsHandler.RequestAccountDeletion).Methods("POST")
	
	// Dashboard routes
	protected.HandleFunc("/dashboard/items", dashboardHandler.CreateItem).Methods("POST")
	protected.HandleFunc("/dashboard/items", dashboardHandler.ListItems).Methods("GET")
	protected.HandleFunc("/dashboard/items/{id}", dashboardHandler.GetItem).Methods("GET")
	protected.HandleFunc("/dashboard/items/{id}", dashboardHandler.UpdateItem).Methods("PUT")
	protected.HandleFunc("/dashboard/items/{id}", dashboardHandler.DeleteItem).Methods("DELETE")
	protected.HandleFunc("/dashboard/items/{id}/archive", dashboardHandler.SoftDeleteItem).Methods("POST")
	
	// Notification routes
	protected.HandleFunc("/notifications", notificationHandler.GetNotifications).Methods("GET")
	protected.HandleFunc("/notifications/unread-count", notificationHandler.GetUnreadCount).Methods("GET")
	protected.HandleFunc("/notifications/read", notificationHandler.MarkAsRead).Methods("POST")
	protected.HandleFunc("/notifications/read-all", notificationHandler.MarkAllAsRead).Methods("POST")
	protected.HandleFunc("/notifications", notificationHandler.DeleteNotification).Methods("DELETE")
	
	// Messaging routes
	protected.HandleFunc("/messages", messagingHandler.SendMessage).Methods("POST")
	protected.HandleFunc("/messages/conversations", messagingHandler.GetConversations).Methods("GET")
	protected.HandleFunc("/messages", messagingHandler.GetMessages).Methods("GET")
	protected.HandleFunc("/messages/read", messagingHandler.MarkAsRead).Methods("POST")
	protected.HandleFunc("/messages/unread-count", messagingHandler.GetUnreadCount).Methods("GET")
	
	// Account switching routes
	protected.HandleFunc("/account/switch", accountSwitchHandler.SwitchAccount).Methods("POST")
	protected.HandleFunc("/account/switch/history", accountSwitchHandler.GetSwitchHistory).Methods("GET")
	
	// Search routes
	protected.HandleFunc("/search", searchHandler.Search).Methods("GET", "POST")
	protected.HandleFunc("/search/history", searchHandler.GetSearchHistory).Methods("GET")
	protected.HandleFunc("/search/history", searchHandler.ClearSearchHistory).Methods("DELETE")
	
	// File upload routes
	protected.HandleFunc("/files/upload/image", fileUploadHandler.UploadImage).Methods("POST")
	protected.HandleFunc("/files/upload/document", fileUploadHandler.UploadDocument).Methods("POST")
	protected.HandleFunc("/files/download", fileUploadHandler.DownloadFile).Methods("GET")
	protected.HandleFunc("/files/delete", fileUploadHandler.DeleteFile).Methods("DELETE")

	// Serve uploaded files
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))

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
	
	// Admin settings routes
	adminProtected.HandleFunc("/settings", adminHandler.GetSettings).Methods("GET")
	adminProtected.HandleFunc("/settings", adminHandler.UpdateSettings).Methods("PUT")
	
	// Admin custom CRUD routes
	adminProtected.HandleFunc("/cruds/entities", adminHandler.CreateCRUDEntity).Methods("POST")
	adminProtected.HandleFunc("/cruds/entities", adminHandler.ListCRUDEntities).Methods("GET")
	adminProtected.HandleFunc("/cruds/entities/{id}", adminHandler.GetCRUDEntity).Methods("GET")
	adminProtected.HandleFunc("/cruds/entities/{id}", adminHandler.UpdateCRUDEntity).Methods("PUT")
	adminProtected.HandleFunc("/cruds/entities/{id}", adminHandler.DeleteCRUDEntity).Methods("DELETE")
	adminProtected.HandleFunc("/cruds/entities/{id}/data", adminHandler.CreateCRUDData).Methods("POST")
	adminProtected.HandleFunc("/cruds/entities/{id}/data", adminHandler.ListCRUDData).Methods("GET")
	adminProtected.HandleFunc("/cruds/data/{id}", adminHandler.GetCRUDData).Methods("GET")
	adminProtected.HandleFunc("/cruds/data/{id}", adminHandler.UpdateCRUDData).Methods("PUT")
	adminProtected.HandleFunc("/cruds/data/{id}", adminHandler.DeleteCRUDData).Methods("DELETE")
	
	// CRUD Templates routes (for easy entity creation)
	adminProtected.HandleFunc("/cruds/templates", adminHandler.GetCRUDTemplates).Methods("GET")
	adminProtected.HandleFunc("/cruds/templates/{name}", adminHandler.GetCRUDTemplate).Methods("GET")
	adminProtected.HandleFunc("/cruds/templates/{name}/create", adminHandler.CreateEntityFromTemplate).Methods("POST")
	
	// Enhanced admin user CRUD routes
	adminProtected.HandleFunc("/users", adminHandler.CreateUser).Methods("POST")
	adminProtected.HandleFunc("/users/{id}", adminHandler.UpdateUser).Methods("PUT")
	adminProtected.HandleFunc("/users/{id}", adminHandler.DeleteUser).Methods("DELETE")
	adminProtected.HandleFunc("/users/{id}/sessions", adminHandler.GetUserSessions).Methods("GET")
	adminProtected.HandleFunc("/users/{id}/sessions", adminHandler.RevokeUserSessions).Methods("DELETE")

	// Optional: Static frontend serving (can be disabled for API-only mode)
	// Backend works completely independently - frontend is optional
	// Set FRONTEND_DIR environment variable to enable frontend serving
	// If not set, backend runs in API-only mode (recommended for production)
	frontendDir := os.Getenv("FRONTEND_DIR")
	if frontendDir != "" {
		if _, err := os.Stat(frontendDir); err == nil {
			staticServer := http.FileServer(http.Dir(frontendDir))
			router.PathPrefix("/").Handler(staticServer)
			logger.Info("Serving frontend from directory", zap.String("dir", frontendDir))
		} else {
			logger.Warn("Frontend directory not found, running in API-only mode", zap.String("dir", frontendDir))
		}
	} else {
		// API-only mode - no frontend serving
		// This is the recommended mode for production
		// Frontend can be served separately (e.g., CDN, separate server, etc.)
		logger.Info("Running in API-only mode (no frontend serving)")
		logger.Info("Backend is ready to accept requests from any frontend")
		// Add a simple API info endpoint for root
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"Base App API","version":"1.0","docs":"/v1","status":"running"}`))
		}).Methods("GET")
	}

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

// Helper function for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
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
