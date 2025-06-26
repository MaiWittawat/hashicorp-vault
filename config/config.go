package appcore_config

import (
	"log"

	"github.com/spf13/viper"
)

var Config *Configuration

type Configuration struct {
	// Develop or production
	Mode string `json:"mode"`

	// Hashicorp Vault
	VaultAddr        string `json:"vault_addr"`
	VaultRoleID      string `json:"vault_role_id"`
	VaultSecretID    string `json:"vault_secret_id"`
	VaultClientToken string `json:"vault_client_token"`

	// Database
	PsqlConn string `json:"psql_conn"`

	// Redis
	RedisUrl  string `json:"redis_url"`
	RedisPass string `json:"redis_pass"`

	// Message broker (rabbitmq)
	RabbitmqUrl string `json:"rabbitmq_url"`

	// Storage
	MinioAccessKey string `json:"minio_access_key"`
	MinioSecretKey string `json:"minio_secret_key"`
}

func InitConfiguration() {
	viper.AutomaticEnv()
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No .env file found, relying on system environment variables.")
		} else {
			log.Fatalf("Fatal error config file: %s \n", err)
		}
	}

	Config = &Configuration{
		Mode: viper.GetString("MODE"),

		VaultAddr:     viper.GetString("VAULT_ADDR"),
		VaultRoleID:   viper.GetString("VAULT_ROLE_ID"),
		VaultSecretID: viper.GetString("VAULT_SECRET_ID"),

		PsqlConn: viper.GetString("POSTGRES_URL"),

		RedisUrl:  viper.GetString("REDIS_URL"),
		RedisPass: viper.GetString("REDIS_PASS"),

		RabbitmqUrl: viper.GetString("RABBITMQ_URL"),

		MinioSecretKey: viper.GetString("MINIO_ROOT_USER"),
		MinioAccessKey: viper.GetString("MINIO_ROOT_PASSWORD"),
	}
}
