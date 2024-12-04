package rabbitmq

import (
	"SomersaultCloud/app/somersaultcloud-chat/constant/mq"
	"SomersaultCloud/app/somersaultcloud-chat/infrastructure/log"
	log2 "SomersaultCloud/app/somersaultcloud-common/log"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

// Connection amqp.Connection wrapper
type Connection struct {
	Conn *amqp.Connection
}

//TODO 这里的逻辑前后不太一致,获取channel本身是在这里获取的,但是消费者创建新channel所需setup中的args信息
//	所以需到外面创建

// NewChannel 获取channel.
func (c *Connection) NewChannel() (*RabbitMqChannel, error) {
	ch, err := c.Conn.Channel()
	if err != nil {
		return nil, err
	}

	channel := &RabbitMqChannel{
		Channel: ch,
	}

	go func() {
		for {
			reason, ok := <-channel.Channel.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok || channel.IsClosed() {
				log2.GetTextLogger().Error("channel closed")
				_ = channel.Close() // close again, ensure closed flag set when connection closed
				break
			}
			log2.GetJsonLogger().WithFields("channel closed reasons", reason).Error("channel closed")

			// reconnect if not closed by developer
			for {
				// wait 1s for connection reconnect
				time.Sleep(mq.RabbitMqReconnectDelay * time.Second)

				ch, err := c.Conn.Channel()
				if err == nil {
					log2.GetJsonLogger().Info("channel recreate success")
					channel.Channel = ch
					break
				}

				log2.GetJsonLogger().WithFields("channel recreate failed", err.Error()).Error("channel recreate")
			}
		}

	}()

	return channel, nil
}

// NewConsumer 实例化一个消费者, 会单独用一个channel.
func (c *Connection) NewConsumer(ch *RabbitMqChannel, queue string, handler func([]byte) error) error {

	deliveries, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume err: %v, queue: %s", err, queue)
	}

	for msg := range deliveries {
		err = handler(msg.Body)
		if err != nil {
			_ = msg.Reject(true)
			continue
		}
		_ = msg.Ack(false)
	}

	return nil
}
