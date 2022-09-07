package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	Model
	UserID           uint            `json:"user_id"`
	Serial           string          `json:"serial"`
	OrderItems       []OrderItem     `json:"order_items"`
	Status           byte            `json:"status"`
	TotalAmount      decimal.Decimal `json:"total_amount" gorm:"type:numeric"`
	TransactionChain string          `json:"transaction_chain"`
	TransactionID    string          `json:"transaction_id"`
	PaymentLimitTime time.Time       `json:"payment_limit_time"`
	Comment          string          `json:"comment"`
}

const (
	Cancel = iota + 1
	Ordered
	WaitConfirmPayment
	PaymentFailed
	Completed
)
