package utils

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail checks if email format is valid
func ValidateEmail(email string) bool {
	email = strings.TrimSpace(strings.ToLower(email))
	return emailRegex.MatchString(email)
}

// ValidatePassword checks if password meets requirements
func ValidatePassword(password string) bool {
	password = strings.TrimSpace(password)
	return len(password) >= 8
}

// ValidateRequired checks if a string field is not empty
func ValidateRequired(field, fieldName string) error {
	if strings.TrimSpace(field) == "" {
		return &AppError{
			Message:    fieldName + " is required",
			StatusCode: 400,
		}
	}
	return nil
}

