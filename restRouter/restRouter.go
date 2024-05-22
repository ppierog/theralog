package restRouter

import (
	"database/sql"
	"log"
	"theraLog/dataRepository/dataModel"
	"theraLog/dataRepository/dbLayer"

	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
	"github.com/swaggest/openapi-go/openapi3"
)

const patientsURL = "/patients"
const patientsByIdURL = "/patients/:id"
const notesURL = "/notes"
const notesByIdURL = "/notes/:id"
const usersURL = "/users"
const usersByIdURL = "/users/:id"

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

func deleteObject[T dbLayer.DbOps](dbHandler *sqlx.DB, c *gin.Context, t T) {
	id := c.Param("id")
	rowId, _ := strconv.ParseInt(id, 10, 64)
	t.SetRowId(rowId)
	if 0 == dbLayer.DeleteByRowId(dbHandler, t) {
		c.JSON(400, gin.H{"code": "BAD_REQUEST", "message": "Bad Request : Could not delete object"})
	}
}

func (r *RestRouter) resetDB(c *gin.Context) {
	dbLayer.DeleteBy(r.dbHandler, nil, &dataModel.Patient{})
	dbLayer.DeleteBy(r.dbHandler, nil, &dataModel.Note{})
	dbLayer.DeleteBy(r.dbHandler, nil, &dataModel.User{})
}

func (r *RestRouter) initDB(c *gin.Context) {
	testData := dataModel.InitialTestVevtor{}

	pateintsTestVector := testData.Patients()
	notesTestVector := testData.Notes()
	userTestVector := testData.Users()

	for i := 0; i < len(pateintsTestVector); i++ {
		dbLayer.Insert(r.dbHandler, &pateintsTestVector[i])
	}

	for i := 0; i < len(notesTestVector); i++ {
		dbLayer.Insert(r.dbHandler, &notesTestVector[i])
	}

	for _, user := range userTestVector {
		dbLayer.Insert(r.dbHandler, &user)
	}
}

func AddOperation[REQ any, RES any](reflector *openapi3.Reflector, req *REQ, res *RES, method string, url string) {
	op := openapi3.Operation{}
	reflector.SetRequest(&op, req, method)
	reflector.SetJSONResponse(&op, new(RES), http.StatusOK)
	reflector.Spec.AddOperation(method, url, op)
}

func (r *RestRouter) getApi(c *gin.Context) {
	reflector := openapi3.Reflector{}
	reflector.Spec = &openapi3.Spec{Openapi: "3.0.3"}
	reflector.Spec.Info.
		WithTitle("TheraLog Api").
		WithVersion("1.00").
		WithDescription("TheraLog Api description")

	type blankReq struct{}
	type blankRes struct{}

	type idReq struct {
		ID string `path:"id" example:"1"`
	}

	type patientsResp struct {
		patients []dataModel.Patient
	}

	AddOperation[blankReq](&reflector, nil, new([]dataModel.Patient), http.MethodGet, patientsURL)
	AddOperation(&reflector, new(idReq), new(dataModel.Patient), http.MethodGet, patientsURL+"/{id}")
	AddOperation[blankReq, blankRes](&reflector, nil, nil, http.MethodPost, patientsURL)
	AddOperation[idReq, blankRes](&reflector, new(idReq), nil, http.MethodDelete, patientsURL+"/{id}")

	AddOperation[blankReq](&reflector, nil, new([]dataModel.Note), http.MethodGet, notesURL)
	AddOperation(&reflector, new(idReq), new(dataModel.Note), http.MethodGet, notesURL+"/{id}")
	AddOperation[blankReq, blankRes](&reflector, nil, nil, http.MethodPost, notesURL)
	AddOperation[idReq, blankRes](&reflector, new(idReq), nil, http.MethodDelete, notesURL+"/{id}")

	AddOperation[blankReq](&reflector, nil, new([]dataModel.User), http.MethodGet, usersURL)
	AddOperation(&reflector, new(idReq), new(dataModel.User), http.MethodGet, usersURL+"/{id}")
	AddOperation[blankReq, blankRes](&reflector, nil, nil, http.MethodPost, usersURL)
	AddOperation[idReq, blankRes](&reflector, new(idReq), nil, http.MethodDelete, usersURL+"/{id}")

	schema, err := reflector.Spec.MarshalYAML()
	if err != nil {
		log.Fatal(err)
	}

	s := string(schema)

	c.String(200, s)

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

func (r *RestRouter) deletePatient(c *gin.Context) {
	deleteObject(r.dbHandler, c, &dataModel.Patient{})
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
func (r *RestRouter) deleteNote(c *gin.Context) {
	deleteObject(r.dbHandler, c, &dataModel.Note{})
}

func (r *RestRouter) getUsers(c *gin.Context) {
	getObjects[dataModel.User](r.dbHandler, c)
}

func (r *RestRouter) getUserById(c *gin.Context) {
	getObjectById[dataModel.User](r.dbHandler, c)
}

func (r *RestRouter) postUser(c *gin.Context) {
	newUser := dataModel.User{}
	postObject(r.dbHandler, c, &newUser)
}

func (r *RestRouter) deleteUser(c *gin.Context) {
	deleteObject(r.dbHandler, c, &dataModel.User{})
}

func (r *RestRouter) GetEngine() *gin.Engine {
	return r.engine
}

func (r *RestRouter) Init(dbHandler *sqlx.DB) *RestRouter {
	r.dbHandler = dbHandler

	//https: //go.dev/doc/tutorial/web-service-gin
	r.engine = gin.Default()
	r.engine.POST("/reset", r.resetDB)
	r.engine.POST("/init", r.initDB)
	r.engine.GET("api", r.getApi)

	r.engine.GET(patientsURL, r.getPatients)
	r.engine.GET(patientsByIdURL, r.getPatientById)
	r.engine.POST(patientsURL, r.postPatient)
	r.engine.DELETE(patientsByIdURL, r.deletePatient)

	r.engine.GET(notesURL, r.getNotes)
	r.engine.GET(notesByIdURL, r.getNoteById)
	r.engine.POST(notesURL, r.postNote)
	r.engine.DELETE(notesByIdURL, r.deleteNote)

	r.engine.GET(usersURL, r.getUsers)
	r.engine.GET(usersByIdURL, r.getUserById)
	r.engine.POST(usersURL, r.postUser)
	r.engine.DELETE(usersByIdURL, r.deleteUser)

	return r

}
