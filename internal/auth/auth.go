package auth

import (
	"net/http"
	"fmt"
	"strings"
	"crypto/rand"
	"encoding/hex"
)

func GetBearerToken(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("No authorization header")
	}
	authHeaderSplit := strings.Split(authHeader, " ")

	if authHeaderSplit[0] != "Bearer" {
		return "", fmt.Errorf("Not a bearer token")
	}

	return authHeaderSplit[1], nil
}

func MakeRefreshToken() string {

	key := make([]byte, 32)
	bytes, err := rand.Read(key)
	if err != nil || bytes != 32 {
		panic(err)
	}
	stringKey := hex.EncodeToString(key)

	return stringKey
}

func GetAPIKey(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("No authorization header")
	}
	authHeaderSplit := strings.Split(authHeader, " ")

	if authHeaderSplit[0] != "ApiKey" {
		return "", fmt.Errorf("Not an API Key")
	}

	return authHeaderSplit[1], nil
}