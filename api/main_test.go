//  *@createTime    2022/3/21 1:27
//  *@author        hay&object
//  *@version       v1.0.0

package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
