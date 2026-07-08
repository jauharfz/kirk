package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", errors.New("failed to get bearer token")
	}
	return strings.TrimPrefix(auth, "Bearer "), nil
}

func MakeRefreshToken() string {
	return hex.EncodeToString([]byte(rand.Text()))
}
