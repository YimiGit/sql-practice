package practiceModel

// ProfileModel MySQL性能表模型
type ProfileModel struct {
	QueryID  string  `gorm:"column:Query_ID"`
	Duration float64 `gorm:"column:Duration"`
	Query    string  `gorm:"column:Query"`
}
