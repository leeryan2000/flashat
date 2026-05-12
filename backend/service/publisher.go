package service

import "github.com/leeryan2000/flashat/wire"

// Publisher is the interface MessageService depends on for outbound message queuing.
// Any concrete implementation (RabbitMQ, Kafka, in-memory) satisfies this.
type Publisher interface {
	Publish(env *wire.Msg) error
}
