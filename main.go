package main

import (
	"database/sql"
	"log"
	_"github.com/lib/pq"
	"github.com/AlhajiTaibu/simplebank/api"
	db "github.com/AlhajiTaibu/simplebank/sqlc"
)

const (
	dbDriver = "postgres"
	source   = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	address  = "0.0.0.0:8000"
)

func main() {

	conn, err := sql.Open(dbDriver, source)

	if err != nil {
		log.Fatal("cannot connect")
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	errs := server.Start(address)

	if errs != nil {
		log.Fatal("cannot connect to server", errs)
	}

}
