package models

import (
	pq "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Product struct {
	Model
	Stock           uint            `json:"stock" form:"stock"`
	Status          bool            `json:"status" form:"status"`
	Price           decimal.Decimal `json:"price" form:"price" gorm:"type:numeric"`
	Title           string          `json:"title" form:"title"`
	Description     string          `json:"description" form:"description"`
	Image           string          `json:"image" form:"image"`
	PreviewImage    string          `json:"preview_image" form:"preview_image"`
	ContractChain   string          `json:"contract_chain" form:"contract_chain"`
	ContractType    string          `json:"contract_type" form:"contract_type"`
	ContractAddress string          `json:"contract_address" form:"contract_address"`
	BuyPermissions  pq.Int32Array   `json:"buy_permissions" form:"buy_permissions" gorm:"type:integer[]" swaggertype:"array,number"`
}
