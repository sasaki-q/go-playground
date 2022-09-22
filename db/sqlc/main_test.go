package db

import (
	"database/sql"
	"dbapp/utils"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testQueries *Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("ERROR: cannot load env === ", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("ERROR: cannot connect db === ", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
