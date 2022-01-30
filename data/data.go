package data

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func OpenDb() error {
	var err error

	db, err = sql.Open("sqlite3", "./sqlite-db.db")
	if err != nil {
		return err
	}
	return db.Ping()
}
