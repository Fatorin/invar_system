package models

import (
	pq "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Stack struct {
	Model
	ProductID           uint            `json:"product_id"`
	Profit              decimal.Decimal `json:"profit" gorm:"type:numeric"`
	ContractTimeBound   uint            `json:"contract_time_bound"`
	ProfitIntervalMonth uint            `json:"profit_interval_month"`
	StackPermissions    pq.Int32Array   `json:"stack_permissions" gorm:"type:integer[]" swaggertype:"array,number"`
}
