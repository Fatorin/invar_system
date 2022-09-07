package models

import "github.com/shopspring/decimal"

type Token struct {
	Model
	BankID          uint            `json:"bank_id"`
	Status          byte            `json:"status"`
	ProductID       uint            `json:"product_id"`
	ContractChain   string          `json:"contract_chain"`
	ContractType    string          `json:"contract_type"`
	ContractAddress string          `json:"contract_address"`
	OwnID           uint            `json:"own_id"`
	Quantity        decimal.Decimal `json:"quantity" gorm:"type:numeric"`
}

const (
	Unlock = iota
	Stacking
	StackExpired
)
