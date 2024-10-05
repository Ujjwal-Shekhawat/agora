package config

import "os"

type Config struct {
	UserServiceAddr string
	ServerPort      string
	VaultHost       string
	VaultToken      string
	KVStore         string
	KVPath          string
	KafkaBrokers    string
}

func LoadConfig() *Config {
	return &Config{
		UserServiceAddr: getEnvOrDefault("USER_SERVICE_ADDRESS", ":5051"),
		ServerPort:      getEnvOrDefault("SERVICE_PORT", ":8080"),
		VaultHost:       getEnvOrDefault("VAULT_HOST", "localhost:8200"),
		VaultToken:      getEnvOrDefault("VAULT_TOKEN", ""),
		KVStore:         getEnvOrDefault("KVSTORE_NAME", "signing_key"),
		KVPath:          getEnvOrDefault("KVPATH_NAME", "keys"),
		KafkaBrokers:    getEnvOrDefault("KAFKA_BROKERS", ""),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
