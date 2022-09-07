package models

import (
	pq "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	Model
	Premissions pq.Int32Array `json:"premissions" gorm:"type:integer[]" swaggertype:"array,number"`
	Account     string        `json:"account"`
	Password    []byte        `json:"-"`
	TFACode     []byte        `json:"-"`
	TFAEnable   bool          `json:"-"`
}

func (admin *Admin) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	admin.Password = hashedPassword
}

func (admin *Admin) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(admin.Password, []byte(password))
}
