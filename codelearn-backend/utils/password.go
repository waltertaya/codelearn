package utils

import "golang.org/x/crypto/bcrypt"

func HashingPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	return hashedPassword, err
}

func ComparePassword(hashedPassword, password []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)

	return err
}
