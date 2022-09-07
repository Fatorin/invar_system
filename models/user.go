package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Model
	Role         uint          `json:"role"`
	Status       byte          `json:"status"`
	Email        string        `json:"email" gorm:"not null;size:50;unique"`
	UserName     string        `json:"username" gorm:"not null;size:50"`
	Password     []byte        `json:"-"`
	Referrer     string        `json:"referrer" gorm:"size:50"`
	ReferrerCode string        `json:"referrer_code" gorm:"size:50;unique"`
	Comment      string        `json:"commet"`
	UserKYC      UserKYC       `json:"user_kyc"`
	TFACode      []byte        `json:"-"`
	TFAEnable    bool          `json:"-"`
	Orders       []Order       `json:"-"`
	WhiteLists   []WhiteList   `json:"-"`
	StackRecords []StackRecord `json:"-"`
	Bank         Bank          `json:"-"`
}

func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user.Password = hashedPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}

const (
	Registed = iota + 1
	Auditing
	AuditFailed
	AuditSuccess
	Disabled
)
