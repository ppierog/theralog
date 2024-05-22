package dataModel

import (
	"database/sql"
	"fmt"
)

type User struct {
	RowId           int64  `json:"id"`
	Name            string `json:"name" binding:"required"`
	LastName        string `json:"lastName" binding:"required"`
	Email           string `json:"email" binding:"required"`
	TelephoneNumber string `json:"telephoneNumber"`
	PasswordSalt    string `json:"passwordSalt"`
	PasswordSha     string `json:"passwordSha"`
	PubKey          string `json:"pubKey"`
}

const UserTableName string = "User"

func (p *User) TableName() string {
	return UserTableName
}

func (p *User) Init(rows *sql.Rows) error {
	if rows != nil {
		err := rows.Scan(&p.RowId, &p.Name, &p.LastName, &p.Email, &p.TelephoneNumber, &p.PasswordSalt, &p.PasswordSha, &p.PubKey)
		return err

	} else {
		return nil
	}
}

func (p *User) GetName() string {
	return p.Name
}

func (p *User) GetRowId() int64 {
	return p.RowId
}

func (p *User) SetRowId(rowId int64) {
	p.RowId = rowId
}

func (p *User) Insert() string {
	const INSERT_QRY = "INSERT INTO %s VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s');"
	return fmt.Sprintf(INSERT_QRY, p.TableName(), p.Name, p.LastName, p.Email, p.TelephoneNumber, p.PasswordSalt, p.PasswordSha, p.PubKey)
}

func (p *User) Update() string {
	panic("Not Implemented")
}
