package config

import "os"

type Config struct {
	CassandraCluster string
	KafkaBrokers     string
}

func LoadConfig() *Config {
	return &Config{
		CassandraCluster: getEnvOrDefault("CASSANDRA_CLUSTER", "localhost:9042,localhost:9043,localhost:9044"),
		KafkaBrokers:     getEnvOrDefault("KAFKA_BROKERS", ""),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
