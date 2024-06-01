package restRouter

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"theraLog/cred"
	"theraLog/dataRepository/dataModel"
	"theraLog/dataRepository/dbLayer"

	"time"

	"github.com/gin-contrib/cors"

	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
	"github.com/swaggest/openapi-go/openapi3"
)

const loginURL = "/login"
const apiURL = "/api"
const resetURL = "/reset"
const initURL = "/init"
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

type Id struct {
	RowId int64 `json:"id"`
}

func toJson[T dbLayer.DbTable](c *gin.Context, o dbLayer.DbObjectWrapper[T]) {
	if o.Err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": o.Err})
	} else {
		c.IndentedJSON(http.StatusOK, o.Data)
	}
}

func getObjects[T dbLayer.DbTable, PT interface {
	Init(rows *sql.Rows) error
	TableName() string
	*T
}](dbHandler *sqlx.DB, filter func(T) bool, modifier func(*T)) []T {
	objects := dbLayer.FindBy[T, PT](dbHandler, nil)
	ret := []T{}
	if filter != nil {
		for _, o := range objects {
			if filter(o) {
				ret = append(ret, o)
			}
		}
	} else {
		ret = objects
	}

	if modifier != nil {
		for i := 0; i < len(ret); i++ {
			modifier(&ret[i])
		}
	}

	return ret
}

func getObjectById[T dbLayer.DbTable, PT interface {
	Init(rows *sql.Rows) error
	TableName() string
	*T
}](dbHandler *sqlx.DB, c *gin.Context, modifier func(*T)) (T, error) {

	id := c.Param("id")
	rowId, _ := strconv.ParseInt(id, 10, 64)
	var object T
	if dbLayer.FindByRowId[T, PT](dbHandler, rowId, &object) {
		if modifier != nil {
			modifier(&object)
		}
		return object, nil
	} else {
		return object, errors.New("object not found")
	}
}

func postObject[T dbLayer.DbOps](dbHandler *sqlx.DB, c *gin.Context, t T, prepare func(t T)) {

	if err := c.BindJSON(t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "Bad Request : Could not deserialize data : " + err.Error()})
		return
	}
	if prepare != nil {
		prepare(t)
	}
	dbLayer.Insert(dbHandler, t)
	c.IndentedJSON(http.StatusOK, Id{RowId: t.GetRowId()})
}

func deleteObject[T dbLayer.DbOps](dbHandler *sqlx.DB, c *gin.Context, t T) {
	id := c.Param("id")
	rowId, _ := strconv.ParseInt(id, 10, 64)
	t.SetRowId(rowId)
	if dbLayer.DeleteByRowId(dbHandler, t) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "Bad Request : Could not delete object"})
	}
}

func putObject[T dbLayer.DbOps](dbHandler *sqlx.DB, c *gin.Context, t T, prepare func(t T)) {

	if err := c.BindJSON(t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "Bad Request : Could not deserialize data : " + err.Error()})
		return
	}
	id := c.Param("id")
	rowId, _ := strconv.ParseInt(id, 10, 64)
	t.SetRowId(rowId)
	if prepare != nil {
		prepare(t)
	}

	dbLayer.Update(dbHandler, t)
}

func (r *RestRouter) resetDB(c *gin.Context) {
	dbLayer.DeleteBy(r.dbHandler, nil, &dataModel.Patient{})
	dbLayer.DeleteBy(r.dbHandler, nil, &dataModel.Note{})
	dbLayer.DeleteBy(r.dbHandler, nil, &dataModel.User{})
	dbLayer.DeleteBy(r.dbHandler, nil, &dataModel.PatientManifest{})
	if err := os.RemoveAll(appdataPath); err != nil {
		log.Fatal(err)
	}

	if err := os.Mkdir(appdataPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

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

// @TODO pp :  Be more precise here
func AddOperation[REQ any, RES any](reflector *openapi3.Reflector, req *REQ, res *RES, method string, url string) {
	op := openapi3.Operation{}

	reflector.SetRequest(&op, req, method)
	reflector.SetJSONResponse(&op, new(RES), http.StatusOK)
	reflector.Spec.AddOperation(method, url, op)
}

// @TODO pp :  Be more precise here
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
		WithVersion("1.05").
		WithDescription("TheraLog Api description")

	op := openapi3.Operation{}
	reflector.Spec.AddOperation(http.MethodGet, apiURL, *op.WithDescription("Get YAML description of API"))

	AddOperation(&reflector, &dataModel.UserCred{}, &dataModel.UserToken{}, http.MethodPost, loginURL)
	reflector.Spec.AddOperation(http.MethodPost, resetURL, *op.WithDescription("Reset System to initial state"))
	reflector.Spec.AddOperation(http.MethodPost, initURL, *op.WithDescription("Init  System with default values"))

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
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "Bad Request : Could not deserialize data : " + err.Error()})
		return
	}

	qry := dbLayer.QryBuilder{}
	// qry := dbLayer.QryBuilder{}.Get().Where("email").IsEqual(dbLayer.SqlWrapValue(creditionals.Email)).Latch()
	qry.Where("email").IsEqual(dbLayer.SqlWrapValue(creditionals.Email))
	users := dbLayer.FindBy[dataModel.User](r.dbHandler, &qry)
	if len(users) != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "NOT_FOUND", "message": "Could not login"})
		return
	}
	user := users[0]
	saltedPasswdSha256 := cred.CalcSha256(creditionals.Password, user.Salt)
	if saltedPasswdSha256 != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"code": "UNAUTHORIZED", "message": "Could not login"})
		return

	}

	tokenRepo := cred.TokenRepository{Secret: r.secret}
	expiresAt := time.Now().Add(time.Minute * 15)
	jwt, err := tokenRepo.NewAccessToken(cred.UserClaims{Email: creditionals.Email, StandardClaims: jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: expiresAt.Unix(),
	}})

	if err != nil {
		c.JSON(500, gin.H{"code": "Internal Server Error", "message": "Could not generate token"})
		return
	}
	c.IndentedJSON(http.StatusOK, cred.Token{Jwt: jwt, ExpiresAt: expiresAt})
}

