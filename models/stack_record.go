package models

import (
	"time"
)

type StackRecord struct {
	Model
	Serial             string              `json:"serial"`
	Status             byte                `json:"status"`
	UserID             uint                `json:"user_id"`
	StackID            uint                `json:"stack_id"`
	Quantity           uint                `json:"quantity"`
	NextGetProfitTime  time.Time           `json:"next_get_profit_time"`
	EndGetProfitTime   time.Time           `json:"end_get_profit_time"`
	StackProfitRecords []StackProfitRecord `json:"stack_profit_records"`
}
