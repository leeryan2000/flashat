package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/models"
	"github.com/leeryan2000/flashat/repo"
	"github.com/leeryan2000/flashat/wire"
	"github.com/stretchr/testify/mock"
)

type ConversationRepoMock struct {
	mock.Mock
}

func (m *ConversationRepoMock) CreateGroupConversation(ctx context.Context, conv *models.Conversation, msg *models.Message, creatorUID uuid.UUID, participantsUID []uuid.UUID) error {
	args := m.Called(ctx, conv, msg, creatorUID, participantsUID)
	return args.Error(0)
}

func (m *ConversationRepoMock) CreateDirectConversation(ctx context.Context, conv *models.Conversation, uidA, uidB uuid.UUID) error {
	args := m.Called(ctx, conv, uidA, uidB)
	return args.Error(0)
}

func (m *ConversationRepoMock) ListConversationByUID(ctx context.Context, uid uuid.UUID) ([]*models.Conversation, error) {
	args := m.Called(ctx, uid)
	convs, ok := args.Get(0).([]*models.Conversation)
	if !ok {
		return nil, args.Error(1)
	}
	return convs, args.Error(1)
}

func (m *ConversationRepoMock) GetConversationByID(ctx context.Context, conversationID uuid.UUID) (*models.Conversation, error) {
	args := m.Called(ctx, conversationID)
	conv, ok := args.Get(0).(*models.Conversation)
	if !ok {
		return nil, args.Error(1)
	}
	return conv, args.Error(1)
}

func (m *ConversationRepoMock) ListParticipantByID(ctx context.Context, conversationID uuid.UUID) ([]*models.ConversationParticipant, error) {
	args := m.Called(ctx, conversationID)
	participants, ok := args.Get(0).([]*models.ConversationParticipant)
	if !ok {
		return nil, args.Error(1)
	}
	return participants, args.Error(1)
}

func (m *ConversationRepoMock) AddParticipant(ctx context.Context, conversationID uuid.UUID, participantUID uuid.UUID) error {
	args := m.Called(ctx, conversationID, participantUID)
	return args.Error(0)
}

func (m *ConversationRepoMock) ModifyParticipant(ctx context.Context, conversationID uuid.UUID, participantUID uuid.UUID, role string) error {
	args := m.Called(ctx, conversationID, participantUID, role)
	return args.Error(0)
}

func (m *ConversationRepoMock) RemoveParticipant(ctx context.Context, conversationID uuid.UUID, participantUID uuid.UUID) error {
	args := m.Called(ctx, conversationID, participantUID)
	return args.Error(0)
}

func (m *ConversationRepoMock) LeaveGroup(ctx context.Context, conversationID uuid.UUID, uid uuid.UUID) error {
	args := m.Called(ctx, conversationID, uid)
	return args.Error(0)
}

func (m *ConversationRepoMock) UpdateLastReadSeq(ctx context.Context, conversationID uuid.UUID, uid uuid.UUID, seq int64) error {
	args := m.Called(ctx, conversationID, uid, seq)
	return args.Error(0)
}

func (m *ConversationRepoMock) GetLastReadSeq(ctx context.Context, conversationID uuid.UUID, uid uuid.UUID) (int64, error) {
	args := m.Called(ctx, conversationID, uid)
	seq, ok := args.Get(0).(int64)
	if !ok {
		return 0, args.Error(1)
	}
	return seq, args.Error(1)
}

func (m *ConversationRepoMock) GetSummary(ctx context.Context, uid uuid.UUID) ([]*wire.ConversationSummary, error) {
	args := m.Called(ctx, uid)
	summaries, ok := args.Get(0).([]*wire.ConversationSummary)
	if !ok {
		return nil, args.Error(1)
	}
	return summaries, args.Error(1)
}

var _ repo.ConversationRepo = (*ConversationRepoMock)(nil)
