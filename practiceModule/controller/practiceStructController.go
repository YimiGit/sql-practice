package controller

import (
	"common/results"
	"github.com/gin-gonic/gin"
	"practice/service"
)

// QuestionList 题目列表 grpc调用admin模块
func QuestionList(c *gin.Context) {
	results.ResultHandle(c, service.QuestionList(c))
}

// TableStruct 练习表结构 grpc调用admin模块
func TableStruct(c *gin.Context) {
	results.ResultHandle(c, service.TableStruct(c))
}

// LevelList 难度等级列表 grpc调用admin模块
func LevelList(c *gin.Context) {
	results.ResultHandle(c, service.LevelList(c))
}

// PaperList 试卷列表 grpc调用admin模块
func PaperList(c *gin.Context) {
	results.ResultHandle(c, service.PaperList(c))
}

// CommitSQL 提交练习
func CommitSQL(c *gin.Context) {
	results.ResultHandle(c, service.CommitSQL(c))
}
