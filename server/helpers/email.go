package helpers

import (
	"context"
	"log"
	"matchme-server/internal"
	"regexp"
	"strings"
)

func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	emailRe := regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)
	return emailRe.MatchString(email)
}

func IsUniqEmail(email string) bool {
	const q = `SELECT NOT EXISTS (SELECT 1 FROM users WHERE email = $1);`
	var unique bool
	err := internal.DB.QueryRow(context.Background(), q, email).Scan(&unique)
	if err != nil {
		log.Println("db error:", err)
		return false
	}
	return unique
}