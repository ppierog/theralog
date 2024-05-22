package dbLayer

import (
	"database/sql"
	"fmt"
	"log"
	"theraLog/dataRepository/dataModel"

	"github.com/jmoiron/sqlx"
)

type DbTable interface {
	dataModel.Patient | dataModel.Note | dataModel.User | dataModel.PatientManifest
}

type DbOps interface {
	*dataModel.Patient | *dataModel.Note | *dataModel.User | *dataModel.PatientManifest

	Insert() string
	Update() string
	SetRowId(rowId int64)
	Init(rows *sql.Rows) error
	GetRowId() int64
	GetName() string
	TableName() string
}

func executeIf(condition bool, f func()) {
	if condition {
		f()
	}
}

type QryBuilder struct {
	Qry string
}

func (qryBuilder *QryBuilder) SelectFrom(tableName string) *QryBuilder {
	qryBuilder.Qry = fmt.Sprintf("SELECT rowid, * FROM  %s ", tableName)
	return qryBuilder
}

func (qryBuilder *QryBuilder) DeleteFrom(tableName string) *QryBuilder {
	qryBuilder.Qry = fmt.Sprintf("DELETE FROM %s", tableName)
	return qryBuilder
}

func (qryBuilder *QryBuilder) WhereName(name string) *QryBuilder {
	qryBuilder.Qry = fmt.Sprintf("%s WHERE name = '%s' ", qryBuilder.Qry, name)
	return qryBuilder
}

func (qryBuilder *QryBuilder) Where(name string) *QryBuilder {
	qryBuilder.Qry = fmt.Sprintf("%s WHERE %s ", qryBuilder.Qry, name)
	return qryBuilder
}

func (qryBuilder *QryBuilder) Is(name string) *QryBuilder {
	qryBuilder.Qry = fmt.Sprintf("%s = '%s'", qryBuilder.Qry, name)
	return qryBuilder
}

func (qryBuilder *QryBuilder) WhereRowId(rowId int64) *QryBuilder {
	qryBuilder.Qry = fmt.Sprintf("%s WHERE rowid = %d ", qryBuilder.Qry, rowId)
	return qryBuilder
}
func (qryBuilder QryBuilder) Get() *QryBuilder {
	return &qryBuilder
}

type DbLayer struct {
	dbHandler *sqlx.DB
}

func FindBy[T DbTable,
	PT interface {
		Init(rows *sql.Rows) error
		TableName() string
		*T
	}](handler *sqlx.DB, qry *QryBuilder) []T {

	var ret []T
	var t T

	initializer := func(t PT, rows *sql.Rows) error {
		return t.Init(rows)
	}

	tabeleName := func(t PT) string {
		return t.TableName()
	}

	selectQry := QryBuilder{}
	selectQry.SelectFrom(tabeleName(&t))
	finalQry := QryBuilder{Qry: selectQry.Qry}
	if nil != qry {
		finalQry.Qry += qry.Qry
	}

	rows, err := handler.Query(finalQry.Qry)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err = initializer(&t, rows)
		if err != nil {
			log.Fatal(err)
		}
		ret = append(ret, t)

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)

	}
	return ret

}

func FindByName[T DbTable, PT interface {
	Init(rows *sql.Rows) error
	TableName() string
	*T
}](handler *sqlx.DB, name string, t PT) bool {
	qry := QryBuilder{}

	ret := FindBy[T, PT](handler, qry.WhereName(name))
	if len(ret) > 0 {

		executeIf(len(ret) > 1,
			func() { log.Println("Warning : more than 1 elements, returning index 0") })

		*t = ret[0]

		return true
	}
	return false
}

func FindByRowId[T DbTable, PT interface {
	Init(rows *sql.Rows) error
	TableName() string
	*T
}](handler *sqlx.DB, rowId int64, t PT) bool {
	qry := QryBuilder{}

	ret := FindBy[T, PT](handler, qry.WhereRowId(rowId))
	if len(ret) > 0 {
		executeIf(len(ret) > 1,
			func() { log.Println("Warning : more than 1 elements, returning index 0") })
		*t = ret[0]
		return true
	}
	return false
}

func Insert[T DbOps](handler *sqlx.DB, obj T) {
	res, err := handler.Exec(obj.Insert())
	if err != nil {
		log.Fatal("Could not Exec Insert()", err)

	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Could not Insert()", err)

	}
	obj.SetRowId(id)

}

func Exec(handler *sqlx.DB, qry *QryBuilder) sql.Result {

	result, err := handler.Exec(qry.Qry)

	if err != nil {
		log.Fatal("Could not Exec ", qry.Qry, err)
	}
	return result
}

func Update[T DbOps](handler *sqlx.DB, obj T) {

	qry := QryBuilder{Qry: obj.Update()}

	Exec(handler, qry.Get())
}

func DeleteBy[T DbOps](handler *sqlx.DB, qry *QryBuilder, obj T) int64 {

	mainQry := QryBuilder{}
	mainQry.DeleteFrom(obj.TableName())
	if qry != nil {
		mainQry.Qry += qry.Qry
	}

	result := Exec(handler, &mainQry)
	obj.SetRowId(-1)
	num, err := result.RowsAffected()
	if err != nil {
		return 0
	}
	return num

}

func DeleteByName[T DbOps](handler *sqlx.DB, obj T) int64 {
	qry := QryBuilder{}
	return DeleteBy(handler, qry.WhereName(obj.GetName()), obj)
}

func DeleteByRowId[T DbOps](handler *sqlx.DB, obj T) int64 {
	qry := QryBuilder{}
	return DeleteBy(handler, qry.WhereRowId(obj.GetRowId()), obj)
}
