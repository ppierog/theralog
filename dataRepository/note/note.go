package note

import (
	"database/sql"
)

type Note struct {
	RowId        int64  `json:"id"`
	PatientRowId int64  `json:"patientRowId"`
	Name         string `json:"name"`
	SessionDate  int    `json:"sessionDate"`
	NoteDate     int    `json:"noteDate"`
	FileName     string `json:"fileName"`
}

const NoteTableName string = "Patient"

func (p *Note) TableName() string {
	return NoteTableName
}

func (n *Note) Init(rows *sql.Rows) error {
	if rows != nil {
		err := rows.Scan(&n.RowId, &n.PatientRowId, &n.Name, &n.SessionDate, &n.NoteDate, &n.FileName)
		return err
	} else {
		return nil
	}
}

func (s *Note) GetName() string {
	return s.Name
}

func (s *Note) GetRowId() int64 {
	return s.RowId
}

func (s *Note) SetRowId(rowId int64) {
	s.RowId = rowId
}

func (s *Note) Insert() string {
	const INSERT_QRY = "INSERT INTO Note VALUES('%s','%s',%d,'%s','%s');"
	//return fmt.Sprintf(INSERT_QRY, p.Name, p.Occupation, p.BirthYear, p.City, p.TelephoneNumber)
	return ""
}

func (n *Note) Update() string {
	return ""
}

func Init(n *Note, rows *sql.Rows) error {
	//return p.Init(rows)
	return nil
}
