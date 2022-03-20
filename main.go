//  *@createTime    2022/3/20 17:46
//  *@author        hay&object
//  *@version       v1.0.0

package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/haylove/small_bank/api"
	db "github.com/haylove/small_bank/db/sqlc"
	"github.com/haylove/small_bank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("loadConfig err,", err)
	}

	conn, err := sql.Open(config.DBDrive, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db,", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)
	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("cannot start server,", err)
	}
}
