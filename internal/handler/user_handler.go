package handler

import (
	appcore_config "hashicorp-vault/cmd/hashicorp-vault/config"
	"hashicorp-vault/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	vault "github.com/hashicorp/vault/api"
)

type UserHandler struct {
	service     services.UserService
	appConfig   *appcore_config.Configuration
	clientVault *vault.Client
}

func NewUserHandler(service services.UserService, appConfig *appcore_config.Configuration, clientVault *vault.Client) *UserHandler {
	return &UserHandler{
		service:     service,
		appConfig:   appConfig,
		clientVault: clientVault,
	}
}

func (h *UserHandler) GetSecret(c *gin.Context) {
	resp, err := h.service.GetSecret(c.Request.Context(), c.Param("service"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "get token success", "response": resp})
}

func (h *UserHandler) GetConfig(c *gin.Context) {
	resp, err := h.service.GetConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "get config success", "response": resp})
}


