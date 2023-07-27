package simplebank

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/AlhajiTaibu/simplebank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain( m *testing.M){
	var err error

	config, err := util.LoadConfig("..")

	if err != nil {
		log.Fatal("Unable to load config", err)
	}
	testDB, err = sql.Open(config.DbDriver , config.DbSource)

	if err != nil{
		log.Fatal("cannot connect")
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}