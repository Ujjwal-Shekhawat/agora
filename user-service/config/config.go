package config

import "os"

type ServiceConfig struct {
	ServerPort       string
	CassandraCluster string
	MongoURI         string
}

func LoadConfig() *ServiceConfig {
	return &ServiceConfig{
		ServerPort:       getEnvOrDefault("SERVICE_PORT", ":5051"),
		MongoURI:         getEnvOrDefault("MONGO_DB_URI", "mongodb://kamisama:root@localhost:27017/?authSource=admin"),
		CassandraCluster: getEnvOrDefault("CASSANDRA_CLUSTER", "localhost:9042,localhost:9043,localhost:9044"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
