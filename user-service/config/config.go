package config

import "os"

type ServiceConfig struct {
	ServerPort       string
	CassandraCluster string
}

func LoadConfig() *ServiceConfig {
	return &ServiceConfig{
		ServerPort:       getEnvOrDefault("SERVICE_PORT", ":5051"),
		CassandraCluster: getEnvOrDefault("CASSANDRA_CLUSTER", "localhost:9042,localhost:9043,localhost:9044")}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
