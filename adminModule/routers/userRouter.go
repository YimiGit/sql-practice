package routers

import (
	"github.com/gin-gonic/gin"
	"user/controller"
)

func userRouter(e *gin.Engine) {
	e.GET("/user/login/log", controller.UserLoginMonth)
}
