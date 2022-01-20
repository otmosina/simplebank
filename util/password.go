package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Error when hash password %s", password)
	}
	return string(hashedPassword), nil
}

func CheckPassword(password string, inputHashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(inputHashedPassword), []byte(password))
	// hashedPassword, err := HashPassword(password)
	// if err != nil {
	// 	return err
	// }
	// if inputHashedPassword != hashedPassword {
	// 	return fmt.Errorf("Error when check password %s", password)
	// }
	// return nil
}
