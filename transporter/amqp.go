package transporter

import (
	"github.com/paperclicks/go-rabbitmq"
	"github.com/streadway/amqp"
)

type AMQPTransporter struct {
	RabbitMQ *rabbitmq.RabbitMQ
	Queue    string
}

func (t *AMQPTransporter) Write(data []byte) (int, error) {

	qInfo := rabbitmq.QueueInfo{
		Name:       t.Queue,
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}

	err := t.RabbitMQ.PublishAssert(qInfo, string(data),amqp.Table{})

	return len(data), err
}
