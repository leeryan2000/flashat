package service

import (
	"errors"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestDecideOutcome(t *testing.T) {
	someErr := errors.New("boom")

	tests := []struct {
		name       string
		processErr error
		retryCount int
		want       deliveryOutcome
	}{
		{"success acks regardless of retry count", nil, 0, outcomeAck},
		{"success acks even at the retry ceiling", nil, maxDeliveryRetries, outcomeAck},
		{"failure with no prior retries is retried", someErr, 0, outcomeRetry},
		{"failure just under the ceiling is retried", someErr, maxDeliveryRetries - 1, outcomeRetry},
		{"failure at the ceiling is discarded", someErr, maxDeliveryRetries, outcomeDiscard},
		{"failure past the ceiling is discarded", someErr, maxDeliveryRetries + 100, outcomeDiscard},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decideOutcome(tt.processErr, tt.retryCount)
			if got != tt.want {
				t.Errorf("decideOutcome(err=%v, retryCount=%d) = %v, want %v", tt.processErr, tt.retryCount, got, tt.want)
			}
		})
	}
}

func TestRetryCountFromHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers amqp.Table
		want    int
	}{
		{"nil headers default to zero", nil, 0},
		{"missing header defaults to zero", amqp.Table{}, 0},
		{"int32 as published by PublishRaw", amqp.Table{retryCountHeader: int32(3)}, 3},
		{"int64 as some AMQP brokers round-trip it", amqp.Table{retryCountHeader: int64(3)}, 3},
		{"plain int", amqp.Table{retryCountHeader: 3}, 3},
		{"unexpected type falls back to zero rather than panicking", amqp.Table{retryCountHeader: "3"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := retryCountFromHeaders(tt.headers)
			if got != tt.want {
				t.Errorf("retryCountFromHeaders(%v) = %d, want %d", tt.headers, got, tt.want)
			}
		})
	}
}
