package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateJWT(t *testing.T) {
	_, err := MakeJWT(uuid.New(), "supersecret", 1*time.Second)
	if err != nil {
		t.Errorf("could not create JWT! %s", err)
	}
}

func TestValidateJWT(t *testing.T) {
	userId := uuid.New()
	jwt, _ := MakeJWT(userId, "supersecret", 1*time.Second)
	validatedId, err := ValidateJWT(jwt, "supersecret")
	if err != nil {
		t.Errorf("validation of JWT failed: %s", err)
	}
	if userId != validatedId {
		t.Errorf("id mismatch: %s vs. %s", userId, validatedId)
	}
}

func TestValidateExpiredJWT(t *testing.T) {
	userId := uuid.New()
	jwt, _ := MakeJWT(userId, "supersecret", 1*time.Microsecond)

	time.Sleep(10 * time.Microsecond)

	_, err := ValidateJWT(jwt, "supersecret")
	if err == nil {
		t.Errorf("expected token to be expired!")
	}
}

func TestValidateSignedJWTWithDifferentKey(t *testing.T) {
	userId := uuid.New()
	jwt, _ := MakeJWT(userId, "supersecret", 1*time.Second)
	_, err := ValidateJWT(jwt, "differentKey")
	if err == nil {
		t.Errorf("expected validation to fail due to different sign key!")
	}
}

func TestGetBearerToken(t *testing.T) {
	header := http.Header{}
	header.Set("Authorization", "Bearer some-token")

	token, _ := GetBearerToken(header)
	if token != "some-token" {
		t.Errorf("extracted unexpected token: %s", token)
	}
}

func TestGetBearerTokenFromEmptyAuthorizationHeader(t *testing.T) {
	header := http.Header{}

	_, err := GetBearerToken(header)
	if err == nil {
		t.Errorf("expected error because of missing authorization header!")
	}
}
