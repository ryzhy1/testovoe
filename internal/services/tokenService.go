package services

import (
	"awesomeProject/internal/lib/jwt"
	"awesomeProject/internal/middlewares"
	"awesomeProject/internal/models"
	"fmt"
	"log/slog"
	"time"
)

type TokenStorage interface {
	SaveRefreshToken(token *models.RefreshToken) error
	GetRefreshTokenByUserID(userID string) (*models.RefreshToken, error)
	DeleteRefreshToken(userID string) error
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenService struct {
	log       *slog.Logger
	db        TokenStorage
	jwtSecret string
}

func NewTokenService(log *slog.Logger, db TokenStorage) *TokenService {
	return &TokenService{
		log: log,
		db:  db,
	}
}

func (ts *TokenService) CreateTokens(userID, ip string) (*models.TokenPair, error) {
	const op = "TokenService.CreateTokens"

	ts.log.With(
		slog.String("op", op),
		slog.String("userID", userID),
		slog.String("ip", ip),
	)

	ts.log.Info("Creating tokens", "userID", userID, "ip", ip)

	tokenPair, err := jwt.NewTokens(userID, ip, time.Hour)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ts.log.Info("Tokens created", "userID", userID, "ip", ip)

	refreshTokenHash, err := middlewares.HashRefreshToken(tokenPair.RefreshToken.TokenHash)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tokens := models.RefreshToken{
		UserID:    userID,
		ClientIP:  ip,
		TokenHash: refreshTokenHash,
	}

	fmt.Println(tokenPair.RefreshToken)

	ts.log.Info("Saving tokens", "userID", userID, "ip", ip)

	err = ts.db.SaveRefreshToken(&tokens)
	if err != nil {
		return nil, err
	}

	ts.log.Info("Tokens saved", "userID", userID, "ip", ip)

	ts.log.Info("Returning tokens", "userID", userID, "ip", ip)

	return &models.TokenPair{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}

func (ts *TokenService) StoreRefreshToken(userID, refreshToken, ip string) error {
	const op = "TokenService.StoreRefreshToken"

	ts.log.With(
		slog.String("op", op),
		slog.String("userID", userID),
		slog.String("refreshToken", refreshToken),
		slog.String("ip", ip),
	)

	ts.log.Info("Store refresh token", "userID", userID, "ip", ip)

	storedToken, err := ts.db.GetRefreshTokenByUserID(userID)
	if err != nil {
		sendEmailWarning(userID, ip)
		return fmt.Errorf("could not get stored token: %w", err)
	}

	ts.log.Info("Stored token found", "userID", userID, "ip", ip)

	if storedToken.ClientIP != ip {
		sendEmailWarning(userID, ip)
		return fmt.Errorf("IP address mismatch")
	}

	ts.log.Info("IP address matched", "userID", userID, "ip", ip)

	if !middlewares.CheckHash(storedToken.TokenHash, refreshToken) {
		sendEmailWarning(userID, ip)
		return fmt.Errorf("invalid refresh token")
	}

	if err = ts.db.DeleteRefreshToken(userID); err != nil {
		sendEmailWarning(userID, ip)
		return fmt.Errorf("could not delete old refresh token: %w", err)
	}

	return nil
}

func sendEmailWarning(userID, ip string) {
	fmt.Println("send email warning", userID, ip)
}
