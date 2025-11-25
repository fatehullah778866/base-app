package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"base-app-service/internal/config"
	"base-app-service/internal/database"
	"base-app-service/internal/handlers"
	"base-app-service/internal/middleware"
	"base-app-service/internal/repositories"
	"base-app-service/internal/services"
)

func main() {
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
		Host:                 cfg.Database.Host,
		Port:                 cfg.Database.Port,
		User:                 cfg.Database.User,
		Password:             cfg.Database.Password,
		Name:                 cfg.Database.Name,
		SSLMode:              cfg.Database.SSLMode,
		MaxConnections:       cfg.Database.MaxConnections,
		MaxIdleConnections:   cfg.Database.MaxIdleConnections,
		ConnectionMaxLifetime: cfg.Database.ConnectionMaxLifetime,
	}
	db, err := database.NewConnection(dbConfig, logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	deviceRepo := repositories.NewDeviceRepository(db)
	themeRepo := repositories.NewThemeRepository(db)

	// Initialize services
	authService := services.NewAuthService(
		userRepo, sessionRepo, deviceRepo,
		cfg.JWT.Secret, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry,
		logger,
	)
	themeService := services.NewThemeService(themeRepo, logger)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, logger)
	userHandler := handlers.NewUserHandler(userRepo, logger)
	themeHandler := handlers.NewThemeHandler(themeService, logger)

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

	// Protected routes
	protected := v1.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret, logger))
	protected.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")
	protected.HandleFunc("/users/me", userHandler.GetCurrentUser).Methods("GET")
	protected.HandleFunc("/users/me", userHandler.UpdateProfile).Methods("PUT")
	protected.HandleFunc("/users/me/settings/theme", themeHandler.GetTheme).Methods("GET")
	protected.HandleFunc("/users/me/settings/theme", themeHandler.UpdateTheme).Methods("PUT")
	protected.HandleFunc("/users/me/settings/theme/sync", themeHandler.SyncTheme).Methods("POST")

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

