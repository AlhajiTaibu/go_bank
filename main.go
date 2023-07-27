package main

import (
	"database/sql"
	"log"

	"github.com/AlhajiTaibu/simplebank/api"
	db "github.com/AlhajiTaibu/simplebank/sqlc"
	"github.com/AlhajiTaibu/simplebank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Unable to load config", err)
	}

	conn, err := sql.Open(config.DbDriver, config.DbSource)

	if err != nil {
		log.Fatal("cannot connect")
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	errs := server.Start(config.ServerAddress)

	if errs != nil {
		log.Fatal("cannot connect to server", errs)
	}

}
