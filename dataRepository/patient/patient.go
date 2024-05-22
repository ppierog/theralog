package patient

import (
	"database/sql"
	"fmt"
)

type Patient struct {
	RowId           int64  `json:"id"`
	Name            string `json:"name"`
	Occupation      string `json:"occupation"`
	BirthYear       int    `json:"birthYear"`
	City            string `json:"city"`
	TelephoneNumber string `json:"telephoneNumber"`
}

const PatientTableName string = "Patient"

func (p *Patient) TableName() string {
	return PatientTableName
}

func (p *Patient) Init(rows *sql.Rows) error {
	if rows != nil {
		err := rows.Scan(&p.RowId, &p.Name, &p.Occupation, &p.BirthYear, &p.City, &p.TelephoneNumber)
		return err
	} else {
		return nil
	}
}

func (p *Patient) GetName() string {
	return p.Name
}

func (p *Patient) GetRowId() int64 {
	return p.RowId
}

func (p *Patient) SetRowId(rowId int64) {
	p.RowId = rowId
}

func (p *Patient) Insert() string {
	const INSERT_QRY = "INSERT INTO Patient VALUES('%s','%s',%d,'%s','%s');"
	return fmt.Sprintf(INSERT_QRY, p.Name, p.Occupation, p.BirthYear, p.City, p.TelephoneNumber)
}

func (p *Patient) Update() string {
	const UPDATE_QRY = "UPDATE Patient SET name='%s', occupation='%s', birth_year=%d, city='%s', telephone_number='%s' WHERE rowid=%d"
	return fmt.Sprintf(UPDATE_QRY, p.Name, p.Occupation, p.BirthYear, p.City, p.TelephoneNumber, p.RowId)
}

func Init(p *Patient, rows *sql.Rows) error {
	return p.Init(rows)
}
