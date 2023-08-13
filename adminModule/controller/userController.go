package controller

import (
	"common/results"
	"github.com/gin-gonic/gin"
	"user/service"
)

// UserLoginMonth 用户登录 当月统计
func UserLoginMonth(c *gin.Context) {
	results.ResultHandle(c, service.UserLoginMonth(c))
}
