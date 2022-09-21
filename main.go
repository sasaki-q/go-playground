package main

import (
	"database/sql"
	db "dbapp/db/sqlc"
	"dbapp/src"
	"log"

	_ "github.com/lib/pq"
)

const (
	address  = "0.0.0.0:80"
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:postgres@localhost:5433/bank?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("ERROR cannot connect db === ", err)
	}

	store := db.NewStore(conn)
	server := src.NewServer(store)

	serverErr := server.StartServer(address)

	if serverErr != nil {
		log.Fatal("ERROR cannot start server === ", serverErr)
	}
}
