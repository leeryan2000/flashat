package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/wire"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/mock"
)

func TestWorkerIndexFor(t *testing.T) {
	t.Run("deterministic for the same id", func(t *testing.T) {
		id := uuid.New().String()
		first := workerIndexFor(id, workerPoolSize)
		for i := 0; i < 100; i++ {
			if got := workerIndexFor(id, workerPoolSize); got != first {
				t.Fatalf("workerIndexFor(%q) changed between calls: %d then %d", id, first, got)
			}
		}
	})

	t.Run("always within pool bounds", func(t *testing.T) {
		for i := 0; i < 200; i++ {
			idx := workerIndexFor(uuid.New().String(), workerPoolSize)
			if idx < 0 || idx >= workerPoolSize {
				t.Fatalf("workerIndexFor returned out-of-range index %d for pool size %d", idx, workerPoolSize)
			}
		}
	})

	t.Run("pool size of 1 always routes to worker 0", func(t *testing.T) {
		if idx := workerIndexFor(uuid.New().String(), 1); idx != 0 {
			t.Fatalf("expected index 0 for pool size 1, got %d", idx)
		}
	})

	t.Run("spreads different ids across more than one worker", func(t *testing.T) {
		seen := map[int]bool{}
		for i := 0; i < 50; i++ {
			seen[workerIndexFor(uuid.New().String(), workerPoolSize)] = true
		}
		if len(seen) < 2 {
			t.Fatalf("expected 50 random conversation ids to land on more than one worker, got %d distinct index(es)", len(seen))
		}
	})
}

// fakeAcknowledger is the mock implementation amqp091-go's own docs point
// at for testing Delivery handlers — a zero-value amqp.Delivery has a nil
// Acknowledger and panics on Ack/Nack, so processDelivery tests need this.
type fakeAcknowledger struct {
	acked        []uint64
	nacked       []uint64
	nackRequeues []bool
}

func (f *fakeAcknowledger) Ack(tag uint64, multiple bool) error {
	f.acked = append(f.acked, tag)
	return nil
}

func (f *fakeAcknowledger) Nack(tag uint64, multiple, requeue bool) error {
	f.nacked = append(f.nacked, tag)
	f.nackRequeues = append(f.nackRequeues, requeue)
	return nil
}

func (f *fakeAcknowledger) Reject(tag uint64, requeue bool) error {
	return nil
}

func newTestRoutedDelivery(ack *fakeAcknowledger) routedDelivery {
	return routedDelivery{
		d: amqp.Delivery{
			Acknowledger: ack,
			DeliveryTag:  1,
		},
		env: wire.Msg{
			ConversationID: uuid.New().String(),
			FromUID:        uuid.New().String(),
			ServerMsgID:    uuid.New().String(),
		},
	}
}

func TestProcessDelivery(t *testing.T) {
	t.Run("succeeds on first attempt: acks once, SaveMessage called once", func(t *testing.T) {
		svc, _, _, msgRepo, convRepo := newTestMessageService()
		msgRepo.On("SaveMessage", mock.Anything, mock.AnythingOfType("*models.Message")).Return(nil)
		convRepo.On("ListParticipantByID", mock.Anything, mock.Anything).Return([]*models.ConversationParticipant{}, nil)

		w := &MessageWorker{Service: svc}
		ack := &fakeAcknowledger{}
		rd := newTestRoutedDelivery(ack)

		w.processDelivery(context.Background(), rd)

		msgRepo.AssertNumberOfCalls(t, "SaveMessage", 1)
		if len(ack.acked) != 1 {
			t.Fatalf("expected 1 ack, got %d", len(ack.acked))
		}
		if len(ack.nacked) != 0 {
			t.Fatalf("expected 0 nacks, got %d", len(ack.nacked))
		}
	})

	t.Run("recovers after transient failures within the retry budget", func(t *testing.T) {
		svc, _, _, msgRepo, convRepo := newTestMessageService()
		msgRepo.On("SaveMessage", mock.Anything, mock.AnythingOfType("*models.Message")).
			Return(errors.New("transient")).Twice()
		msgRepo.On("SaveMessage", mock.Anything, mock.AnythingOfType("*models.Message")).
			Return(nil)
		convRepo.On("ListParticipantByID", mock.Anything, mock.Anything).Return([]*models.ConversationParticipant{}, nil)

		w := &MessageWorker{Service: svc}
		ack := &fakeAcknowledger{}
		rd := newTestRoutedDelivery(ack)

		start := time.Now()
		w.processDelivery(context.Background(), rd)
		elapsed := time.Since(start)

		msgRepo.AssertNumberOfCalls(t, "SaveMessage", processRetryAttempts)
		if len(ack.acked) != 1 {
			t.Fatalf("expected 1 ack after eventual success, got %d", len(ack.acked))
		}
		if len(ack.nacked) != 0 {
			t.Fatalf("expected 0 nacks, got %d", len(ack.nacked))
		}
		if elapsed < 2*processRetryDelay {
			t.Fatalf("expected at least 2 retry delays (%v) between 3 attempts, only took %v", 2*processRetryDelay, elapsed)
		}
	})

	t.Run("gives up after exhausting retries: nacks without requeue", func(t *testing.T) {
		svc, _, _, msgRepo, _ := newTestMessageService()
		msgRepo.On("SaveMessage", mock.Anything, mock.AnythingOfType("*models.Message")).
			Return(errors.New("permanent"))

		w := &MessageWorker{Service: svc}
		ack := &fakeAcknowledger{}
		rd := newTestRoutedDelivery(ack)

		w.processDelivery(context.Background(), rd)

		msgRepo.AssertNumberOfCalls(t, "SaveMessage", processRetryAttempts)
		if len(ack.acked) != 0 {
			t.Fatalf("expected 0 acks, got %d", len(ack.acked))
		}
		if len(ack.nacked) != 1 {
			t.Fatalf("expected exactly 1 nack, got %d", len(ack.nacked))
		}
		if ack.nackRequeues[0] != false {
			t.Fatalf("expected nack with requeue=false, got requeue=%v", ack.nackRequeues[0])
		}
	})
}
