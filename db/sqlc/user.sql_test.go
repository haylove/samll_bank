//  *@createTime    2022/3/21 4:25
//  *@author        hay&object
//  *@version       v1.0.0

package db

import (
	"context"
	"testing"
	"time"

	"github.com/haylove/small_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	hashPassword, err := util.HashPassword(util.RandString(6))
	require.NoError(t, err)
	arg := CreateUserParams{
		Username:       util.RandOwner(),
		HashedPassword: hashPassword,
		FullName:       util.RandOwner(),
		Email:          util.RandEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, user)
	require.NotZero(t, user.CreatedAt)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	return user
}

func TestQueries_CreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestQueries_GetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
