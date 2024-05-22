package dataModel

import (
	"database/sql"
	"fmt"
)

type PatientManifest struct {
	RowId        int64  `json:"id"`
	PatientId    int64  `json:"patientId"`
	UserId       int64  `json:"userId"`
	CrudMask     int    `json:"crudMask"`
	EncryptedAes string `json:"encryptedAes"`
}

const PatientManifestTableName string = "PatientManifest"

func (p *PatientManifest) TableName() string {
	return PatientManifestTableName
}

func (p *PatientManifest) Init(rows *sql.Rows) error {
	if rows != nil {
		err := rows.Scan(&p.RowId, &p.PatientId, &p.UserId, &p.CrudMask, &p.EncryptedAes)
		return err

	} else {
		return nil
	}
}

func (p *PatientManifest) GetName() string {
	panic("Not Implemented")
}

func (p *PatientManifest) GetRowId() int64 {
	return p.RowId
}

func (p *PatientManifest) SetRowId(rowId int64) {
	p.RowId = rowId
}

func (p *PatientManifest) Insert() string {
	const INSERT_QRY = "INSERT INTO PatientManifest VALUES(%d, %d, %d, '%s');"
	return fmt.Sprintf(INSERT_QRY, p.PatientId, p.UserId, p.CrudMask, p.EncryptedAes)
}

func (p *PatientManifest) Update() string {
	panic("Not Implemented")
}
