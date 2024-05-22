package main

import (
	"fmt"
	"theraLog/dataRepository/dbLayer"
	noteDataModel "theraLog/dataRepository/note"
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
	dbLayer.DeleteBy(dbHandler, nil, &noteDataModel.Note{})
}

func initializeDB(c *gin.Context) {
	pateintsTestVector :=
		[]patientDataModel.Patient{
			{Name: "Patient 1", Occupation: "Sw Developer",
				City: "Krakow Krowodrza", TelephoneNumber: "+486111111111", BirthYear: 1982},

			{Name: "Patient 2", Occupation: "Student",
				City: "Krakow Pradnik", TelephoneNumber: "+48222222222", BirthYear: 2007},

			{Name: "Patient 3", Occupation: "Therapist",
				City: "Krakow Bronowice", TelephoneNumber: "+48760300300", BirthYear: 1979},
		}
	notesTestVector := []noteDataModel.Note{
		{Name: "Note 1", PatientRowId: 1,
			SessionDate: 1, NoteDate: 1, FileName: "test1.txt", IsCrypted: false},
		{Name: "Note 2", PatientRowId: 1,
			SessionDate: 1, NoteDate: 1, FileName: "test1.txt", IsCrypted: false},
		{Name: "Note 3", PatientRowId: 1,
			SessionDate: 1, NoteDate: 1, FileName: "test1.txt", IsCrypted: false},
	}

	for i := 0; i < len(pateintsTestVector); i++ {
		dbLayer.Insert(dbHandler, &pateintsTestVector[i])
		dbLayer.Insert(dbHandler, &notesTestVector[i])
	}

}

func getPatients(c *gin.Context) {

	patients := dbLayer.FindBy[patientDataModel.Patient](dbHandler, nil)
	c.IndentedJSON(http.StatusOK, patients)
}

func getNotes(c *gin.Context) {
	notes := dbLayer.FindBy[noteDataModel.Note](dbHandler, nil)
	c.IndentedJSON(http.StatusOK, notes)
}

func getPatientById(c *gin.Context) {

	id := c.Param("id")
	rowId, _ := strconv.ParseInt(id, 10, 64)
	patient := patientDataModel.Patient{}
	if dbLayer.FindByRowId(dbHandler, rowId, &patient) {
		c.IndentedJSON(http.StatusOK, patient)
	} else {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
}

func getNoteById(c *gin.Context) {

	id := c.Param("id")
	rowId, _ := strconv.ParseInt(id, 10, 64)
	note := noteDataModel.Note{}
	if dbLayer.FindByRowId(dbHandler, rowId, &note) {
		c.IndentedJSON(http.StatusOK, note)
	} else {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
}

func postPatient(c *gin.Context) {

	var newPatient patientDataModel.Patient

	if err := c.BindJSON(&newPatient); err != nil {
		c.JSON(400, gin.H{"code": "BAD_REQUEST", "message": "Bad Request : Could not deserialize patient data"})
		return
	}
	dbLayer.Insert(dbHandler, &newPatient)

}

func postNote(c *gin.Context) {

	/*
		var newPatient patientDataModel.Patient

		if err := c.BindJSON(&newPatient); err != nil {
			c.JSON(400, gin.H{"code": "BAD_REQUEST", "message": "Bad Request : Could not deserialize patient data"})
			return
		}
		dbLayer.Insert(dbHandler, &newPatient)
	*/
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

	router.GET("/notes", getNotes)
	router.GET("/notess/:id", getNoteById)
	router.POST("/notes", postNote)

	router.Run("localhost:8080")

}
