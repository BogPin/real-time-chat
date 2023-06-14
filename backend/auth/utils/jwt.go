package utils

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTStrategy struct {
	jwtSecret     string
	signingMethod jwt.SigningMethod
}

func NewJWTStrategy(secret string, method jwt.SigningMethod) JWTStrategy {
	return JWTStrategy{secret, method}
}

type payload struct {
	UserId int `json:"userId"`
}

type userClaims struct {
	*jwt.StandardClaims
	Payload payload `json:"payload"`
}

func (s *JWTStrategy) CreateJWT(id int) (string, HttpError) {
	token := jwt.New(s.signingMethod)
	token.Claims = userClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		payload{id},
	}
	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", NewHttpError(err, http.StatusInternalServerError)
	}
	return tokenStr, nil
}

func (s *JWTStrategy) DecodeJWT(tokenStr string) (*payload, HttpError) {
	var claims userClaims
	_, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("wrong signing method")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, NewHttpError(err, http.StatusUnauthorized)
	}
	return &claims.Payload, nil
}
