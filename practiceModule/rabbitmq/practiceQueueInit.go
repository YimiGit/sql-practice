package rabbitmq

import (
	"common/config"
	"github.com/streadway/amqp"
)

// PracticeAnswerChannel 练习答案统计 生产者channel
var PracticeAnswerChannel *amqp.Channel

func practiceQueueInit() {
	//创建channel
	channel, channelErr := config.RabbitConnect.Channel()
	if channelErr != nil {
		panic("practiceChannel创建失败")
	}
	PracticeAnswerChannel = channel
}
