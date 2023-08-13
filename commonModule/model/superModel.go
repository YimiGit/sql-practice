package model

import "time"

// SuperModel 公共模型
type SuperModel struct {
	ID        string    `gorm:"id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}
