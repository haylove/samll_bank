//  *@createTime    2022/3/20 0:12
//  *@author        hay&object
//  *@version       v1.0.0

package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/haylove/small_bank/util"
	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config,", err)
	}
	testDB, err = sql.Open(config.DBDrive, config.DBSource)
	if err != nil {
		log.Fatal("connect db err")
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
