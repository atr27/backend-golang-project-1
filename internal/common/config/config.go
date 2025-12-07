package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
	NATS     NATSConfig
	CORS     CORSConfig
	Security SecurityConfig
	Upload   UploadConfig
	Email    EmailConfig
	External ExternalConfig
	FHIR     FHIRConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Name        string
	Environment string
	Port        string
	Version     string
	LogLevel    string
	LogFormat   string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL            string
	Host           string
	Port           string
	User           string
	Password       string
	Name           string
	SSLMode        string
	MaxConnections int
	MaxIdleConns   int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret                string
	ExpirationHours       int
	RefreshExpirationHours int
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// NATSConfig holds NATS configuration
type NATSConfig struct {
	URL       string
	ClusterID string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// SecurityConfig holds security configuration
type SecurityConfig struct {
	EncryptionKey          string
	MFAIssuer              string
	SessionTimeoutMinutes  int
	DataEncryptionEnabled  bool
	AuditLogRetentionYears int
	RateLimitPerMinute     int
}

// UploadConfig holds file upload configuration
type UploadConfig struct {
	MaxSizeMB  int
	UploadPath string
}

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	EmailFrom    string
}

// ExternalConfig holds external system configuration
type ExternalConfig struct {
	ERPAPIUrl string
	ERPAPIKey string
	LISAPIUrl string
	LISAPIKey string
	RISAPIUrl string
	RISAPIKey string
}

// FHIRConfig holds FHIR server configuration
type FHIRConfig struct {
	ServerURL string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	config := &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "Hospital-EMR-System"),
			Environment: getEnv("APP_ENV", "development"),
			Port:        getEnv("APP_PORT", "8080"),
			Version:     getEnv("APP_VERSION", "1.0.0"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
			LogFormat:   getEnv("LOG_FORMAT", "json"),
		},
		Database: DatabaseConfig{
			URL:            getEnv("DATABASE_URL", ""),
			Host:           getEnv("DB_HOST", "localhost"),
			Port:           getEnv("DB_PORT", "5432"),
			User:           getEnv("DB_USER", "postgres"),
			Password:       getEnv("DB_PASSWORD", ""),
			Name:           getEnv("DB_NAME", "emr_database"),
			SSLMode:        getEnv("DB_SSLMODE", "disable"),
			MaxConnections: getEnvAsInt("DB_MAX_CONNECTIONS", 100),
			MaxIdleConns:   getEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 10),
		},
		JWT: JWTConfig{
			Secret:                 getEnv("JWT_SECRET", "your_jwt_secret_key"),
			ExpirationHours:        getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
			RefreshExpirationHours: getEnvAsInt("JWT_REFRESH_EXPIRATION_HOURS", 168),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		NATS: NATSConfig{
			URL:       getEnv("NATS_URL", ""),
			ClusterID: getEnv("NATS_CLUSTER_ID", "emr-cluster"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "https://frontend-hospital-gules.vercel.app"}),
			AllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}),
			AllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Accept", "Authorization"}),
		},
		Security: SecurityConfig{
			EncryptionKey:          getEnv("ENCRYPTION_KEY", ""),
			MFAIssuer:              getEnv("MFA_ISSUER", "Hospital-EMR"),
			SessionTimeoutMinutes:  getEnvAsInt("SESSION_TIMEOUT_MINUTES", 30),
			DataEncryptionEnabled:  getEnvAsBool("DATA_ENCRYPTION_ENABLED", true),
			AuditLogRetentionYears: getEnvAsInt("AUDIT_LOG_RETENTION_YEARS", 25),
			RateLimitPerMinute:     getEnvAsInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 100),
		},
		Upload: UploadConfig{
			MaxSizeMB:  getEnvAsInt("MAX_UPLOAD_SIZE_MB", 50),
			UploadPath: getEnv("UPLOAD_PATH", "./uploads"),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnv("SMTP_PORT", "587"),
			SMTPUser:     getEnv("SMTP_USER", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			EmailFrom:    getEnv("EMAIL_FROM", "noreply@hospital-emr.com"),
		},
		External: ExternalConfig{
			ERPAPIUrl: getEnv("ERP_API_URL", ""),
			ERPAPIKey: getEnv("ERP_API_KEY", ""),
			LISAPIUrl: getEnv("LIS_API_URL", ""),
			LISAPIKey: getEnv("LIS_API_KEY", ""),
			RISAPIUrl: getEnv("RIS_API_URL", ""),
			RISAPIKey: getEnv("RIS_API_KEY", ""),
		},
		FHIR: FHIRConfig{
			ServerURL: getEnv("FHIR_SERVER_URL", ""),
		},
	}

	// Validate critical configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.URL == "" && c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required when DATABASE_URL is not set")
	}

	if c.JWT.Secret == "" || c.JWT.Secret == "your_jwt_secret_key" {
		return fmt.Errorf("JWT_SECRET must be set to a secure value")
	}

	if c.Security.EncryptionKey == "" && c.Security.DataEncryptionEnabled {
		return fmt.Errorf("ENCRYPTION_KEY is required when data encryption is enabled")
	}

	return nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	if c.Database.URL != "" {
		return c.Database.URL
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// GetRedisAddr returns Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

// GetJWTExpiration returns JWT expiration duration
func (c *Config) GetJWTExpiration() time.Duration {
	return time.Duration(c.JWT.ExpirationHours) * time.Hour
}

// GetJWTRefreshExpiration returns JWT refresh expiration duration
func (c *Config) GetJWTRefreshExpiration() time.Duration {
	return time.Duration(c.JWT.RefreshExpirationHours) * time.Hour
}

// IsProduction returns true if running in production
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// IsDevelopment returns true if running in development
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
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
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		var result []string
		for _, v := range splitByComma(value) {
			if trimmed := trimSpace(v); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}
	return defaultValue
}

func splitByComma(s string) []string {
	var result []string
	current := ""
	for _, char := range s {
		if char == ',' {
			result = append(result, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}
	return s[start:end]
}
