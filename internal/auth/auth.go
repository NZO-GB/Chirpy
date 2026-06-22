package auth

import (
	argon2id "github.com/alexedwards/argon2id"
	uuid "github.com/google/uuid"
	jwt "github.com/golang-jwt/jwt/v5"
	"time"
	"net/http"
	"fmt"
	"strings"
	"crypto/rand"
	"encoding/hex"
)

func HashPassword(password string) (string, error) {
	hashed, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hashed, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	jwtTimeNow := jwt.NewNumericDate(time.Now())
	jwtTimeExpire := jwt.NewNumericDate(time.Now().Add(expiresIn))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:		"chirpy-access",
		IssuedAt: 	jwtTimeNow,
		ExpiresAt:	jwtTimeExpire,
		Subject:	userID.String(),
	})

	key := []byte(tokenSecret)

	tokenString, err := token.SignedString(key)

	return tokenString, err
}

func ValidateJWT(tokenString string, tokenSecret string) (uuid.UUID, error) {

	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims,
	func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("No authorization header")
	}
	authHeaderSplit := strings.Split(authHeader, " ")

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