package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Jiang-Gianni/chat/dfrr"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nats-io/nuid"
)

var (
	StoreDB = "store.db"
	sqlite3 = "sqlite3"
)

func Sqlite() *sql.DB {
	return NewSqlite(StoreDB)
}

func NewSqlite(filename string) *sql.DB {
	db, err := sql.Open(sqlite3, filename)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// Return a *sql.DB for a sqlite databate to use in testing.
func GetSqliteTest() (testSqlite *sql.DB, cleanup func(), err error) {
	defer dfrr.Wrap(&err, "GetSqliteTest")
	dbName := nuid.New().Next() + ".db"
	testSqlite, err = sql.Open(sqlite3, dbName)
	if err != nil {
		return
	}
	err = testSqlite.Ping()
	if err != nil {
		return
	}
	err = SqliteInit(testSqlite)
	if err != nil {
		return
	}
	cleanup = func() {
		os.Remove(dbName)
		testSqlite.Close()
	}
	return
}

// Function to initialize the tables
func SqliteInit(sqlite *sql.DB) (rerr error) {
	defer dfrr.Wrap(&rerr, "SqliteInit")
	fileList := []string{
		"../sql/drop.sql",
		"../sql/message.sql",
		"../sql/room.sql",
		"../sql/user.sql",
	}
	for i := range fileList {
		b, err := os.ReadFile(fileList[i])
		if err != nil {
			return fmt.Errorf("os.ReadFile: %w", err)
		}
		_, err = sqlite.Exec(string(b))
		if err != nil {
			return fmt.Errorf("sqlite.Exec: %w", err)
		}
	}
	return nil
}
