package auth

import (
	"errors"
	"net/http"
	"strings"
)

/*
Extracts an api key from the headers from HTTP request.

Example:
Authorization: ApiKey {insert api key here}
*/
func GetApiKey(headers http.Header) (apiKey string, err error) {
	head := headers.Get("Authorization")
	
	if strings.TrimSpace(head) == "" {
		return "", errors.New("no auth info found") 
	}

	heads := strings.Split(head, " ")
	if len(heads) != 2 {
		return "", errors.New("malformed auth header")
	}

	if heads[0] != "ApiKey" {
		return "", errors.New("malformed auth first part header")
	}

	return heads[1], nil

}