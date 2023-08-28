package model

import (
	"github.com/shopspring/decimal"
	"time"
)

// Order 委托订单
type Order struct {
	ID string `json:"id" gorm:"id"` // 订单ID

	UID string `json:"uid" gorm:"uid"` // 用户ID

	GID string `json:"gid" gorm:"gid"` // 股票编号

	Type int `json:"type" gorm:"type"` // 订单类型，1：市价单，2：限价单

	Direction int `json:"direction" gorm:"direction"` //买卖方向，1：买入，2：卖出

	Count int `json:"count" gorm:"count"` //数量

	DealCount int `json:"dealCount" gorm:"deal_count"` //成交数量

	DealPrice decimal.Decimal `json:"dealPrice" gorm:"deal_price"` //成交价格

	Price decimal.Decimal `json:"price" gorm:"price"` //请求价格

	Status int `json:"status" gorm:"status"` //订单状态，1：未成交，2：部分成交，3：全部成交，4：手动撤单，5：自动撤单

	CreatedTime time.Time `json:"createdTime" gorm:"created_time"` //下单时间

	DealTime time.Time `json:"dealTime" gorm:"deal_time"` //成交时间

	CancelTime time.Time `json:"cancelTime" gorm:"cancel_time"` //撤单时间
}

// QueueOrder 队列订单
type QueueOrder struct {
	ID string `json:"id" gorm:"id"` // 订单ID

	UID string `json:"uid" gorm:"uid"` // 用户ID

	GID string `json:"gid" gorm:"gid"` // 股票编号
}
