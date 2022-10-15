package jwt

import (
	"go-mongo-auth/internal/config"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type CustomClaims struct {
	Entity any
	jwt.RegisteredClaims
}

func CreateToken(o any, expiry time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, createClaims(o, expiry))

	ss, err := token.SignedString([]byte(viper.GetString("jwt.signing-key")))

	if err != nil {
		log.Println(err)
		return "", err
	}
	return ss, err
}

func createClaims(o any, expiry time.Duration) CustomClaims {
	return CustomClaims{
		o,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    config.Get("app.name"),
		},
	}
}

func ValidateToken(tokenString string) (any, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(config.Get("jwt.signing-key")), nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.Entity, nil
	}
	return nil, err
}
