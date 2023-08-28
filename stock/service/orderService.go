package service

import (
	"common/results"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"stock/config"
	"stock/model"
	"strconv"
	"time"
)

const (
	MarketPrice = 1 // 市价
	LimitPrice  = 2 // 限价
	Buy         = 1 // 买入
	Sell        = 2 // 卖出
)

// CTX 操作redis的 空上下文
var CTX = context.Background()

func PriceOrder(c *gin.Context) *results.JsonResult {

	createdTime := time.Now()

	var requestData model.Order
	err := c.ShouldBindJSON(&requestData)
	if err != nil {
		return results.Error("参数错误", requestData)
	}

	requestData.CreatedTime = createdTime

	queueKey := requestData.GID + strconv.Itoa(requestData.Direction)
	// 获取分布式锁
	redisLock, redisLockErr := config.RedisLockClient.AcquireLock(queueKey, 5*time.Second)
	if redisLockErr != nil {
		return results.Error("redis锁定失败", requestData)
	}

	//从redis获取当前股票的 下单队列 (买/卖 方向)
	redisOrderSlice, redisErr := config.RedisClient.Get(CTX, queueKey).Result()
	if redisErr != nil && redisErr.Error() != redis.Nil.Error() {
		return results.Error("redis获取失败", requestData)
	}

	// redis中没有该股票的下单队列, 创建
	var orderSlice *[]model.Order
	*orderSlice = append(*orderSlice, requestData)

	redisErr = config.RedisClient.Set(CTX, queueKey, redisOrderSlice, -1).Err()
	if redisErr != nil {
		return results.Error("redis创建失败", requestData)
	}

	_, unlockErr := redisLock.Unlock()
	if unlockErr != nil {
		return results.Error("redis解锁失败", requestData)
	}

	return results.Success("下单成功", requestData)
}

func CancelOrder(c *gin.Context) *results.JsonResult {
	return results.Success("撤单成功", nil)
}
