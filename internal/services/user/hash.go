package services

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ошибка при хешировании пароля: %v", err)
		return "", err
	}

	return string(bytes), nil
}

func CheckPassword(hashPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		log.Printf("ошибка при проверке пароля")
	}

	return err
}
