package db

import (
	"database/sql"
	"friendfy-api/src/config"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.StringDatabaseConn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
