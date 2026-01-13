package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	TorProxy       string
	TorControlPort string

	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string

	ScraperInterval time.Duration
	MaxRetries      int
	RequestTimeOut  time.Duration

	HealthPort string

	LogLevel string
}

func Load() *Config {
	return &Config{
		TorProxy:       getEnv("TOR_PROXY", "tor:9050"),
		TorControlPort: getEnv("TOR_CONTROL_PORT", "tor:9051"),

		DBHost:     getEnv("DB_HOST", "postgres"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBName:     getEnv("DB_NAME", "cti_data"),
		DBUser:     getEnv("DB_USER", "cti_user"),
		DBPassword: getEnv("DB_PASSWORD", ""),

		ScraperInterval: getEnvAsDuration("SCRAPER_INTERVAL", 1*time.Hour),
		MaxRetries:      getEnvAsInt("MAX_RETRIES", 3),
		RequestTimeOut:  getEnvAsDuration("REQUEST_TIMEOUT", 30*time.Second),

		HealthPort: getEnv("HEALTH_PORT", "8080"),

		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}

	return fallback
}

func (c *Config) Validate() error {
	if c.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD boş olamaz")
	}
	if c.TorProxy == "" {
		return fmt.Errorf("TOR_PROXY boş olamaz")
	}
	return nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}
