package dbLayer

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"theraLog/dataRepository/dataModel"

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
		e := os.Remove(dbFile)
		if e != nil {
			log.Fatal(e)
		}

	}()

	m.Run()

}
func TestPing(t *testing.T) {

	err := dbHandler.Ping()
	if err != nil {
		t.Fatalf("Could not ping DB err %s: ", err)
	}
}

func testNoObject[T DbTable, PT interface {
	Init(rows *sql.Rows) error
	TableName() string
	*T
}](t *testing.T) {
	objects := FindBy[T, PT](dbHandler, QryBuilder{}.Get())
	if len(objects) != 0 {
		t.Fatalf("Wrong DB initial state")
	}
}

func testNoObjectByName[T DbTable, PT interface {
	Init(rows *sql.Rows) error
	TableName() string
	*T
}](b *testing.B, prefix string) {
	var obj T
	if FindByName[T, PT](dbHandler, prefix, &obj) {
		b.Fatalf("Wrong DB initial state")
	}
}

func testNoObjectByRowId[T DbTable, PT interface {
	Init(rows *sql.Rows) error
	TableName() string
	*T
}](b *testing.B, id int64) {
	var obj T
	if FindByRowId[T, PT](dbHandler, id, &obj) {
		b.Fatalf("Wrong DB initial state")
	}
}

func TestNoObjects(t *testing.T) {
	testNoObject[dataModel.Patient](t)
	testNoObject[dataModel.Note](t)
	testNoObject[dataModel.User](t)
	testNoObject[dataModel.PatientManifest](t)
}
func BenchmarkNoObjects(b *testing.B) {

	for i := 0; i < b.N; i++ {
		testNoObjectByName[dataModel.Patient](b, fmt.Sprintf("Patient %d", i))
		testNoObjectByName[dataModel.Note](b, fmt.Sprintf("Note %d", i))
		testNoObjectByName[dataModel.User](b, fmt.Sprintf("User %d", i))
		testNoObjectByRowId[dataModel.Patient](b, int64(i))
		testNoObjectByRowId[dataModel.Note](b, int64(i))
		testNoObjectByRowId[dataModel.User](b, int64(i))
	}
}

func TestPatients(t *testing.T) {
	patients := FindBy[dataModel.Patient](dbHandler, QryBuilder{}.Get())

	if len(patients) != 0 {
		t.Fatalf("Wrong DB initial state")
	}
	pateintsTestVector := dataModel.InitialTestVevtor{}.Patients()

	for i := 0; i < len(pateintsTestVector); i++ {
		Insert(dbHandler, &pateintsTestVector[i])
	}

	patients = FindBy[dataModel.Patient](dbHandler, QryBuilder{}.Get())

	fatalIfNotEqual := func(p1 *dataModel.Patient, p2 *dataModel.Patient) {
		if !dataModel.Equal(p1, p2) {
			t.Fatalf("Objects are not Equal %v vs %v", p1, p2)
		}
	}
	fatalIfEqual := func(p1 *dataModel.Patient, p2 *dataModel.Patient) {
		if dataModel.Equal(p1, p2) {
			t.Fatalf("Objects are Equal %v vs %v", p1, p2)
		}
	}
	fatalIfNotEqualTables := func(p1 []dataModel.Patient, p2 []dataModel.Patient) {
		if len(p1) != len(p2) {
			t.Fatalf("Len p1 %d vs Len p2 %d", len(p1), len(p2))
		}
		for i := 0; i < len(p1); i++ {
			fatalIfNotEqual(&p1[i], &p2[i])
		}
	}

	fatalIfNotEqualTables(pateintsTestVector, patients)

	var newPatient dataModel.Patient

	for i := 0; i < len(pateintsTestVector); i++ {
		newPatient = dataModel.Patient{}
		fatalIfEqual(&pateintsTestVector[i], &newPatient)

		FindByName(dbHandler, pateintsTestVector[i].Name+"1", &newPatient)

		fatalIfEqual(&pateintsTestVector[i], &newPatient)

		FindByName(dbHandler, pateintsTestVector[i].Name, &newPatient)

		fatalIfNotEqual(&pateintsTestVector[i], &newPatient)

		newPatient = dataModel.Patient{}
		FindByRowId(dbHandler, int64(i+1), &newPatient)

		fatalIfNotEqual(&pateintsTestVector[i], &newPatient)

	}

	pateintsTestVector[0].Name = "Patient 111"
	pateintsTestVector[0].City = "Poznan Lawica"

	pateintsTestVector[1].Name = "Patient 222"
	pateintsTestVector[1].Name = "Poznan Wola"

	pateintsTestVector[2].Name = "Patient 333"
	pateintsTestVector[2].Name = "Poznan Winiary"

	for i := 0; i < len(pateintsTestVector); i++ {
		Update(dbHandler, &pateintsTestVector[i])
	}
	patients = FindBy[dataModel.Patient](dbHandler, QryBuilder{}.Get())

	fatalIfNotEqualTables(pateintsTestVector, patients)

	DeleteByName(dbHandler, &pateintsTestVector[0])
	patients = FindBy[dataModel.Patient](dbHandler, QryBuilder{}.Get())
	if len(patients) != 2 {
		t.Fatalf("Wrong length of patients , expected 2")
	}

	DeleteByRowId(dbHandler, &pateintsTestVector[1])
	patients = FindBy[dataModel.Patient](dbHandler, QryBuilder{}.Get())
	if len(patients) != 1 {
		t.Fatalf("Wrong length of patients , expected 2")
	}

	DeleteBy(dbHandler, QryBuilder{}.Get(), &dataModel.Patient{})
	patients = FindBy[dataModel.Patient](dbHandler, QryBuilder{}.Get())
	if len(patients) != 0 {
		t.Fatalf("Wrong length of patients , expected 2")
	}

}
