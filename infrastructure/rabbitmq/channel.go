package rabbitmq

import (
	"SomersaultCloud/constant/mq"
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync/atomic"
	"time"
)

type MessageQueue interface {
	// ExchangeDeclare 创建交换机.
	ExchangeDeclare(name string, kind string) (err error)

	// Publish 发布消息.
	Publish(exchange, key string, body []byte) (err error)

	// PublishWithDelay 发布消息带TTL.
	PublishWithDelay(exchange, key string, body []byte, timer time.Duration) (err error)

	// QueueDeclare 创建队列.
	QueueDeclare(name string) (err error)

	// QueueDeclareDeadLetter 创建死信队列.
	QueueDeclareDeadLetter(name, exchange, key string) (err error)

	// QueueBind 绑定队列.
	QueueBind(name, key, exchange string) (err error)
}

// RabbitMqChannel amqp.RabbitMqChannel wapper
type RabbitMqChannel struct {
	*amqp.Channel
	closed int32
}

// ExchangeDeclare 创建交换机.
func (ch *RabbitMqChannel) ExchangeDeclare(name string, kind string) (err error) {
	return ch.Channel.ExchangeDeclare(name, kind, true, false, false, false, nil)
}

// Publish 发布消息.
func (ch *RabbitMqChannel) Publish(exchange, key string, body []byte) (err error) {
	_, err = ch.Channel.PublishWithDeferredConfirmWithContext(context.Background(), exchange, key, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: body})
	return err
}

// PublishWithDelay 发布消息 带TTL.
func (ch *RabbitMqChannel) PublishWithDelay(exchange, key string, body []byte, timer time.Duration) (err error) {
	_, err = ch.Channel.PublishWithDeferredConfirmWithContext(context.Background(), exchange, key, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: body, Expiration: fmt.Sprintf("%d", timer.Milliseconds())})
	return err
}

// QueueDeclare 创建队列.
func (ch *RabbitMqChannel) QueueDeclare(name string) (err error) {
	_, err = ch.Channel.QueueDeclare(name, true, false, false, false, nil)
	return
}

// QueueDeclareDeadLetter 创建死信队列.
func (ch *RabbitMqChannel) QueueDeclareDeadLetter(name, exchange, key string) (err error) {
	_, err = ch.Channel.QueueDeclare(name, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    exchange,
		"x-dead-letter-routing-key": key,
	})
	return
}

// QueueBind 绑定队列.
func (ch *RabbitMqChannel) QueueBind(name, key, exchange string) (err error) {
	return ch.Channel.QueueBind(name, key, exchange, false, nil)
}

// IsClosed indicate closed by developer
func (ch *RabbitMqChannel) IsClosed() bool {
	return atomic.LoadInt32(&ch.closed) == 1
}

// Close ensure closed flag set
func (ch *RabbitMqChannel) Close() error {
	if ch.IsClosed() {
		return amqp.ErrClosed
	}

	atomic.StoreInt32(&ch.closed, 1)

	return ch.Channel.Close()
}

// Consume wrap amqp.RabbitMqChannel.Consume, the returned delivery will end only when channel closed by developer
func (ch *RabbitMqChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	deliveries := make(chan amqp.Delivery)

	go func() {
		for {
			d, err := ch.Channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
			if err != nil {
				log.Printf("consume failed, err: %v", err)
				time.Sleep(mq.RabbitMqReconnectDelay * time.Second)
				continue
			}

			for msg := range d {
				deliveries <- msg
			}

			// sleep before IsClose call. closed flag may not set before sleep.
			time.Sleep(mq.RabbitMqReconnectDelay * time.Second)

			if ch.IsClosed() {
				break
			}
		}
	}()

	return deliveries, nil
}
