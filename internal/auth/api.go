package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if !strings.HasPrefix(auth, "ApiKey ") {
		return "", errors.New("failed to get bearer token")
	}
	return strings.TrimPrefix(auth, "ApiKey "), nil

}
