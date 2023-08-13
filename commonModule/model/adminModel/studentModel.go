package adminModel

import "common/model"

// Student 学生表模型
type Student struct {
	model.SuperModel
	//姓名
	Name string `gorm:"column:name" json:"name"`
	//班级id
	ClassId string `gorm:"column:class_id" json:"classId"`
	//性别(1-男，2-女)
	Gender int `gorm:"column:gender" json:"gender"`
}

func (Student) TableName() string {
	return "student"
}
