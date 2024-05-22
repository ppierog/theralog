package main

import (
	"fmt"
	"log"
	"os"
	"theraLog/restRouter"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
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
		log.Fatalf("Could not ping DB : %s", err)
	}
	err = godotenv.Load("test.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	secret := os.Getenv("SECRET")
	router := restRouter.RestRouter{}

	const APP_PORT = 8080
	const APP_IP_ADDRESS_ENV = "APP_IP_ADDRESS"
	appIP := fmt.Sprintf("localhost:%d", APP_PORT)

	if os.Getenv(APP_IP_ADDRESS_ENV) != "" {
		appIP = fmt.Sprintf("%s:%d", os.Getenv(APP_IP_ADDRESS_ENV), APP_PORT)
	}

	router.Init(dbHandler, secret).GetEngine().Run(appIP)

}
