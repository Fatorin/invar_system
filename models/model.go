package models

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `json:"id" gorm:"primarykey" form:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" swaggerignore:"true"`
}
