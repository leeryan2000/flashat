package db

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/leeryan2000/flashat/config"
	"github.com/leeryan2000/flashat/wire"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	conn *amqp.Connection
	Ch   *amqp.Channel
}

func NewRabbitMQClient(cfg config.Configuration) (*RabbitMQClient, error) {
	addr := fmt.Sprintf("amqp://%s:%s@%s/", cfg.MQ_USER, cfg.MQ_PASS, cfg.MQ_ADDR)
	conn, err := amqp.Dial(addr)
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

func (r *RabbitMQClient) Publish(env *wire.Msg) error {
	body, err := json.Marshal(env)
	if err != nil {
		return err
	}

	log.Printf("Publishing message to RabbitMQ: Type=%s FromUID=%s ConvID=%s\n", env.Type, env.FromUID, env.ConversationID)

	return r.Ch.Publish(
		"",              // default exchange
		"chat_messages", // routing key (match queue name, when using default exchange)
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
}

func (r *RabbitMQClient) Close() {
	if r == nil {
		return
	}
	r.Ch.Close()
	r.conn.Close()
}
