package utils

import (
	"github.com/google/uuid"
)

func ValidateUUID(value string) bool {
	_, err := uuid.Parse(value)
	return err == nil
}
