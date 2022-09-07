package models

import (
	"github.com/shopspring/decimal"
)

type OrderItem struct {
	Model
	ProductID uint            `json:"product_id"`
	OrderID   uint            `json:"order_id"`
	Price     decimal.Decimal `json:"price" gorm:"type:numeric"`
	Quantity  uint            `json:"quantity"`
}
