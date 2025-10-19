// Package config содержит структуру конфига всего микросервиса и функцию его загрузки
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	RabbitMQ RabbitMQConfig
	Telegram TelegramConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type RabbitMQConfig struct {
	URL       string
	QueueName string
	Exchange  string
}

type TelegramConfig struct {
	Token string
}

func (d DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

type JWTConfig struct {
	Secret        string `yaml:"secret" env:"JWT_SECRET" default:"your-secret-key"`
	TokenLifespan time.Duration
}

func Load() *Config {
	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		fmt.Println("!!!\nNo .env file found, using environment variables\n!!!")
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},

		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "eventbooker"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		RabbitMQ: RabbitMQConfig{
			URL:       getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
			QueueName: getEnv("RABBITMQ_QUEUE", "event-reminders"),
			Exchange:  getEnv("RABBITMQ_EXCHANGE", "reminders-exchange"),
		},
		Telegram: TelegramConfig{
			Token: getEnv("TELEGRAM_BOT_TOKEN", ""),
		},
	}
}

// Вспомогательные функции
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultValue
	}
	val := strings.Split(valStr, sep)
	return val
}
