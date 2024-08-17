package jwt

import (
	"awesomeProject/internal/models"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)
import "time"

func NewTokens(userID, ip string, tokenTTL time.Duration) (*models.TokenPair, error) {
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["ip"] = ip
	claims["exp"] = time.Now().Add(tokenTTL).Unix()

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, fmt.Errorf("can't sign token: %w", err)
	}

	guid, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("can't create uuid: %w", err)
	}

	refreshToken := models.RefreshToken{
		ClientIP:  ip,
		TokenHash: guid.String(),
		UserID:    userID,
	}

	return &models.TokenPair{
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
	}, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ValidateRefreshToken(refreshToken string, tokenHash uuid.UUID) bool {
	return refreshToken == tokenHash.String()
}
