package middleware

import (
	"common/results"
	"common/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

// CheckToken 检查token中间件
func CheckToken(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		// 未传token
		if len(token) <= 0 {
			results.ResultHandle(c, results.Error("请登录", nil))
			//立即中断请求处理流程 并 停止后续中间件 和 处理程序的执行, 它不会传递请求给下一个中间件或处理程序
			c.Abort()
			return
		}

		userId, err := utils.ParseToken(token)
		if err != nil {
			// token解析失败
			results.ResultHandle(c, results.Fail("登录信息已失效", err))
			c.Abort()
			return
		}

		result, err := redisClient.Get(c, strconv.Itoa(int(userId))).Result()
		if err != nil {
			// redis中不存在
			results.ResultHandle(c, results.Fail("登录信息有误", err))
			c.Abort()
			return
		}

		//续期
		err = redisClient.Set(c, strconv.Itoa(int(userId)), result, time.Minute*30).Err()
		if err != nil {
			// redis中不存在
			results.ResultHandle(c, results.Fail("登录信息续期失败", err))
			c.Abort()
			return
		}

		//请求中存入userId，方便后续使用
		c.Set("userId", userId)

		//流转请求
		c.Next()
	}
}
