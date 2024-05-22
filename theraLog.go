package main

import (
	"fmt"
	"theraLog/dataRepository/dbLayer"
	patientDataModel "theraLog/dataRepository/patient"

	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
)

const dbFile string = "theraLog.sqlite3"

var dbHandler *sqlx.DB = nil

func resetDB(c *gin.Context) {
	dbLayer.DeleteBy(dbHandler, nil, &patientDataModel.Patient{})
}

func initializeDB(c *gin.Context) {

	ppierog := patientDataModel.Patient{Name: "Piotr Pierog", Occupation: "Sw Developer",
		City: "Krakow", TelephoneNumber: "+48660345416", BirthYear: 1982}

	zpierog := patientDataModel.Patient{Name: "Zuzanna Pierog", Occupation: "Student",
		City: "Krakow", TelephoneNumber: "+48760300300", BirthYear: 2007}

	mpierog := patientDataModel.Patient{Name: "Marta Pierog", Occupation: "Therapist",
		City: "Krakow", TelephoneNumber: "+48760300300", BirthYear: 1979}

	dbLayer.Insert(dbHandler, &ppierog)
	dbLayer.Insert(dbHandler, &zpierog)
	dbLayer.Insert(dbHandler, &mpierog)
}

func getPatients(c *gin.Context) {

	patients := dbLayer.FindBy[patientDataModel.Patient](dbHandler, nil)
	c.IndentedJSON(http.StatusOK, patients)
}

func getPatientById(c *gin.Context) {

	id := c.Param("id")
	rowId, _ := strconv.ParseInt(id, 10, 64)
	patient := patientDataModel.Patient{}
	dbLayer.FindByRowId(dbHandler, rowId, &patient)
	c.IndentedJSON(http.StatusOK, patient)
}

func postPatient(c *gin.Context) {

	var newPatient patientDataModel.Patient

	if err := c.BindJSON(&newPatient); err != nil {
		return
	}
	dbLayer.Insert(dbHandler, &newPatient)

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

	router.POST("/reset", resetDB)
	router.POST("/initialize", initializeDB)
	router.GET("/patients", getPatients)
	router.GET("/patients/:id", getPatientById)
	router.POST("/patients", postPatient)
	router.Run("localhost:8080")

}
