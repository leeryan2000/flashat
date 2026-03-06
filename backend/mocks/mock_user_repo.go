package mocks

import (
	"github.com/google/uuid"
	"github.com/leeryan2000/flashat/models"
	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *UserRepoMock) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *UserRepoMock) GetUserByUID(uid uuid.UUID) (*models.User, error) {
	args := m.Called(uid)
	user, ok := args.Get(0).(*models.User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *UserRepoMock) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	user, ok := args.Get(0).(*models.User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *UserRepoMock) UpdateName(uid uuid.UUID, name string) error {
	args := m.Called(uid, name)
	return args.Error(0)
}
