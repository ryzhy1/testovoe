package routes

import (
	"awesomeProject/internal/handlers/token_handlers"
	"awesomeProject/internal/services"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, tokenService *services.TokenService) {

	handlers := token_handlers.NewTokenHandlers(tokenService)

	r.POST("/token", handlers.CreateTokens)
	r.POST("/refresh", handlers.RefreshToken)
}
