package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/repo"
	"github.com/stretchr/testify/mock"
)

type MessageRepoMock struct {
	mock.Mock
}

func (m *MessageRepoMock) SaveMessage(ctx context.Context, msg *models.Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *MessageRepoMock) ListLatest(ctx context.Context, uid uuid.UUID, conversationID uuid.UUID, limit int) ([]models.Message, error) {
	args := m.Called(ctx, uid, conversationID, limit)
	msgs, ok := args.Get(0).([]models.Message)
	if !ok {
		return nil, args.Error(1)
	}
	return msgs, args.Error(1)
}

func (m *MessageRepoMock) ListBefore(ctx context.Context, uid uuid.UUID, conversationID uuid.UUID, beforeSeq int64, limit int) ([]models.Message, error) {
	args := m.Called(ctx, uid, conversationID, beforeSeq, limit)
	msgs, ok := args.Get(0).([]models.Message)
	if !ok {
		return nil, args.Error(1)
	}
	return msgs, args.Error(1)
}

func (m *MessageRepoMock) ListAfter(ctx context.Context, conversationID uuid.UUID, afterSeq int64, limit int) ([]models.Message, error) {
	args := m.Called(ctx, conversationID, afterSeq, limit)
	msgs, ok := args.Get(0).([]models.Message)
	if !ok {
		return nil, args.Error(1)
	}
	return msgs, args.Error(1)
}

var _ repo.MessageRepo = (*MessageRepoMock)(nil)
