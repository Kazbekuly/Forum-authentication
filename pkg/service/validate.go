package service

import (
	"errors"
	"net/mail"
	"unicode"
)

func checkEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("Email is not valid")
	}
	return nil
}

func checkPassword(password string) bool {
	var (
		minLen     = false
		maxLen     = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(password) >= 7 {
		minLen = true
	}
	if len(password) <= 30 {
		maxLen = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return minLen && maxLen && hasLower && hasNumber && hasSpecial && hasUpper
}

func checkUsername(username string) error {
	for _, char := range username {
		if char < 32 || char > 126 {
			return errors.New("Username is not valid")
		}
	}
	return nil
}
