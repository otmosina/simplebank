package main

import (
	"database/sql"
	"log"

	"github.com/otmosina/simplebank/api"
	db "github.com/otmosina/simplebank/db/sqlc"
	"github.com/otmosina/simplebank/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can not load config file", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can not connect to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can not connect to server", err)
	}
}
