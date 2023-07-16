package simplebank

import (
	"database/sql"
	"log"
	"testing"
	"os"
	_"github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB
const (
	dbDriver = "postgres"
	source = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func TestMain( m *testing.M){
	var err error
	testDB, err = sql.Open(dbDriver , source)

	if err != nil{
		log.Fatal("cannot connect")
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}