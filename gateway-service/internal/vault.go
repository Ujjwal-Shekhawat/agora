package internal

import (
	"context"
	"gateway_service/config"
	"log"

	"github.com/hashicorp/vault/api"
)

func GetCredentials() map[string]interface{} {
	cfg := config.LoadConfig()
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalf("Failed to create Vault client: %v", err)
	}

	vault_hostname := cfg.VaultHost

	client.SetAddress("http://" + vault_hostname)

	client.SetToken(cfg.VaultToken)

	secret, err := client.KVv2(cfg.KVStore).Get(context.Background(), cfg.KVPath)
	if err != nil {
		log.Fatalf("Failed to read secret: %v", err)
	}

	return secret.Data
}
