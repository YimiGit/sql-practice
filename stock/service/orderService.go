package service

import (
	"common/results"
	"context"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"sort"
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
		return results.Fail("redis锁定失败", redisLockErr)
	}

	var orderSlice *[]*model.Order

	//从redis获取当前股票的 下单队列 (买/卖 方向)
	redisOrderSlice, redisErr := config.RedisClient.Get(CTX, queueKey).Result()
	if redisErr != nil {
		if redisErr.Error() != redis.Nil.Error() {
			return results.Fail("redis获取失败", redisErr)
		}
		// redis中没有该股票的下单队列, 创建
		orderSlice = &[]*model.Order{}
	} else {
		// redis中有该股票的下单队列, 反序列化
		redisErr = sonic.Unmarshal([]byte(redisOrderSlice), orderSlice)
		if redisErr != nil {
			return results.Fail("redis反序列化失败", redisErr)
		}
	}

	// 下单队列 排序
	orderSlice = sortOrder(orderSlice, requestData.Direction)

	//撮合交易
	//集合竞价

	redisErr = config.RedisClient.Set(CTX, queueKey, orderSlice, -1).Err()
	if redisErr != nil {
		return results.Fail("redis创建失败", redisErr)
	}

	_, unlockErr := redisLock.Unlock()
	if unlockErr != nil {
		return results.Fail("redis解锁失败", unlockErr)
	}

	return results.Success("下单成功", requestData)
}

// 下单队列 排序
// 买入: 价格从高到低、时间从早到晚
// 卖出: 价格从低到高、时间从早到晚
func sortOrder(slice *[]*model.Order, direction int) *[]*model.Order {
	slices := *slice
	if direction == Buy {
		sort.SliceStable(slices, func(i, j int) bool {
			cmp := slices[i].Price.Cmp(slices[j].Price)
			if cmp == 1 {
				return true
			}
			if cmp == -1 {
				return false
			}
			return slices[i].CreatedTime.Before(slices[j].CreatedTime)
		})
		return &slices
	}
	sort.SliceStable(slices, func(i, j int) bool {
		cmp := slices[i].Price.Cmp(slices[j].Price)
		if cmp == 1 {
			return false
		}
		if cmp == -1 {
			return true
		}
		return slices[i].CreatedTime.Before(slices[j].CreatedTime)
	})
	return &slices
}

func CancelOrder(c *gin.Context) *results.JsonResult {
	return results.Success("撤单成功", nil)
}

func KLineData(c *gin.Context) *results.JsonResult {
	return results.Success("获取成功", nil)
}
