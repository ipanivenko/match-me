package helpers

import (
	"strings"

	"github.com/google/uuid"
)

func IsValidID(id string) bool {
	_, err := uuid.Parse(strings.TrimSpace(id))
	return err == nil
}