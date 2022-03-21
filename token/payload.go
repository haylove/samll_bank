//  *@createTime    2022/3/21 16:16
//  *@author        hay&object
//  *@version       v1.0.0

package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

//Payload is the payload of token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

//Valid verify if the payload is invalid or not
func (p *Payload) Valid(helper *jwt.ValidationHelper) error {
	if helper.After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

//NewPayload create a new payload
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        id,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}
