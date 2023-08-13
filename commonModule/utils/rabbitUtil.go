package utils

import (
	"errors"
	"github.com/streadway/amqp"
	"log"
)

// PushMessage 发送消息
func PushMessage(message []byte, channel *amqp.Channel, exchangeName string, routingKey string) {
	err := channel.Publish(exchangeName, routingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	})
	if err != nil {
		log.Println("消息发送失败", err)
	}
}

// ChannelPool RabbitMQ信道池
type ChannelPool struct {
	pool     chan *amqp.Channel
	capacity int // 通道池的容量
}

// NewChannelPool 创建一个通道池
func NewChannelPool(capacity int, conn *amqp.Connection) (*ChannelPool, error) {
	chPool := &ChannelPool{
		pool:     make(chan *amqp.Channel, capacity),
		capacity: capacity,
	}

	for i := 0; i < capacity; i++ {
		channel, err := conn.Channel()
		if err != nil {
			return nil, err
		}
		chPool.pool <- channel
	}

	return chPool, nil
}

// GetChannel 从通道池中获取一个信道
func (cp *ChannelPool) GetChannel() (*amqp.Channel, error) {
	select {
	case ch := <-cp.pool:
		return ch, nil
	default:
		log.Println("登录统计channel pool is empty")
		return nil, errors.New("channel pool is empty")
	}
}

// ReleaseChannel 释放一个信道
func (cp *ChannelPool) ReleaseChannel(ch *amqp.Channel) {
	cp.pool <- ch
}
