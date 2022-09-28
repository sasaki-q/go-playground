package main

import (
	"database/sql"
	db "dbapp/db/sqlc"
	"dbapp/src"
	"dbapp/utils"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("ERROR cannot load env === ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("ERROR cannot connect db === ", err)
	}

	store := db.NewStore(conn)
	server, err := src.NewServer(config, store)

	if err != nil {
		log.Fatal("ERROR cannot create server: ", err)
	}

	serverErr := server.StartServer(config.ServerAdress)

	if serverErr != nil {
		log.Fatal("ERROR cannot start server === ", serverErr)
	}
}
