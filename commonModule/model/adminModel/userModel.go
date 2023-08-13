package adminModel

import "common/model"

// User 用户表模型
type User struct {
	model.SuperModel
	Username string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"password"`
}

func (User) TableName() string {
	return "user"
}
