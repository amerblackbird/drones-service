package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	NATS     NATSConfig     `json:"nats"`
	Jwt      JwtConfig      `json:"auth"`
}

// JwtConfig holds JWT configuration
type JwtConfig struct {
	Secret    string `json:"secret"`
	ExpiresIn string `json:"expires_in"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	ApiPrefix string `json:"api_prefix"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSLMode  string `json:"ssl_mode"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// NATSConfig holds NATS configuration
type NATSConfig struct {
	URL        string       `json:"url"`
	Servers    []string     `json:"servers"`
	QueueGroup string       `json:"queue_group"`
	Subjects   NATSSubjects `json:"subjects"`
}

// NATSSubjects defines all NATS subjects
type NATSSubjects struct {
	OrdersEvents      string `json:"orders_events"`
	DronesEvents      string `json:"drones_events"`
	UsersEvents       string `json:"users_events"`
	LogActivityEvents string `json:"log_activity_events"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Host:      getEnv("SERVER_HOST", "localhost"),
			Port:      getEnvAsInt("SERVER_PORT", 8080),
			ApiPrefix: getEnv("API_PREFIX", "/api/v1"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "auth_service_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		NATS: NATSConfig{
			URL:        getEnv("NATS_URL", "nats://localhost:4222"),
			Servers:    []string{getEnv("NATS_SERVERS", "nats://localhost:4222")},
			QueueGroup: getEnv("NATS_QUEUE_GROUP", "drones.service"),
			Subjects: NATSSubjects{
				OrdersEvents:      getEnv("NATS_SUBJECT_ORDERS_EVENTS", "orders.events"),
				DronesEvents:      getEnv("NATS_SUBJECT_DRONES_EVENTS", "drones.events"),
				UsersEvents:       getEnv("NATS_SUBJECT_USERS_EVENTS", "users.events"),
				LogActivityEvents: getEnv("NATS_SUBJECT_LOG_ACTIVITY_EVENTS", "log_activity.events"),
			},
		},
		Jwt: JwtConfig{
			Secret:    getEnv("AUTH_JWT_SECRET", "secret"),
			ExpiresIn: getEnv("AUTH_JWT_EXPIRES_IN", "24h"),
		},
	}

	return config, nil
}

// DatabaseURL returns the database connection URL
func (c *DatabaseConfig) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

// RedisURL returns the Redis connection URL
func (c *RedisConfig) RedisURL() string {
	if c.Password != "" {
		return fmt.Sprintf("%s@%s:%d/%d", c.Password, c.Host, c.Port, c.DB)
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Helper functions
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}
