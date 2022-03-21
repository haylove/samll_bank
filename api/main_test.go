//  *@createTime    2022/3/21 1:27
//  *@author        hay&object
//  *@version       v1.0.0

package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/haylove/small_bank/db/sqlc"
	"github.com/haylove/small_bank/util"
	"github.com/stretchr/testify/require"
)

func createRandomServer(t *testing.T, store db.Store) (*Server, error) {
	config := util.Config{
		TokenSymmetricKey:   util.RandString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server, nil
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
