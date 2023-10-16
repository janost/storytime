package model

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

func Setup(dbhost, dbname, dbuser, dbpass string) {

	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", dbhost, dbuser, dbpass, dbname))
	if err != nil {
		fmt.Println("Could not connect to db", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Could not ping db", err)
	}
}
