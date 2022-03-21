//  *@createTime    2022/3/21 17:40
//  *@author        hay&object
//  *@version       v1.0.0

package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/o1egl/paseto/v2"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

type PasetoGenerator struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func (p *PasetoGenerator) GenerateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return p.paseto.Encrypt(p.symmetricKey, payload, nil)
}

func (p *PasetoGenerator) VerifyToken(token string) (*Payload, error) {
	var payload = &Payload{}
	err := p.paseto.Decrypt(token, p.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	err = payload.Valid(jwt.NewValidationHelper())
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func NewPasetoGenerator(symmetricKey string) (Generator, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key:must be exactly %d characters", chacha20poly1305.KeySize)
	}
	return &PasetoGenerator{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}
