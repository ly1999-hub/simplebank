package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ly1999-hub/simplebank/util"
)

var testQuery *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("./..")
	if err != nil {
		log.Fatal("Cannot load file config env :", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to postgresql :", err)
	}
	testQuery = New(testDB)
	os.Exit(m.Run())
}
