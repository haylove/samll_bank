//  *@createTime    2022/3/21 16:28
//  *@author        hay&object
//  *@version       v1.0.0

package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"time"
)

const minSecretKeySize = 32

//JWTGenerator is a web JSON token generator
type JWTGenerator struct {
	secretKey string
}

func (J *JWTGenerator) GenerateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(J.secretKey))
}

func (J *JWTGenerator) VerifyToken(token string) (*Payload, error) {
	var keyFunc = func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(J.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		if errors.Is(err, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, err
	}
	return payload, nil
}

//NewJWTGenerator creates a new JWTGenerator
func NewJWTGenerator(secretKey string) (Generator, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key:must be at least %d characters", minSecretKeySize)
	}
	return &JWTGenerator{secretKey: secretKey}, nil
}
