package middlewares

import (
	"golang.org/x/crypto/bcrypt"
)

func HashRefreshToken(password string) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return string(passHash), nil
}
