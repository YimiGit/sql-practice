package routers

import (
	"common/config"
	"common/middleware"
	"github.com/gin-gonic/gin"
	"practice/controller"
)

// AllRoutersInit 初始化所有路由
func AllRoutersInit(e *gin.Engine) {

	e.GET("/practice/list", controller.QuestionList)
	e.GET("/practice/level/list", controller.LevelList)
	e.GET("/practice/paper/list", controller.PaperList)
	e.GET("/practice/table/struct", controller.TableStruct)
	e.POST("/practice/commit/sql", controller.CommitSQL)

	//鉴权
	e.Use(middleware.CheckToken(config.RedisClient))
	{

	}
}
