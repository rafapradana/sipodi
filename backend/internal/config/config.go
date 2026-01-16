package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	MinIO    MinIOConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Env  string
	Port string
	Name string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
	PublicURL string
}

type CORSConfig struct {
	Origins string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
			Name: getEnv("APP_NAME", "SIPODI"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "sipodi"),
			Password: getEnv("DB_PASSWORD", "sipodi_secret"),
			Name:     getEnv("DB_NAME", "sipodi"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "secret"),
			AccessExpiry:  parseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m")),
			RefreshExpiry: parseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h")),
		},
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "sipodi"),
			UseSSL:    getEnvBool("MINIO_USE_SSL", false),
			PublicURL: getEnv("MINIO_PUBLIC_URL", "http://localhost:9000"),
		},
		CORS: CORSConfig{
			Origins: getEnv("CORS_ORIGINS", "http://localhost:3000"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return b
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 15 * time.Minute
	}
	return d
}
