package utils

import "strings"

func GetEmailDomain(email string) string {
	arr := strings.Split(email, "@")
	if len(arr) < 2 {
		return ""
	}

	return arr[1]
}
