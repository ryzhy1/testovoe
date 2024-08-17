package services

import (
	"awesomeProject/internal/models"
	"errors"
)

type MockTokenStorage struct {
	tokens map[string]*models.RefreshToken
}

func NewMockTokenStorage() *MockTokenStorage {
	return &MockTokenStorage{
		tokens: make(map[string]*models.RefreshToken),
	}
}

func (m *MockTokenStorage) SaveRefreshToken(token *models.RefreshToken) error {
	m.tokens[token.UserID] = token
	return nil
}

func (m *MockTokenStorage) GetRefreshTokenByUserID(userID string) (*models.RefreshToken, error) {
	token, ok := m.tokens[userID]
	if !ok {
		return nil, errors.New("token not found")
	}
	return token, nil
}

func (m *MockTokenStorage) DeleteRefreshToken(userID string) error {
	delete(m.tokens, userID)
	return nil
}
