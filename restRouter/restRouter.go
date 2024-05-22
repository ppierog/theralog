package restRouter

import (
	"database/sql"
	"theraLog/dataRepository/dataModel"
	"theraLog/dataRepository/dbLayer"

	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
)

type RestRouter struct {
	dbHandler *sqlx.DB
	engine    *gin.Engine
}

func getObjects[T dbLayer.DbTable, PT interface {
	Init(rows *sql.Rows) error
	TableName() string
	*T
}](dbHandler *sqlx.DB, c *gin.Context) {
	objects := dbLayer.FindBy[T, PT](dbHandler, nil)
	c.IndentedJSON(http.StatusOK, objects)
}

func getObjectById[T dbLayer.DbTable, PT interface {
	Init(rows *sql.Rows) error
	TableName() string
	*T
}](dbHandler *sqlx.DB, c *gin.Context) {

	id := c.Param("id")
	rowId, _ := strconv.ParseInt(id, 10, 64)
	var object T
	if dbLayer.FindByRowId[T, PT](dbHandler, rowId, &object) {
		c.IndentedJSON(http.StatusOK, object)
	} else {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
}

func postObject[T dbLayer.DbOps](dbHandler *sqlx.DB, c *gin.Context, t T) {

	if err := c.BindJSON(t); err != nil {
		c.JSON(400, gin.H{"code": "BAD_REQUEST", "message": "Bad Request : Could not deserialize data"})
		return
	}
	dbLayer.Insert(dbHandler, t)
}

func (r *RestRouter) resetDB(c *gin.Context) {
	dbLayer.DeleteBy(r.dbHandler, nil, &dataModel.Patient{})
	dbLayer.DeleteBy(r.dbHandler, nil, &dataModel.Note{})
}

func (r *RestRouter) initializeDB(c *gin.Context) {
	pateintsTestVector :=
		[]dataModel.Patient{
			{Name: "Patient 1", Occupation: "Sw Developer",
				City: "Krakow Krowodrza", TelephoneNumber: "+486111111111", BirthYear: 1982},

			{Name: "Patient 2", Occupation: "Student",
				City: "Krakow Pradnik", TelephoneNumber: "+48222222222", BirthYear: 2007},

			{Name: "Patient 3", Occupation: "Therapist",
				City: "Krakow Bronowice", TelephoneNumber: "+48760300300", BirthYear: 1979},
		}
	notesTestVector := []dataModel.Note{
		{Name: "Note 1", PatientRowId: 1,
			SessionDate: 1, NoteDate: 1, FileName: "test1.txt", IsCrypted: false},
		{Name: "Note 2", PatientRowId: 1,
			SessionDate: 1, NoteDate: 1, FileName: "test1.txt", IsCrypted: false},
		{Name: "Note 3", PatientRowId: 1,
			SessionDate: 1, NoteDate: 1, FileName: "test1.txt", IsCrypted: false},
	}

	for i := 0; i < len(pateintsTestVector); i++ {
		dbLayer.Insert(r.dbHandler, &pateintsTestVector[i])
		dbLayer.Insert(r.dbHandler, &notesTestVector[i])
	}
}

func (r *RestRouter) getPatients(c *gin.Context) {
	getObjects[dataModel.Patient](r.dbHandler, c)
}

func (r *RestRouter) getPatientById(c *gin.Context) {
	getObjectById[dataModel.Patient](r.dbHandler, c)
}

func (r *RestRouter) postPatient(c *gin.Context) {
	newPatient := dataModel.Patient{}
	postObject(r.dbHandler, c, &newPatient)
}

func (r *RestRouter) getNotes(c *gin.Context) {
	getObjects[dataModel.Note](r.dbHandler, c)
}

func (r *RestRouter) getNoteById(c *gin.Context) {
	getObjectById[dataModel.Note](r.dbHandler, c)
}

func (r *RestRouter) postNote(c *gin.Context) {
	newNote := dataModel.Note{}
	postObject(r.dbHandler, c, &newNote)
}

func (r RestRouter) Get() *RestRouter {
	return &r
}

func (r *RestRouter) GetEngine() *gin.Engine {
	return r.engine
}

func (r *RestRouter) Init(dbHandler *sqlx.DB) *RestRouter {
	r.dbHandler = dbHandler

	//https: //go.dev/doc/tutorial/web-service-gin
	r.engine = gin.Default()
	r.engine.POST("/reset", r.resetDB)
	r.engine.POST("/initialize", r.initializeDB)

	r.engine.GET("/patients", r.getPatients)
	r.engine.GET("/patients/:id", r.getPatientById)
	r.engine.POST("/patients", r.postPatient)

	r.engine.GET("/notes", r.getNotes)
	r.engine.GET("/notes/:id", r.getNoteById)
	r.engine.POST("/notes", r.postNote)

	return r

}
