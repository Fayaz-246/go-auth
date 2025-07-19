package utils

import "regexp"

func CheckEmail(inp string) bool {
	reg := regexp.MustCompile(`^[\w\-\.]+@([\w\-]+\.)+[\w\-]{2,}$`)
	return reg.MatchString(inp)
}
