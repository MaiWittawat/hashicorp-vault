package main

import (
	"context"
	"encoding/json"
	"fmt"
	appcore_config "hashicorp-vault/config"
	"log"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

// convert type config to json
func ConvertConfigJson(conf *appcore_config.Configuration) (map[string]interface{}, error) {
	var result map[string]interface{}
	data, err := json.Marshal(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to convert conf to json: %w", err)
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to convert conf to json: %w", err)
	}

	return result, nil
}

// Write a secret to vault
func PutSecretToVault(path string, secret map[string]interface{}) error {
	_, err := clientVault.KVv2("secret").Put(context.Background(), path, secret)
	if err != nil {
		return err
	}

	return nil
}

// Get secret from vault
func GetSecretFromVault(path string) (*vault.KVSecret, error) {
	secret, err := clientVault.KVv2("secret").Get(context.Background(), path)
	if err != nil {
		return nil, err
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no secret data : %w", err)
	}

	return secret, nil
}

func SubStringPrefix(dataMap map[string]interface{}) map[string]map[string]interface{} {
	grouped := map[string]map[string]interface{}{
		"psql":     {},
		"redis":    {},
		"rabbitmq": {},
		"minio":    {},
	}

	for key, value := range dataMap {
		for _, prefix := range []string{"psql", "redis", "rabbitmq", "minio"} {

			if strings.HasPrefix(strings.ToLower(key), strings.ToLower(prefix)) {
				grouped[prefix][key] = value
				break
			}
		}
	}
	return grouped
}

func PutAllConfig(conf *appcore_config.Configuration) map[string]map[string]interface{} {
	// convert struct to json format
	dataMap, err := ConvertConfigJson(conf)
	if err != nil {
		log.Fatal("error: ", err)
	}

	// Seperate by prefix and put it to vault
	grouped := SubStringPrefix(dataMap)

	// วน grouped แยกเพื่อส่งเข้า Vault
	for prefix, data := range grouped {
		if len(data) == 0 {
			continue
		}

		// ใช้ path ตาม prefix เช่น /dotenv/redis
		path := fmt.Sprintf("dotenv/%s", prefix)
		if err := PutSecretToVault(path, data); err != nil {
			fmt.Printf("error putting secret for %s: %v\n", prefix, err)
		} else {
			fmt.Printf("stored %d secrets at %s\n", len(data), path)
		}
	}
	return grouped
}

func DisplayAllSecretFromVault(grouped map[string]map[string]interface{}) {
	for key, _ := range grouped {
		vaultKv, err := GetSecretFromVault(fmt.Sprintf("dotenv/%s", key))
		if err != nil {
			log.Println("no secret from vault:", err)
		} else {
			fmt.Printf("secret {%s}: %v\n", key, vaultKv.Data)
		}
	}
}
