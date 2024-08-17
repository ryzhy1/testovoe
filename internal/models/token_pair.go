package models

type TokenPair struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken RefreshToken `json:"refresh_token"`
}
