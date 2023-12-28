package config

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Sqlite *sql.DB
	err    error
)

func init() {
	Sqlite, err = sql.Open("sqlite3", "store.db")
	if err != nil {
		log.Fatal(err)
	}
	if err := Sqlite.Ping(); err != nil {
		log.Fatal(err)
	}
}