func (r *RestRouter) checkToken() gin.HandlerFunc {

	return func(c *gin.Context) {
		url := c.Request.URL.String()

		if url == loginURL || url == apiURL || url == initURL || url == resetURL {
			c.Next()
			return
		}
		if url == usersURL && c.Request.Method == http.MethodPost {
			c.Next()
			return
		}
		token := c.Request.Header.Get("Token")
		tokenRepo := cred.TokenRepository{Secret: r.secret}
		user := tokenRepo.ParseAccessToken(token)
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No Auth token in header request"})
			return
		}
		currTime := time.Now().Unix()
		if user.ExpiresAt <= currTime {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token already expired : " + time.Unix(user.ExpiresAt, 0).String()})
			return
		}

		// Set example variable
		// c.Set("example", "12345")
		c.Next()

	}

}

func (r *RestRouter) getPatients(c *gin.Context) {
	objects := getObjects[dataModel.Patient](r.dbHandler, nil, nil)
	c.IndentedJSON(http.StatusOK, objects)
}

func (r *RestRouter) getPatientById(c *gin.Context) {
	object := dbLayer.NewDbObjectWrapper(getObjectById[dataModel.Patient](r.dbHandler, c, nil))
	toJson(c, object)
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
	putObject(r.dbHandler, c, &dataModel.Patient{}, nil)
}

func (r *RestRouter) getNotes(c *gin.Context) {
	objects := getObjects[dataModel.Note](r.dbHandler, nil, nil)
	c.IndentedJSON(http.StatusOK, objects)
}

func (r *RestRouter) getNoteById(c *gin.Context) {
	object := dbLayer.NewDbObjectWrapper(getObjectById[dataModel.Note](r.dbHandler, c, nil))
	toJson(c, object)
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
	putObject(r.dbHandler, c, &dataModel.Note{}, nil)
}

func (r *RestRouter) getUsers(c *gin.Context) {
	objects := getObjects(r.dbHandler, nil, func(o *dataModel.User) {
		o.Password = "︻デ┳═ー"
		o.Salt = "0xdeadbeef"
	})
	c.IndentedJSON(http.StatusOK, objects)
}

func (r *RestRouter) getUserById(c *gin.Context) {
	object := dbLayer.NewDbObjectWrapper(getObjectById(r.dbHandler, c, func(o *dataModel.User) {
		o.Password = "=^..^="
		o.Salt = "0xFEEDFACE"
	}))
	toJson(c, object)
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
	putObject(r.dbHandler, c, &dataModel.User{}, func(user *dataModel.User) {
		user.Salt = cred.GenerateSalt()
		saltedPasswdSha256 := cred.CalcSha256(user.Password, user.Salt)
		user.Password = saltedPasswdSha256
	})
}

func (r *RestRouter) getManifests(c *gin.Context) {
	objects := getObjects[dataModel.PatientManifest](r.dbHandler, nil, nil)
	c.IndentedJSON(http.StatusOK, objects)
}

func (r *RestRouter) getManifestById(c *gin.Context) {
	object := dbLayer.NewDbObjectWrapper(getObjectById[dataModel.PatientManifest](r.dbHandler, c, nil))
	toJson(c, object)
}

func (r *RestRouter) postManifest(c *gin.Context) {
	newUser := dataModel.PatientManifest{}
	postObject(r.dbHandler, c, &newUser, nil)
}

func (r *RestRouter) deleteManifest(c *gin.Context) {
	deleteObject(r.dbHandler, c, &dataModel.PatientManifest{})
}

func (r *RestRouter) putManifest(c *gin.Context) {
	putObject(r.dbHandler, c, &dataModel.PatientManifest{}, nil)
}

func (r *RestRouter) GetEngine() *gin.Engine {
	return r.engine
}

func (r *RestRouter) Init(dbHandler *sqlx.DB, secret string) *RestRouter {
	r.dbHandler = dbHandler
	r.secret = secret

	// https: //go.dev/doc/tutorial/web-service-gin
	r.engine = gin.Default()

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.engine.MaxMultipartMemory = 8 << 20 // 8 MiB

	// CORS probably goes first as middleware
	// https://github.com/gin-contrib/cors
	r.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Origin", "Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowMethods:     []string{"GET", "DELETE", "POST", "PUT"},
		AllowCredentials: true,
		/* AllowOriginFunc: func(origin string) bool {
			return origin == "localhost"
		},
		*/
		MaxAge: 12 * time.Hour,
	}))
	r.engine.Use(r.checkToken())

	r.engine.POST(resetURL, r.resetDB)
	r.engine.POST(initURL, r.initDB)

	r.engine.GET(apiURL, r.getApi)
	r.engine.POST(loginURL, r.login)

	r.engine.GET(patientsURL, r.getPatients)
	r.engine.GET(patientsByIdURL, r.getPatientById)
	r.engine.POST(patientsURL, r.postPatient)
	r.engine.DELETE(patientsByIdURL, r.deletePatient)
	r.engine.PUT(patientsByIdURL, r.putPatient)

	r.engine.GET(notesURL, r.getNotes)
	r.engine.GET(notesByIdURL, r.getNoteById)
	r.engine.POST(notesByIdUploadURL, r.uploadNote)
	r.engine.POST(notesURL, r.postNote)
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
