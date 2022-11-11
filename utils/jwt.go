package utils

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtExpiryDuration = 4 * time.Hour
	jwtSigningKey     = []byte(os.Getenv("JWT_SECRET_KEY"))
)

func GenerateJWT(subject, audience, issuer string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["sub"] = subject
	claims["aud"] = audience
	claims["iss"] = issuer
	claims["exp"] = time.Now().Add(jwtExpiryDuration).Unix()
	claims["iat"] = time.Now().Unix()

	tokenString, err := token.SignedString(jwtSigningKey)

	if err != nil {
		log.Fatalf("Error while signing token: %v", err)
		return "", err
	}

	return tokenString, nil
}
