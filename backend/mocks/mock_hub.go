package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type HubMock struct {
	mock.Mock
}

func (m *HubMock) SendToUID(uid uuid.UUID, payload []byte) {
	m.Called(uid, payload)
}

func (m *HubMock) BroadcastToParticipant(uids []uuid.UUID, fromUID uuid.UUID, payload []byte) {
	m.Called(uids, fromUID, payload)
}
