package config

import (
	"os"
)

type Config struct {
	WeatherApiKey string
	ServerPort    string
}

func Load() *Config {
	cfg := &Config{}

	cfg.WeatherApiKey = getEnv("WEATHER_API_KEY", "")
	cfg.ServerPort = getEnv("SERVER_PORT", "8080")

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
