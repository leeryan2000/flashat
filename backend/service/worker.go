package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/leeryan2000/flashat/db"
	"github.com/leeryan2000/flashat/wire"
	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageWorker owns all RabbitMQ consumer mechanics.
// It knows nothing about business logic — it just dequeues,
// calls the service, and translates the result into ack/nack.
type MessageWorker struct {
	MQ      *db.RabbitMQClient
	Service *MessageService
}

func NewMessageWorker(mq *db.RabbitMQClient, svc *MessageService) *MessageWorker {
	return &MessageWorker{
		MQ:      mq,
		Service: svc,
	}
}

// Launch the consumer in a separate goroutine.
func (w *MessageWorker) Start() {
	deliveries, err := w.MQ.Ch.Consume(
		"chat_messages",
		"",    // consumer tag (auto-generated)
		false, // autoAck=false — we manually ACK after processing
		false, false, false, nil,
	)
	if err != nil {
		log.Fatal("worker: failed to start consumer:", err)
	}

	go func() {
		log.Println("worker: started, waiting for messages")
		for d := range deliveries {
			w.processDelivery(d)
		}
		log.Println("worker: delivery channel closed, shutting down")
	}()
}

func (w *MessageWorker) processDelivery(d amqp.Delivery) {
	ctx := context.Background()

	var env wire.Msg
	if err := json.Unmarshal(d.Body, &env); err != nil {
		log.Println("worker: ❌ malformed message, discarding:", err)
		d.Nack(false, false) // don't requeue — malformed will never succeed
		return
	}

	if err := w.Service.ProcessChat(ctx, &env); err != nil {
		log.Println("worker: ❌ ProcessChat failed, requeueing:", err)
		d.Nack(false, true) // requeue — could be a transient DB or network error
		return
	}

	d.Ack(false)
}
