package bootstrap

import (
	"SomersaultCloud/app/constant/mq"
	"SomersaultCloud/app/infrastructure/log"
	"SomersaultCloud/app/infrastructure/rabbitmq"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

// NewRabbitConnection 获取channel.
func NewRabbitConnection(e *Env) *rabbitmq.Connection {
	if e.RabbitmqAddr == "" {
		return nil
	}
	defaultConn, err := Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		e.RabbitmqUser,
		e.RabbitmqPassword,
		e.RabbitmqAddr,
		e.RabbitmqPort))
	if err != nil {
		log.GetTextLogger().Error("new mq conn err: " + err.Error())
	}
	return defaultConn
}

// Dial wrap amqp.Dial, dial and get reconnect connection
func Dial(url string) (*rabbitmq.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	connection := &rabbitmq.Connection{
		Conn: conn,
	}

	go func() {
		for {
			reason, ok := <-connection.Conn.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok {
				log.GetTextLogger().Error("connection closed")
				break
			}
			log.GetJsonLogger().WithFields("mq connection closed reason", reason).Error("connection closed")
			// reconnect if not closed by developer
			for {
				// wait 1s for reconnect
				time.Sleep(mq.RabbitMqReconnectDelay * time.Second)

				conn, err := amqp.Dial(url)
				if err == nil {
					connection.Conn = conn
					log.GetTextLogger().Info("reconnect success")
					break
				}

				log.GetJsonLogger().WithFields("reconnected failed err", err).Error("mq reconnect")
			}
		}
	}()

	return connection, nil
}
