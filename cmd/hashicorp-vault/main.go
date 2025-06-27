package main

import (
	appcore_config "hashicorp-vault/cmd/hashicorp-vault/config"
	"hashicorp-vault/internal/handler"
	"hashicorp-vault/internal/handler/api"
	userSvc "hashicorp-vault/internal/services/user"

	"github.com/gin-gonic/gin"
	vault "github.com/hashicorp/vault/api"
)

var (
	clientVault *vault.Client
)

func main() {
	// init config
	appcore_config.InitConfiguration()
	appcore_config.InitializeSecretsWithRetry()

	router := gin.Default()

	userService := userSvc.NewUserService(appcore_config.Config, clientVault)
	userHandler := handler.NewUserHandler(userService, appcore_config.Config, clientVault)
	api.RegisterUserAPI(router, userHandler)

	router.Run(":3033")
}



