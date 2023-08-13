package practiceModel

// PracticePaper 试卷表模型
type PracticePaper struct {
	ID          string `gorm:"column:id" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	Type        int    `gorm:"column:type" json:"type"`
	TableStruct string `gorm:"column:table_struct" json:"tableStruct"`
}

func (PracticePaper) TableName() string {
	return "question_paper"
}
