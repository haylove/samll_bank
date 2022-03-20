//  *@createTime    2022/3/20 17:46
//  *@author        hay&object
//  *@version       v1.0.0

package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"small_bank/api"
	db "small_bank/db/sqlc"
)

const (
	dbDrive       = "postgres"
	dbSource      = "postgresql://root:secret@localhost/small_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDrive, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db,", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	if err := server.Start(serverAddress); err != nil {
		log.Fatal("cannot start server,", err)
	}
}
