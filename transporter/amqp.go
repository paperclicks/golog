package transporter

import "github.com/paperclicks/go-rabbitmq"

type AMQPTransporter struct {
	RabbitMQ *rabbitmq.RabbitMQ
	Queue    string
}

func (t *AMQPTransporter) Write(data []byte) (int, error) {

	err := t.RabbitMQ.Publish(t.Queue, string(data))
	if err != nil {
		panic(err)
	}

	return len(data), err
}
