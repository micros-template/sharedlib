package utils

import "strings"

func GetFileNameExtension(filename string) string {
	parts := strings.Split(filename, ".")
	return strings.ToLower(parts[len(parts)-1])
}
