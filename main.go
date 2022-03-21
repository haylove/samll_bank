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
		log.Fatal("cannot load config,", err)
	}

	conn, err := sql.Open(config.DBDrive, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db,", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server,", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server,", err)
	}
}
