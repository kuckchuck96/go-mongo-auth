package jwt

import (
	"go-mongo-auth/internal/config"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type (
	IJwtToken interface {
		CreateToken(any, time.Duration) (string, error)
		ValidateToken(string) (any, error)
	}

	JwtToken struct {
		Config config.Config
	}
)

type CustomClaims struct {
	Entity any
	jwt.RegisteredClaims
}

func NewJwtToken(config config.Config) IJwtToken {
	return &JwtToken{
		config,
	}
}

func (j *JwtToken) CreateToken(o any, expiry time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, createClaims(j, o, expiry))

	ss, err := token.SignedString([]byte(j.Config.Jwt.SigningKey))

	if err != nil {
		log.Println(err)
		return "", err
	}
	return ss, err
}

func createClaims(j *JwtToken, o any, expiry time.Duration) CustomClaims {
	return CustomClaims{
		o,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.Config.App.Name,
		},
	}
}

func (j *JwtToken) ValidateToken(tokenString string) (any, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(j.Config.Jwt.SigningKey), nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.Entity, nil
	}
	return nil, err
}
