package handler

import (
	appcore_config "hashicorp-vault/config"
	"hashicorp-vault/services"
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

func (h *UserHandler) Login(c *gin.Context) {
	resp, err := h.service.GetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "get token success", "token": resp})
}

func (h *UserHandler) GetSecret(c *gin.Context) {
	resp, err := h.service.GetSecret(c.Request.Context(), c.Param("service"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "get token success", "response": resp})
}

func (h *UserHandler) PutSecret(c *gin.Context) {
	var secrets map[string]interface{}

	if err := c.ShouldBindJSON(&secrets); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.PutSecret(c.Request.Context(), c.Param("service"), secrets)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "get token success", "response": resp})
}

func (h *UserHandler) PatchSecret(c *gin.Context) {
	var newSecret map[string]interface{}

	if err := c.ShouldBindJSON(&newSecret); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.PatchSecret(c.Request.Context(), c.Param("service"), newSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "patch success", "response": resp})
}