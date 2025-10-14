package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	FilePath   string
	ConsoleLog bool
}

// При ошибке возвращаются дефолтные значения
func Load() *Config {
	_ = godotenv.Load(".env")
	newConf := &Config{}
	newConf.Port = getEnv("SERVER_PORT", "8080")
	newConf.FilePath = getEnv("LOGGER_FILE_PATH", "file.log")

	b, _ := strconv.ParseBool(getEnv("LOGGER_CONSOLE", "false"))

	newConf.ConsoleLog = b

	return newConf
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
