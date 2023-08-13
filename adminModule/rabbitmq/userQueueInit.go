package rabbitmq

import (
	"common/config"
	model "common/model/adminModel"
	"common/utils"
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

// UserLoginTotalChannel 登录统计channel, 用于 消费 登录统计消息
var userLoginTotalChannel *amqp.Channel

// UserLoginTotalChannelPool 登录统计channel池, 用于 生产 登录统计消息
var UserLoginTotalChannelPool *utils.ChannelPool

// userQueueInit 初始化用户队列
func userQueueInit() {

	//创建channel
	channel, channelErr := config.RabbitConnect.Channel()
	userLoginTotalChannel = channel
	if channelErr != nil {
		panic("userChannel创建失败")
	}

	exchangeName := "userExchange"
	//创建exchange
	exchangeErr := userLoginTotalChannel.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
	if exchangeErr != nil {
		panic("userExchange创建失败")
	}

	//每日登录统计队列, 创建queue, 绑定exchange, 创建consumer
	userLoginQueue(exchangeName)

	//每日登录统计 生产者信道池
	pool, channelErr := utils.NewChannelPool(10, config.RabbitConnect)
	if channelErr != nil {
		panic("userLoginTotalChannelPool创建失败")
	}
	UserLoginTotalChannelPool = pool

}

// userLoginQueue 每日登录统计队列
func userLoginQueue(exchangeName string) {
	//创建queue
	//exclusive: 当将队列声明为独占时，只有声明该队列的连接可以使用它
	//noWait: 当noWait为true时，不等待服务器确认队列的创建。注意，如果队列已经存在，那么该方法会报错
	queue, queueErr := userLoginTotalChannel.QueueDeclare("userLoginQueue", true, false, false, false, nil)
	if queueErr != nil {
		panic("userLoginQueue创建失败")
	}

	//绑定exchange和queue
	queueErr = userLoginTotalChannel.QueueBind(queue.Name, "userLoginKey", exchangeName, false, nil)
	if queueErr != nil {
		panic("userLoginQueue绑定失败")
	}

	//创建consumer, 统计每日登录人数
	go func() {
		c := make(chan struct{}, 0)
		messages, err := userLoginTotalChannel.Consume("userLoginQueue", "", false, false, false, false, nil)
		if err != nil {
			panic("userLoginConsumer创建失败")
		}
		for msg := range messages {

			nx := config.RedisClient.SetNX(context.Background(), "userLogin-"+string(msg.Body), nil, time.Duration(utils.TodayRemainNanosecond()))
			if nx.Err() != nil {

				log.Println("userLoginConsumer存入redis失败", nx.Err())
				//重回队列
				err := msg.Nack(false, true)
				if err != nil {
					log.Println("userLoginConsumer消息拒绝失败", err)
				}
				continue
			}

			userLoginTotalIdKey := "userLogin-" + utils.TodayFormat("yyyyMMdd")
			//当天未登录过
			if nx.Val() {
				//当天统计数据的id
				userLoginTotalId, err := config.RedisClient.Get(context.Background(), userLoginTotalIdKey).Result()

				//当天第一次统计
				if err == redis.Nil {
					id, err := utils.DistributedID(1)
					if err != nil {
						log.Println("userLoginConsumer生成id失败", err)
						//重回队列
						err := msg.Nack(false, true)
						if err != nil {
							log.Println("userLoginConsumer消息拒绝失败", err)
						}
						continue
					}

					userLoginTotalIdValue := strconv.Itoa(int(id))
					result := config.DB.Model(&model.UserLoginTotal{}).Create(&model.UserLoginTotal{ID: userLoginTotalIdValue, CreatedAt: time.Now(), LoginTotal: 1, DateFlag: userLoginTotalIdKey})
					if result.Error != nil {
						log.Println("userLoginConsumer创建当天第一条统计数据失败", result.Error)
						//重回队列
						err := msg.Nack(false, true)
						if err != nil {
							log.Println("userLoginConsumer消息拒绝失败", err)
						}
						continue
					}
					//第一次统计数据的id存入redis
					config.RedisClient.Set(context.Background(), userLoginTotalIdKey, userLoginTotalIdValue, time.Duration(utils.TodayRemainNanosecond()))
					err = msg.Ack(false)
					if err != nil {
						log.Println("userLoginConsumer消息确认失败", err)
					}
					continue

				} else if err != nil {
					log.Println("userLoginConsumer获取当天统计数据的id失败", err)
					//重回队列
					err := msg.Nack(false, true)
					if err != nil {
						log.Println("userLoginConsumer消息拒绝失败", err)
					}
					continue
				}
				//当天已经统计过
				//当天登录人数+1
				config.DB.Model(&model.UserLoginTotal{}).Where("id = ?", userLoginTotalId).Update("login_total", gorm.Expr("login_total + ?", 1))
			}
			//当天已经登录过
			err = msg.Ack(false)
		}
		c <- struct{}{}
	}()
}
