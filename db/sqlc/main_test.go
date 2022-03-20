//  *@createTime    2022/3/20 0:12
//  *@author        hay&object
//  *@version       v1.0.0

package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDrive  = "postgres"
	dbSource = "postgresql://root:secret@localhost/small_bank?sslmode=disable"
)

var testQueries *Queries

var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDrive, dbSource)
	if err != nil {
		log.Fatal("connect db err")
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
