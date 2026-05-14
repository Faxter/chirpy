package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func MakeJWT(userID uuid.UUID, serverSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})
	signedString, err := token.SignedString([]byte(serverSecret))
	if err != nil {
		return "", err
	}
	return signedString, nil
}

func ValidateJWT(tokenString, serverSecret string) (uuid.UUID, error) {
	type MyCustomClaims struct {
		jwt.RegisteredClaims
	}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&MyCustomClaims{},
		func(token *jwt.Token) (any, error) { return []byte(serverSecret), nil })
	if err != nil {
		return uuid.UUID{}, err
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok {
		id, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.UUID{}, err
		}
		return id, nil
	} else {
		return uuid.UUID{}, fmt.Errorf("cannot process claims")
	}
}

func GetBearerToken(header http.Header) (string, error) {
	auth := header.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("could not find authorization value in header")
	}
	return strings.TrimPrefix(auth, "Bearer "), nil
}
