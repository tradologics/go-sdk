package http

import "strings"

func includeProtocol(url string) bool {
	if strings.Contains(url, "://") {
		return true
	}
	return false
}
