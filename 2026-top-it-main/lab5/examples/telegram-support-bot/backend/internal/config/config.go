package config

import (
	"os"
	"strconv"
)

type Config struct {
	// Telegram Bot
	BotToken   string
	WebhookURL string

	// Database
	DatabaseURL      string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string

	// Server
	ServerPort string
	ServerHost string

	// JWT
	JWTSecret string

	// OpenTelemetry (optional; used by otelgin and OTLP exporter)
	OTELServiceName string
}

func Load() *Config {
	return &Config{
		BotToken:         getEnv("BOT_TOKEN", ""),
		WebhookURL:       getEnv("WEBHOOK_URL", ""),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/telegram_bot?sslmode=disable"),
		DatabaseUser:     getEnv("DATABASE_USER", "postgres"),
		DatabasePassword: getEnv("DATABASE_PASSWORD", "postgres"),
		DatabaseName:     getEnv("DATABASE_NAME", "telegram_bot"),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		ServerHost:       getEnv("SERVER_HOST", "0.0.0.0"),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key"),
		OTELServiceName:  getEnv("OTEL_SERVICE_NAME", "telegram-support-backend"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) GetServerAddr() string {
	return c.ServerHost + ":" + c.ServerPort
}

func (c *Config) GetDBPort() int {
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	return port
}
