package utils

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// check password is: longer than 8, uppercase, lowercase, number and symbol
func IsStrongPassword(password string) bool {
	var (
		hasMinLen  = len(password) >= 8
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[!@#~$%^&*()_+|<>?:{}]`).MatchString(password)
	)

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func HashPasswordCompare(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10) // changed cost from 17 to 10
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
