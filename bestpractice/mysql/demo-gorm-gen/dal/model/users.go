package model

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Username string

type Password string

func (p *Password) Scan(src interface{}) error {
	*p = Password(fmt.Sprintf("@%v@", src))
	return nil
}

func (p *Password) Value() (driver.Value, error) {
	*p = Password(strings.Trim(string(*p), "@"))
	return p, nil
}

type User struct {
	gorm.Model        // ID uint CreatAt time.Time UpdateAt time.Time DeleteAt gorm.DeleteAt If it is repeated with the definition will be ignored
	ID         uint   `gorm:"primary_key" json:"id,omitempty"`
	Name       string `gorm:"column:name" json:"name,omitempty"`
	Age        int    `gorm:"column:age;type:varchar(64)" json:"age,omitempty"`
	Role       string `gorm:"column:role;type:varchar(64)" json:"role,omitempty"`
	Friends    []User `gorm:"-" json:"friends,omitempty"` // only local used gorm ignore
}

type Passport struct {
	ID        int
	Username  Username // will be field.String
	Password  Password // will be field.Field because type Password set Scan and Value
	LoginTime time.Time
}
