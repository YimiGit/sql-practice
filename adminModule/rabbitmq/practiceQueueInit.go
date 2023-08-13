package rabbitmq

import (
	"common/config"
	"common/model/practiceModel"
	"common/utils"
	"github.com/bytedance/sonic"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
	"log"
	"strconv"
)

// 练习答案统计 消费者channel
var practiceAnswerChannel *amqp.Channel

// 练习答案统计队列初始化
func practiceQueueInit() {
	//创建channel
	channel, channelErr := config.RabbitConnect.Channel()
	practiceAnswerChannel = channel
	if channelErr != nil {
		panic("practiceChannel创建失败")
	}

	exchangeName := "practiceExchange"
	//创建exchange
	exchangeErr := practiceAnswerChannel.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
	if exchangeErr != nil {
		panic("practiceExchange创建失败")
	}

	//练习答案统计 队列, 创建queue, 绑定exchange, 创建consumer
	practiceAnswerQueue(exchangeName)
}

func practiceAnswerQueue(exchangeName string) {
	//创建queue
	queue, queueErr := practiceAnswerChannel.QueueDeclare("practiceAnswerQueue", true, false, false, false, nil)
	if queueErr != nil {
		panic("practiceAnswerQueue创建失败")
	}

	//绑定exchange和queue
	bindErr := practiceAnswerChannel.QueueBind(queue.Name, "practiceAnswerKey", exchangeName, false, nil)
	if bindErr != nil {
		panic("practiceAnswerQueue绑定失败")
	}

	//消费消息
	go func() {
		c := make(chan struct{}, 0)
		//创建consumer
		messages, consumerErr := practiceAnswerChannel.Consume(queue.Name, "", false, false, false, false, nil)
		if consumerErr != nil {
			panic("practiceAnswerConsumer创建失败")
		}
		for msg := range messages {
			var answer practiceModel.UserCommitAnswerLog
			err := sonic.Unmarshal(msg.Body, &answer)
			if err != nil {
				log.Println("practiceAnswerConsumer消息解析失败")
				msg.Nack(false, true)
				continue
			}
			//是否已有正确作答记录

			//查询题目难度类型, 确定水平分表
			var levelType int
			config.DB.Raw("select type from question_paper t1, practice t2 where t1.id = t2.paper_id and t2.id = ?", answer.QuestionId).Scan(&levelType)
			answer.Type = levelType

			var correctAnswer practiceModel.UserCommitAnswerLog
			result := config.DB.Scopes(practiceModel.UserCommitAnswerLogTableName(&answer)).First(&correctAnswer, "user_id = ? and question_id = ?", answer.UserId, answer.QuestionId)

			if result.Error != nil {
				if result.Error == gorm.ErrRecordNotFound {

					id, err := utils.DistributedID(1)
					if err != nil {
						log.Println("practiceAnswerConsumer生成id失败")
						msg.Nack(false, true)
						continue
					}

					//没有正确作答记录, 插入记录
					answer.ID = strconv.Itoa(int(id))
					begin := config.DB.Begin()
					result := begin.Scopes(practiceModel.UserCommitAnswerLogTableName(&answer)).Create(&answer)
					//垂直分表
					result2 := begin.Create(&practiceModel.UserCommitAnswerSql{ID: answer.ID, CommitSQL: answer.AnswerSql})
					if result.Error != nil || result2.Error != nil {
						log.Println("practiceAnswerConsumer插入失败", result.Error, result2.Error)
						begin.Rollback()
						msg.Nack(false, true)
						continue
					} else {
						begin.Commit()
						msg.Ack(false)
						continue
					}
				} else {
					log.Println("practiceAnswerConsumer查询失败", result.Error)
					msg.Nack(false, true)
					continue
				}
			}

			//已有正确作答记录, 且新答案更优, 更新记录
			if answer.SQLRunTime < correctAnswer.SQLRunTime {
				var answerForUpdate practiceModel.UserCommitAnswerLog
				answerForUpdate.SQLRunTime = answer.SQLRunTime
				answerForUpdate.AnswerSql = answer.AnswerSql
				result := config.DB.Scopes(practiceModel.UserCommitAnswerLogTableName(&answer)).Where("id = ?", correctAnswer.ID).Updates(&answerForUpdate)
				if result.Error != nil {
					msg.Nack(false, true)
					log.Println("practiceAnswerConsumer更新失败", result.Error)
					continue
				}
			}
			msg.Ack(false)
			log.Println("practiceAnswerConsumer收到消息: ", answer)
		}
		c <- struct{}{}
	}()
}
