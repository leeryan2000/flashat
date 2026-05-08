package db

import (
	"fmt"
	"github.com/leeryan2000/flashat/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	conn *amqp.Connection
	Ch   *amqp.Channel
}

func NewRabbitMQClient(cfg config.Configuration) (*RabbitMQClient, error) {
	conn, err := amqp.Dial("amqp://" + cfg.MQ_ADDR)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ at %s: %w", cfg.MQ_ADDR, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	_, err = ch.QueueDeclare("chat_messages", true, false, false, false, nil)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &RabbitMQClient{conn: conn, Ch: ch}, nil
}

func (r *RabbitMQClient) Close() {
	if r == nil {
		return
	}
	r.Ch.Close()
	r.conn.Close()
}
