package token_handlers

import (
	"awesomeProject/internal/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TokenHandlers struct {
	tokenService *services.TokenService
}

func NewTokenHandlers(ts *services.TokenService) *TokenHandlers {
	return &TokenHandlers{
		tokenService: ts,
	}
}

func (th *TokenHandlers) CreateTokens(c *gin.Context) {
	userID := c.Query("user_id")
	ip := c.ClientIP()

	if userID == "" || ip == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id is required",
		})
		return
	}

	tokenPair, err := th.tokenService.CreateTokens(userID, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("can't create token pair: %w", err),
		})
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

func (th *TokenHandlers) RefreshToken(c *gin.Context) {
	refreshToken := c.Query("refresh_token")
	userID := c.Query("user_id")
	ip := c.ClientIP()

	if refreshToken == "" || userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh_token and user_id are required",
		})
		return
	}

	// Попытка сохранить новый refresh token
	if err := th.tokenService.StoreRefreshToken(userID, refreshToken, ip); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store new refresh token: ",
		})
		return
	}

	// Создание новой пары токенов
	tokenPair, err := th.tokenService.CreateTokens(userID, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("can't create token pair: %w", err),
		})
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}
