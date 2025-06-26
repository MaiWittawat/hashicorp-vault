package main

import (
	appcore_config "hashicorp-vault/config"
	"log"

	vault "github.com/hashicorp/vault/api"
)

// var (
// 	clientVault *vault.Client
// )

func test() {
	// init config
	appcore_config.InitConfiguration()
	config := vault.DefaultConfig()

	// define variable
	var err error
	config.Address = appcore_config.Config.VaultAddr

	// new client conenction
	clientVault, err = vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	// set token for authenticate
	clientVault.SetToken("wittawat")

	// Put all config to vault
	grouped := PutAllConfig(appcore_config.Config)

	// Get all secret from vault
	DisplayAllSecretFromVault(grouped)	
}
