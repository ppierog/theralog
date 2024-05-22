package dataModel

import (
	"database/sql"
	"fmt"
)

// Password :In the request post/put its plain password, in the response its sha256 from Salt + Password
// Salt 8 random bytes in a hex representation
type User struct {
	RowId           int64  `json:"id"`
	Name            string `json:"name" binding:"required"`
	LastName        string `json:"lastName" binding:"required"`
	Email           string `json:"email" binding:"required"`
	TelephoneNumber string `json:"telephoneNumber"`
	Salt            string `json:"salt"`
	Password        string `json:"password"`
	PubKey          string `json:"pubKey"`
}

type UserCred struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

const UserTableName string = "User"

func (p *User) TableName() string {
	return UserTableName
}

func (p *User) Init(rows *sql.Rows) error {
	if rows != nil {
		err := rows.Scan(&p.RowId, &p.Name, &p.LastName, &p.Email, &p.TelephoneNumber, &p.Salt, &p.Password, &p.PubKey)
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
	return fmt.Sprintf(INSERT_QRY, p.TableName(), p.Name, p.LastName, p.Email, p.TelephoneNumber, p.Salt, p.Password, p.PubKey)
}

func (p *User) Update() string {
	panic("Not Implemented")
}
