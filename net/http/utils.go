package http

import "strings"

// includeProtocol returns trues if string include protocol
func includeProtocol(url string) bool {
	if strings.Contains(url, "://") {
		return true
	}
	return false
}
