package practiceModel

import "common/model"

// Practice 题目表模型
type Practice struct {
	model.SuperModel
	//题目名称
	Name string `gorm:"column:name" json:"name"`
	//题目描述
	Description string `gorm:"column:description" json:"description"`
	//预定义答案
	Answer string `gorm:"column:answer" json:"answer"`
	//试卷id
	PaperID string `gorm:"column:paper_id" json:"paper_id"`
}

func (Practice) TableName() string {
	return "practice"
}
