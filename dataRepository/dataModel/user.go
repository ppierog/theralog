package dataModel

import (
	"database/sql"
	"fmt"
)

type User struct {
	RowId           int64  `json:"id"`
	Name            string `json:"name"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
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
	const INSERT_QRY = "INSERT INTO User VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s');"
	return fmt.Sprintf(INSERT_QRY, p.Name, p.LastName, p.Email, p.TelephoneNumber, p.PasswordSalt, p.PasswordSha, p.PubKey)
}

func (p *User) Update() string {
	return ""
}
