package adminModel

import "time"

// UserLoginTotal 用户登录统计表模型
type UserLoginTotal struct {
	ID        string    `gorm:"id" json:"id"`
	CreatedAt time.Time `gorm:"created_at" json:"createdAt"`
	//当天登录人数
	LoginTotal int `gorm:"login_total" json:"loginTotal"`
	//日期唯一标识
	DateFlag string `gorm:"date_flag" json:"dateFlag"`
}

func (UserLoginTotal) TableName() string {
	return "user_login_total"
}
