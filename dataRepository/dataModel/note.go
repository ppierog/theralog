package dataModel

import (
	"database/sql"
	"fmt"
)

type Note struct {
	RowId        int64  `json:"id"`
	Name         string `json:"name" binding:"required"`
	PatientRowId int64  `json:"patientRowId" binding:"required"`
	SessionDate  int    `json:"sessionDate" binding:"required"`
	NoteDate     int    `json:"noteDate"`
	FileName     string `json:"fileName"`
	IsCrypted    bool   `json:"isCrypted"`
}

const NoteTableName string = "Note"

func (p *Note) TableName() string {
	return NoteTableName
}

func (p *Note) Init(rows *sql.Rows) error {
	if rows != nil {
		err := rows.Scan(&p.RowId, &p.Name, &p.PatientRowId, &p.SessionDate, &p.NoteDate, &p.FileName, &p.IsCrypted)
		return err

	} else {
		return nil
	}
}

func (p *Note) GetName() string {
	return p.Name
}

func (p *Note) GetRowId() int64 {
	return p.RowId
}

func (p *Note) SetRowId(rowId int64) {
	p.RowId = rowId
}

func (p *Note) Insert() string {
	const INSERT_QRY = "INSERT INTO %s VALUES('%s',%d, %d, %d, '%s', '%t');"
	return fmt.Sprintf(INSERT_QRY, p.TableName(), p.Name, p.PatientRowId, p.SessionDate, p.NoteDate, p.FileName, p.IsCrypted)
}

func (p *Note) Update() string {
	panic("Not Implemented")
}
