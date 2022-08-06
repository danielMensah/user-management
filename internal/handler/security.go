package handler

import (
	"golang.org/x/crypto/bcrypt"
)

func encryptPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
