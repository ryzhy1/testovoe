package services

import (
	"log/slog"
	"testing"
)

func TestNewTokenService(t *testing.T) {
	var logger *slog.Logger
	db := NewMockTokenStorage()

	ts := NewTokenService(logger, db)
	if ts == nil {
		t.Errorf("NewTokenService() returned nil")
	}
}

func TestTokenService_CreateTokens(t *testing.T) {
	var logger *slog.Logger
	db := NewMockTokenStorage()
	ts := NewTokenService(logger, db)

	userID := "123"
	ip := "127.0.0.1"

	t.Run("success", func(t *testing.T) {
		tokenPair, err := ts.CreateTokens(userID, ip)
		if err != nil {
			t.Errorf("CreateTokens() error = %v", err)
		}
		if tokenPair == nil {
			t.Errorf("CreateTokens() returned nil")
		}
		if tokenPair.AccessToken == "" || tokenPair.RefreshToken.TokenHash == "" {
			t.Errorf("CreateTokens() returned empty tokens")
		}
	})
}

func TestTokenService_StoreRefreshToken(t *testing.T) {
	var logger *slog.Logger
	db := NewMockTokenStorage()
	ts := NewTokenService(logger, db)

	userID := "123"
	ip := "127.0.0.1"

	t.Run("success", func(t *testing.T) {
		tokenPair, err := ts.CreateTokens(userID, ip)
		if err != nil {
			t.Fatalf("CreateTokens() error = %v", err)
		}

		err = ts.StoreRefreshToken(userID, tokenPair.RefreshToken.TokenHash, ip)
		if err != nil {
			t.Errorf("StoreRefreshToken() error = %v", err)
		}
	})

	t.Run("IP mismatch", func(t *testing.T) {
		tokenPair, err := ts.CreateTokens(userID, ip)
		if err != nil {
			t.Fatalf("CreateTokens() error = %v", err)
		}

		err = ts.StoreRefreshToken(userID, tokenPair.RefreshToken.TokenHash, "192.168.1.1")
		if err == nil {
			t.Errorf("StoreRefreshToken() expected IP mismatch error, got nil")
		}
	})

	t.Run("Invalid refresh token", func(t *testing.T) {
		_, err := ts.CreateTokens(userID, ip)
		if err != nil {
			t.Fatalf("CreateTokens() error = %v", err)
		}

		err = ts.StoreRefreshToken(userID, "invalid-refresh-token", ip)
		if err == nil {
			t.Errorf("StoreRefreshToken() expected invalid token error, got nil")
		}
	})
}

func Test_sendEmailWarning(t *testing.T) {
	userID := "123"
	ip := "127.0.0.1"

	t.Run("send email", func(t *testing.T) {
		sendEmailWarning(userID, ip)
	})
}
