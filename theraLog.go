package main

import (
	"fmt"
	"theraLog/dataRepository/dbLayer"
	patientDataModel "theraLog/dataRepository/patient"

	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
)

const dbFile string = "therapyLog.sqlite3"

var dbHandler *sqlx.DB = nil

// getPatients responds with the list of all patiens as JSON.
func getPatients(c *gin.Context) {

	patients := dbLayer.FindBy[patientDataModel.Patient](dbHandler, nil)
	c.IndentedJSON(http.StatusOK, patients)
}

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

	router := gin.Default()
	//https: //go.dev/doc/tutorial/web-service-gin
	router.GET("/patients", getPatients)
	router.Run("localhost:8080")

}
