package mocks

import (
	"github.com/leeryan2000/flashat/wire"
	"github.com/stretchr/testify/mock"
)

type PublisherMock struct {
	mock.Mock
}

func (m *PublisherMock) Publish(env *wire.Msg) error {
	args := m.Called(env)
	return args.Error(0)
}
