package routers

import (
	"github.com/gin-gonic/gin"
	"user/controller"
)

func practiceRouter(e *gin.Engine) {
	e.GET("/practice/tree", controller.PracticeTree)
	e.GET("/practice/sqlPassTotal", controller.SqlPassTotal)

	e.POST("/practice/create/levelType", controller.CreateLevelType)
	e.POST("/practice/create/paper", controller.CreatePaper)
	e.POST("/practice/create/question", controller.CreateQuestion)
}
