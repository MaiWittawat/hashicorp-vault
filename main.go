package main

import (
	appcore_config "hashicorp-vault/config"
	"hashicorp-vault/handler"
	"hashicorp-vault/handler/api"
	userSvc "hashicorp-vault/services/user"
	"log"

	"github.com/gin-gonic/gin"
	vault "github.com/hashicorp/vault/api"
)

var (
	clientVault *vault.Client
)

func main() {
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

	router := gin.Default()
	
	userService := userSvc.NewUserService(appcore_config.Config, clientVault)
	userHandler := handler.NewUserHandler(userService, appcore_config.Config, clientVault)
	api.RegisterUserAPI(router, userHandler)

	router.Run(":3033")
}
