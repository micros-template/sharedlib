package utils

import "regexp"

// change to https://github.com/badoux/checkmail if need to make sure valid format, valid host, and valid domain
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
