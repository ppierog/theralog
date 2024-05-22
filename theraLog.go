package main

import (
	"fmt"
	"theraLog/restRouter"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const dbFile string = "theraLog.sqlite3"

var dbHandler *sqlx.DB = nil

func main() {
	// https://earthly.dev/blog/golang-sqlite/
	// https://jmoiron.github.io/sqlx/
	db, err := sqlx.Open("sqlite3", dbFile)

	if err != nil {

		fmt.Println("Could not oper db", err)
		return
	}

	dbHandler = db
	defer func() {
		db.Close()
		dbHandler = nil
	}()

	err = db.Ping()
	if err != nil {
		fmt.Println("Could not ping DB : ", err)
		return
	}

	router := restRouter.RestRouter{}
	router.Init(dbHandler).GetEngine().Run("localhost:8080")

}
