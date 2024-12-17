package services

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	apiErrors "github.com/innovay-software/famapp-main/app/errors"
)

// Generate a JWT access token
func GenerateJwtAccessToken(uuid string) (string, error) {
	validHours, err := strconv.Atoi(os.Getenv("JWT_ACCESSTOKEN_VALID_HOURS"))
	if err != nil {
		validHours = 24 * 60
	}
	return generateJwtToken(uuid, validHours)
}

// Generate a JWT refresh token
func GenerateJwtRefreshToken(uuid string) (string, error) {
	validHours, err := strconv.Atoi(os.Getenv("JWT_REFRESHTOKEN_VALID_HOURS"))
	if err != nil {
		validHours = 24 * 365
	}
	return generateJwtToken(uuid, validHours)
}

func generateJwtToken(sub string, validHours int) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": sub,
			"iss": os.Getenv("APP_HOME"),
			"exp": time.Now().Add(time.Duration(validHours) * time.Hour).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Verifies the token integrity with the JWT_SECRET key
func VerifyToken(tokenString string) (*jwt.Token, error) {
	jwtSecretString := os.Getenv("JWT_SECRET")
	secretKey := []byte(jwtSecretString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	})

	switch {
	case token.Valid:
		return token, nil
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return token, apiErrors.ApiErrorTokenExpired
	}
	return nil, fmt.Errorf("invalid token")
}
