//  *@createTime    2022/3/21 17:52
//  *@author        hay&object
//  *@version       v1.0.0

package token

import (
	"testing"
	"time"

	"github.com/haylove/small_bank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoNormalToken(t *testing.T) {
	generator, err := NewPasetoGenerator(util.RandString(32))
	require.NoError(t, err)

	username := util.RandOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	token, err := generator.GenerateToken(username, duration)
	require.NoError(t, err)

	payload, err := generator.VerifyToken(token)
	require.NoError(t, err)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestPasetoExpiredToken(t *testing.T) {
	generator, err := NewPasetoGenerator(util.RandString(32))
	require.NoError(t, err)

	username := util.RandOwner()

	token, err := generator.GenerateToken(username, -time.Minute)
	require.NoError(t, err)

	payload, err := generator.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
