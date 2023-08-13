package routers

import (
	"common/config"
	"common/middleware"
	"github.com/gin-gonic/gin"
)

// AllRoutersInit 所有路由初始化
func AllRoutersInit(r *gin.Engine) {

	//鉴权
	r.Use(middleware.CheckToken(config.RedisClient))
	{
		practiceRouter(r)
	}
}
