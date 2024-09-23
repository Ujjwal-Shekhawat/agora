package config

import "os"

type ServiceConfig struct {
	ServerPort string
}

func LoadConfig() *ServiceConfig {
	return &ServiceConfig{ServerPort: getEnvOrDefault("SERVICE_PORT", ":5051")}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
