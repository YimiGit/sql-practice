package service

import (
	"common/config"
	"common/model/practiceModel"
	"common/results"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
)

func PracticeTree(c *gin.Context) *results.JsonResult {

	wg := sync.WaitGroup{}
	wg.Add(3)

	//所有难度类型
	var practiceLevel []*practiceModel.PracticeLevel
	var typeErr error
	go func() {
		defer wg.Done()
		typeResult := config.DB.Find(&practiceLevel)
		typeErr = typeResult.Error
	}()

	//所有试卷
	var practicePaper []*practiceModel.PracticePaper
	var paperErr error
	go func() {
		defer wg.Done()
		paperResult := config.DB.Find(&practicePaper)
		paperErr = paperResult.Error
	}()

	//所有题目
	var practiceQuestion []*practiceModel.Practice
	var questionErr error
	go func() {
		defer wg.Done()
		questionResult := config.DB.Find(&practiceQuestion)
		questionErr = questionResult.Error
	}()

	//阻塞主携程
	wg.Wait()

	if typeErr != nil || paperErr != nil || questionErr != nil {
		log.Println("查询失败", typeErr, paperErr, questionErr)
		return results.Fail("查询失败", typeErr)
	}

	//组装树形结构

	//第二层 试卷->题目
	type TreePaper struct {
		*practiceModel.PracticePaper
		Children []*practiceModel.Practice `json:"children"`
	}

	//第一层 难度类型->试卷
	type TreeType struct {
		*practiceModel.PracticeLevel
		Children []*TreePaper `json:"children"`
	}

	//题目分组
	questionGroupMap := make(map[string][]*practiceModel.Practice)
	for _, question := range practiceQuestion {
		if _, ok := questionGroupMap[question.PaperID]; !ok {
			questionGroupMap[question.PaperID] = make([]*practiceModel.Practice, 0)
		}
		questionGroupMap[question.PaperID] = append(questionGroupMap[question.PaperID], question)
	}

	//试卷分组 + 组装第二层
	var treePaperList []*TreePaper
	paperGroupMap := make(map[int][]*practiceModel.PracticePaper)
	for _, paper := range practicePaper {
		if _, ok := paperGroupMap[paper.Type]; !ok {
			paperGroupMap[paper.Type] = make([]*practiceModel.PracticePaper, 0)
		}
		paperGroupMap[paper.Type] = append(paperGroupMap[paper.Type], paper)
		treePaperList = append(treePaperList, &TreePaper{paper, questionGroupMap[paper.ID]})
	}

	//第二层分组
	treePaperGroup := make(map[int][]*TreePaper)
	for _, treePaper := range treePaperList {
		if _, ok := treePaperGroup[treePaper.Type]; !ok {
			treePaperGroup[treePaper.Type] = make([]*TreePaper, 0)
		}
		treePaperGroup[treePaper.Type] = append(treePaperGroup[treePaper.Type], treePaper)
	}

	//组装第一层
	var practiceLevelList []*TreeType
	for _, level := range practiceLevel {
		practiceLevelList = append(practiceLevelList, &TreeType{level, treePaperGroup[level.Type]})
	}

	return results.Success("success", practiceLevelList)
}

func SqlPassTotal(c *gin.Context) *results.JsonResult {
	type sqlPassTotals struct {
		QuestionId string `json:"question_id"`
		Total      int    `json:"total"`
	}
	var sqlPassTotal []*sqlPassTotals
	bytes, err := config.RedisClient.Get(c, "sqlPassTotal").Bytes()
	if err != nil {
		return results.Fail("查询失败", err)
	}
	err = sonic.Unmarshal(bytes, &sqlPassTotal)
	if err != nil {
		return results.Fail("查询失败", err)
	}
	return results.Success("success", sqlPassTotal)
}

func CreateLevelType(c *gin.Context) *results.JsonResult {
	var levelType practiceModel.PracticeLevel

	err := c.ShouldBindJSON(&levelType)
	if err != nil {
		return results.Fail("参数错误", err)
	}

	err = config.DB.Create(&levelType).Error
	if err != nil {
		return results.Fail("创建失败", err)
	}
	return results.Success("success", levelType)
}

func CreatePaper(c *gin.Context) *results.JsonResult {
	var paper practiceModel.PracticePaper
	err := c.ShouldBindJSON(&paper)
	if err != nil {
		return results.Fail("参数错误", err)
	}
	err = config.DB.Create(&paper).Error
	if err != nil {
		return results.Fail("创建失败", err)
	}
	return results.Success("success", paper)
}

func CreateQuestion(c *gin.Context) *results.JsonResult {
	var question practiceModel.Practice
	err := c.ShouldBindJSON(&question)
	if err != nil {
		return results.Fail("参数错误", err)
	}
	err = config.DB.Create(&question).Error
	if err != nil {
		return results.Fail("创建失败", err)
	}
	return results.Success("success", question)
}
