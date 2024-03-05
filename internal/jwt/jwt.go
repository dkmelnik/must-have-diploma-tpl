package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/dkmelnik/go-musthave-diploma/internal/dto"
	"github.com/dkmelnik/go-musthave-diploma/internal/models"
)

type Jwt struct {
	secret   string
	tokenExp time.Duration
}

func NewJwt(secret string, tokenExp time.Duration) *Jwt {
	return &Jwt{
		secret,
		tokenExp,
	}
}

func (j *Jwt) BuildJWTString(userID models.ModelID, jti string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, dto.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenExp)),
		},
		SUB: string(userID),
		JTI: jti,
	})

	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *Jwt) ParseToken(tokenString string) (*dto.Claims, error) {
	claims := &dto.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(j.secret), nil
		})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}
