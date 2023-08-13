package adminModel

import "common/model"

type SchoolTree struct {
	model.SuperModel
	//名称
	Name string `gorm:"column:name" json:"name"`
	//类型(1-学校，2-年级，3-班级)
	Type int `gorm:"column:type" json:"type"`
	//父级id
	ParentId string `gorm:"column:parent_id" json:"parentId"`

	Children []*SchoolTree `gorm:"-" json:"children"`

	Students []*Student `gorm:"-" json:"students"`
}

func (SchoolTree) TableName() string {
	return "school_tree"
}
