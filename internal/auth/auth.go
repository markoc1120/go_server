package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
)

const (
	MissingAuthorization = BearerTokenErr("no Authorization in the header")
	EmptyToken           = BearerTokenErr("empty token")
	WrongAuthorization   = BearerTokenErr("wrong format of Authorization in the header")
)

type BearerTokenErr string

func (e BearerTokenErr) Error() string {
	return string(e)
}

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	currTime := time.Now().UTC()
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    string(TokenTypeAccess),
			IssuedAt:  jwt.NewNumericDate(currTime),
			ExpiresAt: jwt.NewNumericDate(currTime.Add(expiresIn)),
			Subject:   userID.String(),
		},
	)
	return token.SignedString([]byte(tokenSecret))
}

// TODO: write more tests for ValidateJWT, create errors which can be tested in the auth_test.go
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimStruct,
		func(t *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		},
	)

	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("Invalid issuer")
	}
	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Invalid user ID: %w", err)
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", MissingAuthorization
	}
	if cutHeader, found := strings.CutPrefix(authHeader, "Bearer "); found {
		token := strings.TrimSpace(cutHeader)
		if token == "" {
			return "", EmptyToken
		}
		return token, nil
	}
	return "", WrongAuthorization
}
