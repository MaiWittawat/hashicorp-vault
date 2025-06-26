package api

import (
	"hashicorp-vault/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUserAPI(router *gin.Engine, userHandler *handler.UserHandler) {
	public := router.Group("/users")
	public.GET("/login", userHandler.Login)
	public.GET("/env/:service", userHandler.GetSecret)
	public.POST("/env/:service", userHandler.PutSecret)
	public.PATCH("env/:service", userHandler.PatchSecret)
}