package practiceModel

type UserCommitAnswerSql struct {
	ID        string `json:"id" gorm:"column:id"`
	CommitSQL string `json:"commitSql" gorm:"column:commit_sql"`
}

func (UserCommitAnswerSql) TableName() string {
	return "user_commit_answer_sql"
}
