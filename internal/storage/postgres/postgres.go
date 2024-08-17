package postgres

import (
	"awesomeProject/internal/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewPostgresDB(url string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) SaveRefreshToken(token *models.RefreshToken) error {
	_, err := s.db.Exec(context.Background(), `
		INSERT INTO refresh_tokens (user_id, client_ip, token_hash) 
		VALUES ($1, $2, $3)`,
		token.UserID, token.ClientIP, token.TokenHash)
	return err
}

func (s *Storage) GetRefreshTokenByUserID(userID string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := s.db.QueryRow(context.Background(), `
		SELECT user_id, token_hash, client_ip 
		FROM refresh_tokens WHERE user_id = $1`, userID).
		Scan(&token.UserID, &token.TokenHash, &token.ClientIP)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (s *Storage) DeleteRefreshToken(userID string) error {
	_, err := s.db.Exec(context.Background(), `
		DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	return err
}
