package config

import "os"

type Config struct {
	UserServiceAddr string
	ServerPort      string
}

func LoadConfig() *Config {
	return &Config{
		UserServiceAddr: getEnvOrDefault("USER_SERVICE_ADDRESS", "localhost:5051"),
		ServerPort:      getEnvOrDefault("SERVICE_PORT", "localhost:8080"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
