package middlewares

import "golang.org/x/crypto/bcrypt"

func CheckHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
