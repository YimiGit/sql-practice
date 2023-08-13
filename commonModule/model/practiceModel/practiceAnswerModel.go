package practiceModel

import (
	"common/model"
	"gorm.io/gorm"
	"strconv"
)

// UserCommitAnswerLog 答题记录表模型
type UserCommitAnswerLog struct {
	model.SuperModel
	QuestionId string  `gorm:"column:question_id" json:"questionId"`
	UserId     string  `gorm:"column:user_id" json:"userId"`
	SQLRunTime float64 `gorm:"column:sql_run_time" json:"sqlRunTime"`
	Type       int     `gorm:"-" json:"-"`
	AnswerSql  string  `gorm:"-" json:"answerSql"`
}

// UserCommitAnswerLogTableName 水平分表 动态表名
func UserCommitAnswerLogTableName(ucal *UserCommitAnswerLog) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Table("user_commit_answer_log_" + strconv.Itoa(ucal.Type))
	}
}
