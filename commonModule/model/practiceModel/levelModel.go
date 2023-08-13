package practiceModel

// PracticeLevel 试卷难度类型表模型
type PracticeLevel struct {
	ID   string `gorm:"column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
	Type int    `gorm:"column:type" json:"type"`
}

func (PracticeLevel) TableName() string {
	return "practice_level"
}
