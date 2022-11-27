package util

import (
	"net/http"
	"strings"
)
func GetSenderURL(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

// Rules: If localhost, call it true.
// Otherwise compare the pre-port part of the url to see if they match.
func IsOwner(ownerURL string, parsedURL string) bool {
	//aspects of this function may be wrong due to IPv6.
	if strings.Contains(parsedURL, "[::1]") {
		return true
	}
	ownerURL = strings.Split(ownerURL, ":")[0]
	parsedURL = strings.Split(parsedURL, ":")[0]
	if ownerURL == "localhost" || ownerURL == "[::1]" {
		if parsedURL == "localhost" || parsedURL == "[::1]" {
			return true
		}
	}
	return ownerURL == parsedURL
}
