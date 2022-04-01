package mysql

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "mysql"
	dbSource = "root:@tcp(127.0.0.1:3306)/go-restapi-sample_test"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testDB = db
	os.Exit(m.Run())
}

func createTestTable(table string) {
	sql, err := ioutil.ReadFile("./schemas/" + table + ".sql")
	if err != nil {
		panic("cannot load sql file")
	}
	testDB.Exec(string(sql))
}

func deleteTestTable(table string) {
	sql := `DROP TABLE ` + table
	testDB.Exec(sql)
}
