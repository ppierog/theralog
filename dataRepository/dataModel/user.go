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

type UserToken struct {
	Token string `json:"token" binding:"required"`
}

const UserTableName string = "User"

func (e *User) TableName() string {
	return UserTableName
}

func (e *User) Init(rows *sql.Rows) error {
	if rows != nil {
		err := rows.Scan(&e.RowId, &e.Name, &e.LastName, &e.Email, &e.TelephoneNumber, &e.Salt, &e.Password, &e.PubKey)
		return err
	} else {
		return nil
	}
}

func (e *User) GetName() string {
	return e.Name
}

func (e *User) GetRowId() int64 {
	return e.RowId
}

func (e *User) SetRowId(rowId int64) {
	e.RowId = rowId
}

func (e *User) Insert() string {
	const INSERT_QRY = "INSERT INTO %s VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s');"
	return fmt.Sprintf(INSERT_QRY, e.TableName(), e.Name, e.LastName, e.Email, e.TelephoneNumber, e.Salt, e.Password, e.PubKey)
}

func (e *User) Update() string {
	const UPDATE_QRY = "UPDATE %s SET name='%s', last_name='%s', email='%s', telephone_number='%s', salt='%s', password='%s' WHERE rowid=%d"
	return fmt.Sprintf(UPDATE_QRY, e.TableName(), e.Name, e.LastName, e.Email, e.TelephoneNumber, e.Salt, e.Password, e.RowId)
}
