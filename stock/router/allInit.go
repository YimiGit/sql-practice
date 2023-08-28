package router

import (
	"common/middleware"
	"github.com/gin-gonic/gin"
)

func AllInit(engine *gin.Engine) {

	//跨域中间件
	engine.Use(middleware.Cors())
	{
		stockInit(engine)
	}
}
