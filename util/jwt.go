package util

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Claims struct {
	UserID int64
	Email  string
	jwt.RegisteredClaims
}

var jwtSecret []byte

func loadJWTSecret() []byte {
	if len(jwtSecret) == 0 {
		_ = godotenv.Load()

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			log.Fatal("[JWT] Missing JWT_SECRET in environment")
		}
		jwtSecret = []byte(secret)
	}
	return jwtSecret
}

func GenerateJWT(userID int64, email string) (string, error) {
	expirationTime := time.Now().Add(30 * time.Minute)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := loadJWTSecret()
	signedToken, err := token.SignedString(secret)
	if err != nil {
		log.Printf("[GenerateJWT][ERROR] Failed to generate token for userID %d: %v", userID, err)
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	secret := loadJWTSecret()

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expired")
		}
		return nil, errors.New("invalid token")
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
