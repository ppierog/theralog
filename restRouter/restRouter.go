package restRouter

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"theraLog/cred"
	"theraLog/dataRepository/dataModel"
	"theraLog/dataRepository/dbLayer"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
	"github.com/swaggest/openapi-go/openapi3"
)

const patientsURL = "/patients"
const patientsByIdURL = "/patients/:id"
const notesURL = "/notes"
const notesByIdURL = "/notes/:id"
const notesByIdUploadURL = "/notes/:id/upload"
const usersURL = "/users"
const usersByIdURL = "/users/:id"
const manifestsURL = "/manifests"
const manifestsByIdURL = "/manifests/:id"
const appdataPath = "appdata/"

type RestRouter struct {
	dbHandler *sqlx.DB
	engine    *gin.Engine
	secret    string
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
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Object not found"})

	}
}

func postObject[T dbLayer.DbOps](dbHandler *sqlx.DB, c *gin.Context, t T, prepare func(t T)) {

	if err := c.BindJSON(t); err != nil {
		c.JSON(400, gin.H{"code": "BAD_REQUEST", "message": "Bad Request : Could not deserialize data : " + err.Error()})
		return
	}
	if prepare != nil {
		prepare(t)
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

func putObject[T dbLayer.DbOps](dbHandler *sqlx.DB, c *gin.Context, t T) {

	if err := c.BindJSON(t); err != nil {
		c.JSON(400, gin.H{"code": "BAD_REQUEST", "message": "Bad Request : Could not deserialize data : " + err.Error()})
		return
	}
	id := c.Param("id")
	rowId, _ := strconv.ParseInt(id, 10, 64)
	t.SetRowId(rowId)

	dbLayer.Update(dbHandler, t)
}

func (r *RestRouter) resetDB(c *gin.Context) {
	dbLayer.DeleteBy(r.dbHandler, nil, &dataModel.Patient{})
	dbLayer.DeleteBy(r.dbHandler, nil, &dataModel.Note{})
	dbLayer.DeleteBy(r.dbHandler, nil, &dataModel.User{})
	dbLayer.DeleteBy(r.dbHandler, nil, &dataModel.PatientManifest{})
}

func (r *RestRouter) initDB(c *gin.Context) {
	testData := dataModel.InitialTestVevtor{}

	for _, patient := range testData.Patients() {
		dbLayer.Insert(r.dbHandler, &patient)
	}

	for _, note := range testData.Notes() {
		dbLayer.Insert(r.dbHandler, &note)
	}

	for _, user := range testData.Users() {
		dbLayer.Insert(r.dbHandler, &user)
	}

	for _, manifest := range testData.Manifests() {
		dbLayer.Insert(r.dbHandler, &manifest)
	}
}

type blankReq struct{}
type blankRes struct{}

type idReq struct {
	ID string `path:"id" example:"1"`
}

type idReqWithContent[T any] struct {
	ID  string `path:"id" example:"1"`
	OBJ T      `json:"obj"`
}

func AddOperation[REQ any, RES any](reflector *openapi3.Reflector, req *REQ, res *RES, method string, url string) {
	op := openapi3.Operation{}
	reflector.SetRequest(&op, req, method)
	reflector.SetJSONResponse(&op, new(RES), http.StatusOK)
	reflector.Spec.AddOperation(method, url, op)
}

func AddOperations[RES any](reflector *openapi3.Reflector, url string) {
	AddOperation[blankReq](reflector, nil, new([]RES), http.MethodGet, url)
	AddOperation(reflector, new(idReq), new(RES), http.MethodGet, url+"/{id}")
	AddOperation[RES, blankRes](reflector, new(RES), nil, http.MethodPost, url)
	AddOperation[idReq, blankRes](reflector, new(idReq), nil, http.MethodDelete, url+"/{id}")
	AddOperation[idReqWithContent[RES], blankRes](reflector, new(idReqWithContent[RES]), nil, http.MethodPut, url+"/{id}")
}

func (r *RestRouter) getApi(c *gin.Context) {
	reflector := openapi3.Reflector{}
	reflector.Spec = &openapi3.Spec{Openapi: "3.0.3"}
	reflector.Spec.Info.
		WithTitle("TheraLog Api").
		WithVersion("1.00").
		WithDescription("TheraLog Api description")

	AddOperations[dataModel.Patient](&reflector, patientsURL)
	AddOperations[dataModel.Note](&reflector, notesURL)
	AddOperations[dataModel.User](&reflector, usersURL)
	AddOperations[dataModel.PatientManifest](&reflector, manifestsURL)

	schema, err := reflector.Spec.MarshalYAML()
	if err != nil {
		log.Fatal(err)
	}

	s := string(schema)

	c.String(200, s)

}

func (r *RestRouter) login(c *gin.Context) {
	creditionals := dataModel.UserCred{}

	if err := c.BindJSON(&creditionals); err != nil {
		c.JSON(400, gin.H{"code": "BAD_REQUEST", "message": "Bad Request : Could not deserialize data : " + err.Error()})
		return
	}

	qry := dbLayer.QryBuilder{}
	users := dbLayer.FindBy[dataModel.User](r.dbHandler, qry.Where("email").Is(creditionals.Email))
	if len(users) != 1 {
		c.JSON(404, gin.H{"code": "NOT_FOUND", "message": "Could not login"})
		return
	}

	user := users[0]
	saltedPasswdSha256 := cred.CalcSha256(creditionals.Password, user.Salt)
	if saltedPasswdSha256 != user.Password {
		c.JSON(404, gin.H{"code": "NOT_FOUND", "message": "Could not login"})
		return

	}

	tokenRepo := cred.TokenRepository{Secret: r.secret}
	jwt, err := tokenRepo.NewAccessToken(cred.UserClaims{Email: creditionals.Email, StandardClaims: jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
	}})

	if err != nil {
		c.JSON(500, gin.H{"code": "Internal Server Error", "message": "Could not generate token"})
		return
	}
	c.IndentedJSON(http.StatusOK, cred.Token{Jwt: jwt})
}

func (r *RestRouter) checkToken() gin.HandlerFunc {

	return func(c *gin.Context) {

		if c.Request.URL.String() == "/login" {
			c.Next()
		} else {
			token := c.Request.Header.Get("token")
			tokenRepo := cred.TokenRepository{Secret: r.secret}
			user := tokenRepo.ParseAccessToken(token)
			if user == nil {
				c.AbortWithStatusJSON(401, gin.H{"error": "No Auth token in header request"})
				return
			}
			currTime := time.Now().Unix()
			if user.ExpiresAt <= currTime {
				c.AbortWithStatusJSON(401, gin.H{"error": "Token already expired : " + time.Unix(user.ExpiresAt, 0).String()})
				return
			}
			// Set example variable
			// c.Set("example", "12345")
			c.Next()
		}

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
	postObject(r.dbHandler, c, &newPatient, nil)
	patientDirName := appdataPath + strconv.Itoa(int(newPatient.RowId))

	if err := os.Mkdir(patientDirName, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func (r *RestRouter) deletePatient(c *gin.Context) {
	patient := dataModel.Patient{}
	deleteObject(r.dbHandler, c, &patient)

	patientDirName := appdataPath + strconv.Itoa(int(patient.RowId))

	err := os.RemoveAll(patientDirName)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *RestRouter) putPatient(c *gin.Context) {
	putObject(r.dbHandler, c, &dataModel.Patient{})
}

func (r *RestRouter) getNotes(c *gin.Context) {
	getObjects[dataModel.Note](r.dbHandler, c)
}

func (r *RestRouter) getNoteById(c *gin.Context) {
	getObjectById[dataModel.Note](r.dbHandler, c)
}

func (r *RestRouter) postNote(c *gin.Context) {
	newNote := dataModel.Note{}
	postObject(r.dbHandler, c, &newNote, nil)
}

func (r *RestRouter) uploadNote(c *gin.Context) {

	id := c.Param("id")
	rowId, _ := strconv.ParseInt(id, 10, 64)
	note := dataModel.Note{}

	if !dbLayer.FindByRowId(r.dbHandler, rowId, &note) {
		c.String(http.StatusNotFound, "Could not find appropriate note!")
		return
	}

	isCrypted, _ := strconv.ParseBool(c.Request.Header.Get("isCrypted"))
	note.IsCrypted = isCrypted

	file, _ := c.FormFile("file")
	note.FileName = file.Filename
	dbLayer.Update(r.dbHandler, &note)

	dest := fmt.Sprintf("%s%d/%s", appdataPath, note.PatientRowId, note.FileName)
	c.SaveUploadedFile(file, dest)
	c.String(http.StatusOK, fmt.Sprintf("%s uploaded!\n", file.Filename))
}

func (r *RestRouter) deleteNote(c *gin.Context) {
	deleteObject(r.dbHandler, c, &dataModel.Note{})
}

func (r *RestRouter) putNote(c *gin.Context) {
	putObject(r.dbHandler, c, &dataModel.Note{})
}

func (r *RestRouter) getUsers(c *gin.Context) {
	getObjects[dataModel.User](r.dbHandler, c)
}

func (r *RestRouter) getUserById(c *gin.Context) {
	getObjectById[dataModel.User](r.dbHandler, c)
}

func (r *RestRouter) postUser(c *gin.Context) {
	newUser := dataModel.User{}
	postObject(r.dbHandler, c, &newUser, func(user *dataModel.User) {
		user.Salt = cred.GenerateSalt()
		saltedPasswdSha256 := cred.CalcSha256(user.Password, user.Salt)
		user.Password = saltedPasswdSha256
	})
}

func (r *RestRouter) deleteUser(c *gin.Context) {
	deleteObject(r.dbHandler, c, &dataModel.User{})
}

func (r *RestRouter) putUser(c *gin.Context) {
	putObject(r.dbHandler, c, &dataModel.User{})
}

func (r *RestRouter) getManifests(c *gin.Context) {
	getObjects[dataModel.PatientManifest](r.dbHandler, c)
}

func (r *RestRouter) getManifestById(c *gin.Context) {
	getObjectById[dataModel.PatientManifest](r.dbHandler, c)
}

func (r *RestRouter) postManifest(c *gin.Context) {
	newUser := dataModel.PatientManifest{}
	postObject(r.dbHandler, c, &newUser, nil)
}

func (r *RestRouter) deleteManifest(c *gin.Context) {
	deleteObject(r.dbHandler, c, &dataModel.PatientManifest{})
}

func (r *RestRouter) putManifest(c *gin.Context) {
	putObject(r.dbHandler, c, &dataModel.PatientManifest{})
}

func (r *RestRouter) GetEngine() *gin.Engine {
	return r.engine
}

func (r *RestRouter) Init(dbHandler *sqlx.DB, secret string) *RestRouter {
	r.dbHandler = dbHandler
	r.secret = secret

	//https: //go.dev/doc/tutorial/web-service-gin
	r.engine = gin.Default()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.engine.MaxMultipartMemory = 8 << 20 // 8 MiB

	r.engine.Use(r.checkToken())

	r.engine.POST("/reset", r.resetDB)
	r.engine.POST("/init", r.initDB)
	r.engine.GET("api", r.getApi)

	r.engine.POST("/login", r.login)

	r.engine.GET(patientsURL, r.getPatients)
	r.engine.GET(patientsByIdURL, r.getPatientById)
	r.engine.POST(patientsURL, r.postPatient)
	r.engine.DELETE(patientsByIdURL, r.deletePatient)
	r.engine.PUT(patientsByIdURL, r.putPatient)

	r.engine.GET(notesURL, r.getNotes)
	r.engine.GET(notesByIdURL, r.getNoteById)
	r.engine.POST(notesURL, r.postNote)
	r.engine.POST(notesByIdUploadURL, r.uploadNote)

	r.engine.DELETE(notesByIdURL, r.deleteNote)
	r.engine.PUT(notesByIdURL, r.putNote)

	r.engine.GET(usersURL, r.getUsers)
	r.engine.GET(usersByIdURL, r.getUserById)
	r.engine.POST(usersURL, r.postUser)
	r.engine.DELETE(usersByIdURL, r.deleteUser)
	r.engine.PUT(usersByIdURL, r.putUser)

	r.engine.GET(manifestsURL, r.getManifests)
	r.engine.GET(manifestsByIdURL, r.getManifestById)
	r.engine.POST(manifestsURL, r.postManifest)
	r.engine.DELETE(manifestsByIdURL, r.deleteManifest)
	r.engine.PUT(manifestsByIdURL, r.putManifest)

	return r
}
