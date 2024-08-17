package models

type RefreshToken struct {
	UserID    string
	TokenHash string
	ClientIP  string
}
