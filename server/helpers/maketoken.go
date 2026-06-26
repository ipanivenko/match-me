package helpers

import (
	"matchme-server/internal"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func MakeAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(internal.Cfg.JWTSecret))
}
