package util

import "regexp"

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[\w-\.+]+@([\w-]+\.)+[\w-]{2,4}$`)
	return re.MatchString(email)
}
