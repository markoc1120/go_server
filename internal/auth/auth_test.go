package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJWTValidation(t *testing.T) {
	correctSecretToken := "correctToken"
	incorrectSecretToken := "incorrectToken"
	userID := uuid.New()

	tests := []struct {
		name            string
		secretToken     string
		incorrectToken  string
		userID          uuid.UUID
		expiresIn       time.Duration
		wantMakeErr     bool
		wantValidateErr bool
	}{
		{
			name:            "Correct JWT validation",
			secretToken:     correctSecretToken,
			incorrectToken:  correctSecretToken,
			userID:          userID,
			expiresIn:       time.Duration(30 * time.Hour),
			wantMakeErr:     false,
			wantValidateErr: false,
		},
		{
			name:            "Incorrect JWT validation with wrong secretToken",
			secretToken:     correctSecretToken,
			incorrectToken:  incorrectSecretToken,
			userID:          userID,
			expiresIn:       time.Duration(30 * time.Hour),
			wantMakeErr:     false,
			wantValidateErr: true,
		},
		{
			name:            "Expired JWT token validation",
			secretToken:     correctSecretToken,
			incorrectToken:  incorrectSecretToken,
			userID:          userID,
			expiresIn:       time.Duration(0),
			wantMakeErr:     false,
			wantValidateErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.userID, tt.secretToken, tt.expiresIn)
			if (err != nil) != tt.wantMakeErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantMakeErr)
			}

			_, err = ValidateJWT(token, tt.incorrectToken)
			if (err != nil) != tt.wantValidateErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantValidateErr)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		headers       http.Header
		expectedErr   error
		expectedToken string
	}{
		{
			name:          "no header",
			headers:       http.Header{},
			expectedErr:   MissingAuthorization,
			expectedToken: "",
		},
		{
			name:          "wrong authorization in the header",
			headers:       http.Header{"Authorization": []string{"Basic TOKEN_STRING"}},
			expectedErr:   WrongAuthorization,
			expectedToken: "",
		},
		{
			name:          "empty token",
			headers:       http.Header{"Authorization": []string{"Bearer "}},
			expectedErr:   EmptyToken,
			expectedToken: "",
		},
		{
			name:          "check correct token",
			headers:       http.Header{"Authorization": []string{"Bearer TOKEN_STRING"}},
			expectedErr:   nil,
			expectedToken: "TOKEN_STRING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers)
			assertEqual(t, err, tt.expectedErr)
			assertEqual(t, token, tt.expectedToken)
		})
	}
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
