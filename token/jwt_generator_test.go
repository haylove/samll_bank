//  *@createTime    2022/3/21 17:12
//  *@author        hay&object
//  *@version       v1.0.0

package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/haylove/small_bank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTNormalToken(t *testing.T) {
	jwtGenerator, err := NewJWTGenerator(util.RandString(32))
	require.NoError(t, err)

	username := util.RandOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	token, err := jwtGenerator.GenerateToken(username, duration)
	require.NoError(t, err)

	payload, err := jwtGenerator.VerifyToken(token)
	require.NoError(t, err)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestJWTExpiredToken(t *testing.T) {
	jwtGenerator, err := NewJWTGenerator(util.RandString(32))
	require.NoError(t, err)

	username := util.RandOwner()

	token, err := jwtGenerator.GenerateToken(username, -time.Minute)
	require.NoError(t, err)

	payload, err := jwtGenerator.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestJWTInvalidToken(t *testing.T) {
	payload, err := NewPayload(util.RandOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	jwtGenerator, err := NewJWTGenerator(util.RandString(32))
	require.NoError(t, err)

	payload, err = jwtGenerator.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
