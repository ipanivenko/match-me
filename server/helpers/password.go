package helpers

import (
	"errors"
	"log"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func IsValidPassword(password string) bool {
	//just 6 symbols with at least 1 letter
	if len(password) < 6 {
		return false
	}
	
	re := regexp.MustCompile(`[A-Za-z]`)
	return re.MatchString(password)

}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("password hashing failed")
		return "", errors.New("password hashing failed")
	}
	return string(hashed), nil
}

func IsCorrectPassword(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		log.Println("wrong password")
		return false		
	}
	return true
}
