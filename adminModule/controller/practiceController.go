package controller

import (
	"common/results"
	"github.com/gin-gonic/gin"
	"user/service"
)

// PracticeTree 练习树(难度 -> 试卷 -> 题目)
func PracticeTree(c *gin.Context) {
	results.ResultHandle(c, service.PracticeTree(c))
}

// SqlPassTotal sql作答正确的数据
func SqlPassTotal(c *gin.Context) {
	results.ResultHandle(c, service.SqlPassTotal(c))
}

// CreateLevelType 创建难度类型
func CreateLevelType(c *gin.Context) {
	results.ResultHandle(c, service.CreateLevelType(c))
}

// CreatePaper 创建试卷
func CreatePaper(c *gin.Context) {
	results.ResultHandle(c, service.CreatePaper(c))
}

// CreateQuestion 创建题目
func CreateQuestion(c *gin.Context) {
	results.ResultHandle(c, service.CreateQuestion(c))
}
