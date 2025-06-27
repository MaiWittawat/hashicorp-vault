package api

import (
	"hashicorp-vault/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUserAPI(router *gin.Engine, userHandler *handler.UserHandler) {
	public := router.Group("/users")
	public.GET("/env/:service", userHandler.GetSecret)
}