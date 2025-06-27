package appcore_config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type VaultSecrets struct {
	PsqlConn       string `json:"psql_conn"`
	RedisUrl       string `json:"redis_url"`
	RedisPass      string `json:"redis_pass"`
	RabbitmqUrl    string `json:"rabbitmq_url"`
	MinioAccessKey string `json:"minio_access_key"`
	MinioSecretKey string `json:"minio_secret_key"`
}

type VaultSecretResponse struct {
	Data struct {
		Data VaultSecrets `json:"data"`
	} `json:"data"`
}

var Config *Configuration

type Configuration struct {
	// Hashicorp Vault Agent
	VaultAgentAddr        string
	VaultSecretDotenvPath string
	VaultTokenSet         bool

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

	viper.SetDefault("VAULT_AGENT_ADDR", "127.0.0.1:8200")
	viper.SetDefault("VAULT_SECRET_DOTENV_PATH", "secret/data/dotenv")

	Config = &Configuration{
		VaultAgentAddr:        viper.GetString("VAULT_AGENT_ADDR"),
		VaultSecretDotenvPath: viper.GetString("VAULT_SECRET_DOTENV_PATH"),
		VaultTokenSet:         false,
	}
}

// InitSecretConfig จะถูกเปลี่ยนไปเรียก Vault Agent HTTP แทน
func InitSecretConfig(vaultAgentAddr, secretPath string) {
	secrets, err := LoadAppSecretsFromVaultAgent(vaultAgentAddr, secretPath)
	if err != nil {
		log.Fatalf("Failed to read secrets from Vault Agent: %v", err)
	}

	// คัดลอกค่าที่ได้มาใส่ใน Config global variable
	Config.PsqlConn = secrets.PsqlConn
	Config.RedisUrl = secrets.RedisUrl
	Config.RedisPass = secrets.RedisPass
	Config.RabbitmqUrl = secrets.RabbitmqUrl
	Config.MinioAccessKey = secrets.MinioAccessKey
	Config.MinioSecretKey = secrets.MinioSecretKey
	Config.VaultTokenSet = true // ตั้งค่าเป็น true เมื่อโหลดความลับได้สำเร็จ

	log.Printf("Successfully loaded secrets from Vault Agent. DB Host (example): %s", secrets.PsqlConn) // ใช้ PsqlConn เป็นตัวอย่าง
}

func LoadAppSecretsFromVaultAgent(vaultAgentAddr, secretPath string) (*VaultSecrets, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf("http://%s/v1/%s", vaultAgentAddr, secretPath) // <--- Path ต้องเป็น /v1/secret/data/dotenv
	log.Printf("Requesting secrets from Vault Agent at %s", url)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request to Vault Agent: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vault Agent returned non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body from Vault Agent: %w", err)
	}

	var vaultResp VaultSecretResponse
	err = json.Unmarshal(body, &vaultResp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal Vault Agent response: %w", err)
	}

	return &vaultResp.Data.Data, nil
}

func InitializeSecretsWithRetry() {
	maxRetries := 30
	retryInterval := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting to load secrets from Vault Agent (attempt %d/%d)...", i+1, maxRetries)

		// ใช้ closure เพื่อจับ panic จาก InitSecretConfig (ถ้ามี)
		err := func() (initErr error) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("InitSecretConfig panicked: %v", r)
					initErr = fmt.Errorf("InitSecretConfig panicked") // ส่งคืน error เพื่อให้ loop รู้ว่ามีปัญหา
				}
			}()
			InitSecretConfig(Config.VaultAgentAddr, Config.VaultSecretDotenvPath)
			return nil
		}()

		if err == nil && Config.VaultTokenSet {
			log.Println("Successfully loaded initial app secrets from Vault Agent.")
			return // โหลดสำเร็จ ออกจากฟังก์ชัน
		}

		log.Printf("Failed to load initial app secrets from Vault Agent. Retrying in %v...", retryInterval)
		time.Sleep(retryInterval)

		if i == maxRetries-1 {
			log.Fatalf("Failed to load app secrets from Vault Agent after %d retries. Exiting.", maxRetries)
		}
	}
}
