package config

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var sqliteDB *sqlx.DB

func init() {
	dbPath := os.Getenv("SQLITE_DB")

	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("config.sqlite:not_found::%+v", err)
	}

	sqliteDB = db
}

func GetSqliteDB() *sqlx.DB {
	return sqliteDB
}
