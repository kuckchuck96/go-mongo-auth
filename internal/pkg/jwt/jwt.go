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

	jwtToken struct {
		config config.Config
	}

	CustomClaims struct {
		Entity any
		jwt.RegisteredClaims
	}
)

func NewJwtToken(config config.Config) IJwtToken {
	return &jwtToken{
		config,
	}
}

func (j *jwtToken) CreateToken(o any, expiry time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, createClaims(j, o, expiry))

	ss, err := token.SignedString([]byte(j.config.Jwt.SigningKey))

	if err != nil {
		log.Println(err)
		return "", err
	}
	return ss, err
}

func createClaims(j *jwtToken, o any, expiry time.Duration) CustomClaims {
	return CustomClaims{
		o,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.config.App.Name,
		},
	}
}

func (j *jwtToken) ValidateToken(tokenString string) (any, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(j.config.Jwt.SigningKey), nil
	})

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.Entity, nil
	}
	return nil, err
}
