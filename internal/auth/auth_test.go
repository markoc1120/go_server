package auth

import (
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
