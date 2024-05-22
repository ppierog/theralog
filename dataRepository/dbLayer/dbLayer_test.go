package dbLayer

import (
	"log"
	"os"
	"testing"
	noteDataModel "theraLog/dataRepository/note"
	patientDataModel "theraLog/dataRepository/patient"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const dbFile string = "test.sqlite3"
const schemaFile string = "../../db.schemas.sql"

var dbHandler *sqlx.DB = nil

func TestMain(m *testing.M) {

	_, err := os.Stat(dbFile)
	if err == nil {
		e := os.Remove(dbFile)
		if e != nil {
			log.Fatal(e)
		}
	}
	file, err := os.Create(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	db, err := sqlx.Open("sqlite3", dbFile)

	if err != nil {
		log.Fatal("Could not open db", err)
	}

	dbHandler = db
	schemaStr, err := os.ReadFile(schemaFile)
	if err != nil {
		log.Fatal("Could not open schema file ", err)
	}

	Exec(dbHandler, QryBuilder{Qry: string(schemaStr)}.Get())
	defer func() {
		db.Close()
		dbHandler = nil

	}()

	status := m.Run()
	os.Exit(status)
}
func TestPing(t *testing.T) {

	err := dbHandler.Ping()
	if err != nil {
		t.Fatalf("Could not ping DB err %s: ", err)
	}
}
func TestNoObjects(t *testing.T) {
	patients := FindBy[patientDataModel.Patient](dbHandler, QryBuilder{}.Get())
	notes := FindBy[noteDataModel.Note](dbHandler, QryBuilder{}.Get())

	if len(patients) != 0 {
		t.Fatalf("Wrong DB initial state")
	}
	if len(notes) != 0 {
		t.Fatalf("Wrong DB initial state")
	}
}
func TestPatients(t *testing.T) {
	patients := FindBy[patientDataModel.Patient](dbHandler, QryBuilder{}.Get())

	if len(patients) != 0 {
		t.Fatalf("Wrong DB initial state")
	}
	pateintsTestVector :=
		[]patientDataModel.Patient{
			{Name: "Piotr Pierog", Occupation: "Sw Developer",
				City: "Krakow", TelephoneNumber: "+48660345416", BirthYear: 1982},

			{Name: "Zuzanna Pierog", Occupation: "Student",
				City: "Krakow", TelephoneNumber: "+48760300300", BirthYear: 2007},

			{Name: "Marta Pierog", Occupation: "Therapist",
				City: "Krakow", TelephoneNumber: "+48760300300", BirthYear: 1979},
		}

	for i := 0; i < len(pateintsTestVector); i++ {
		Insert(dbHandler, &pateintsTestVector[i])
	}

	patients = FindBy[patientDataModel.Patient](dbHandler, QryBuilder{}.Get())

	if len(patients) != len(pateintsTestVector) {
		t.Fatalf("Could not insert to DB")
	}
	for i := 0; i < len(pateintsTestVector); i++ {
		if !patientDataModel.Equal(&pateintsTestVector[i], &patients[i]) {
			t.Fatalf("Objects are not Equal %v vs %v", pateintsTestVector[i], patients[i])
		}
	}

}
