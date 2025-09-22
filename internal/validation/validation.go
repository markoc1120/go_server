package validation

import (
	"errors"
	"net/mail"
	"strings"
)

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	email = strings.TrimSpace(email)
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format")
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}

	if len(password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}
