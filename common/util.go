package common

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPlaintext(plaintext string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	return string(bytes), err
}
