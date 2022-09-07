package models

import "github.com/shopspring/decimal"

type Bank struct {
	Model
	UserID    uint            `json:"user_id"`
	InVarCoin decimal.Decimal `json:"invar_coin" gorm:"type:numeric"`
	USDTCoin  decimal.Decimal `json:"usdt_coin" gorm:"type:numeric"`
	Tokens    []Token         `json:"tokens"`
}
