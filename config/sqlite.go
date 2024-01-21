package config

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const file string = "alamak.db"

var sqlitedb *sql.DB

// Returns the generated client.
func GetSqliteClient() *sql.DB {
	return sqlitedb
}

func init() {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalf("DB not found, %+v", err)
	}
	sqlitedb = db
}
