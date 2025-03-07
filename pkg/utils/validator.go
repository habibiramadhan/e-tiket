//pkg/utils/validator.go

package utils

import (
	"errors"
	"regexp"
)

func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return errors.New("format email tidak valid")
	}
	return nil
}

func ValidateUsername(username string) error {
	if len(username) < 4 {
		return errors.New("username harus minimal 4 karakter")
	}
	
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !re.MatchString(username) {
		return errors.New("username hanya boleh berisi huruf, angka, dan underscore")
	}
	
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("password harus minimal 6 karakter")
	}
	return nil
}