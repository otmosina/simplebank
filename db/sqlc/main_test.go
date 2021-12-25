package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/otmosina/simplebank/util"
)

// const (
// 	dbDriver = "postgres"
// 	dbSource = "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
// )

var testQueries *Queries

var testDB *sql.DB

func TestMain(m *testing.M) {
	fmt.Println("MAIN_TEST SQLC SQLC SQLC SQLC SQLC SQLC SQLC SQLC SQLC SQLC SQLC SQLC SQLC SQLC SQLC ")
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("error when load config", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
	// testQueries.CreateAccount()

}
