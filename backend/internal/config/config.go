package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	JWT       JWTConfig
	Webhook   WebhookConfig
	RateLimit RateLimitConfig
	Logging   LoggingConfig
	S3        S3Config
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Driver                string
	Host                  string
	Port                  int
	User                  string
	Password              string
	Name                  string
	SSLMode               string
	SQLitePath            string
	MaxConnections        int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret             string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

type WebhookConfig struct {
	Secret                 string
	MaxRetries             int
	RetryBackoffMultiplier float64
}

type RateLimitConfig struct {
	Enabled   bool
	KeyPrefix string
}

type LoggingConfig struct {
	Level  string
	Format string
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Driver:                getEnv("DB_DRIVER", "sqlite"),
			Host:                  getEnv("DB_HOST", "localhost"),
			Port:                  getEnvAsInt("DB_PORT", 5432),
			User:                  getEnv("DB_USER", "baseapp"),
			Password:              getEnv("DB_PASSWORD", ""),
			Name:                  getEnv("DB_NAME", "base_app.db"),
			SSLMode:               getEnv("DB_SSL_MODE", "disable"),
			SQLitePath:            getEnv("DB_SQLITE_PATH", "file:app.db?_pragma=foreign_keys(ON)"),
			MaxConnections:        getEnvAsInt("DB_MAX_CONNECTIONS", 25),
			MaxIdleConnections:    getEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 5),
			ConnectionMaxLifetime: getEnvAsDuration("DB_CONNECTION_MAX_LIFETIME", 300*time.Second),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", "change-me-in-production"),
			AccessTokenExpiry:  getEnvAsDuration("JWT_ACCESS_TOKEN_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry: getEnvAsDuration("JWT_REFRESH_TOKEN_EXPIRY", 30*24*time.Hour),
		},
		Webhook: WebhookConfig{
			Secret:                 getEnv("WEBHOOK_SECRET", ""),
			MaxRetries:             getEnvAsInt("WEBHOOK_MAX_RETRIES", 3),
			RetryBackoffMultiplier: getEnvAsFloat("WEBHOOK_RETRY_BACKOFF_MULTIPLIER", 2.0),
		},
		RateLimit: RateLimitConfig{
			Enabled:   getEnvAsBool("RATE_LIMIT_ENABLED", true),
			KeyPrefix: getEnv("RATE_LIMIT_REDIS_KEY_PREFIX", "ratelimit:"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		S3: S3Config{
			AccessKey: getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			Region:    getEnv("AWS_REGION", ""),
			Bucket:    getEnv("AWS_BUCKET_NAME", ""),
		},
	}

	return cfg, nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

type S3Config struct {
	AccessKey string
	SecretKey string
	Region    string
	Bucket    string
}

func (c S3Config) Enabled() bool {
	return c.AccessKey != "" && c.SecretKey != "" && c.Region != "" && c.Bucket != ""
}
