package config

import "os"

type ServiceConfig struct {
	ServerPort string
}

func LoadConfig() *ServiceConfig {
	return &ServiceConfig{ServerPort: getEnvOrDefault("USER_SERVICE_PORT", "localhost:5051")}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
