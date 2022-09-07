package models

import (
	"github.com/shopspring/decimal"
)

type StackProfitRecord struct {
	Model
	StackRecordID uint            `json:"stack_record_id"`
	Profit        decimal.Decimal `json:"profit" gorm:"type:numeric"`
	Comment       string          `json:"comment"`
}
