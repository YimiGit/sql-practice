package service

import (
	"common/config"
	"common/model/adminModel"
	"common/results"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

func UserLoginMonth(c *gin.Context) *results.JsonResult {
	var loginLogs []adminModel.UserLoginTotal
	bytes, err := config.RedisClient.Get(c, "userLoginMonth").Bytes()
	if err != nil {
		return results.Fail("查询失败", err)
	}
	err = sonic.Unmarshal(bytes, &loginLogs)
	if err != nil {
		return results.Fail("查询失败", err)
	}
	return results.Success("查询成功", loginLogs)
}
