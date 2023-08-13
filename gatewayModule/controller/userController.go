package controller

import (
	"common/results"
	"gateway/service"
	"github.com/gin-gonic/gin"
)

// Login 登录 gRPC调用Admin的登录
func Login(c *gin.Context) {
	results.ResultHandle(c, service.Login(c))
}

// Register 注册 gRPC调用Admin的注册
func Register(c *gin.Context) {
	results.ResultHandle(c, service.Register(c))
}
