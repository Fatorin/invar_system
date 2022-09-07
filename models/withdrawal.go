package models

import "github.com/shopspring/decimal"

type Withdrawal struct {
	Model
	UserID   uint            `json:"-"`
	Quantity decimal.Decimal `json:"quantity" gorm:"type:numeric"`
	Fee      decimal.Decimal `json:"fee" gorm:"type:numeric"`
	Status   uint            `json:"status"`
}
