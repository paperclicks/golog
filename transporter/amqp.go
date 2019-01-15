package transporter

import "github.com/paperclicks/go-rabbitmq"

type AMQPTransporter struct {
	RabbitMQ *rabbitmq.RabbitMQ
	Queue    string
}

// func New(rmq *rabbitmq.RabbitMQ, queue string) *AMQPTransporter {

// 	return &AMQPTransporter{RabbitMQ: rmq, Queue: queue}
// }

func (t AMQPTransporter) Write(data []byte) (int, error) {

	err := t.RabbitMQ.Publish(t.Queue, string(data))

	return len(data), err
}