package commons

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type MDTimeMock struct {
	mock.Mock
}

func (m *MDTimeMock) GetNowUTC() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

type MDCacheMock struct {
	mock.Mock
}

func (m *MDCacheMock) Get(key string) (string, error) {
	args := m.Called(key)
	return args.Get(0).(string), args.Error(1)
}

func (m *MDCacheMock) Set(key string, value string, seconds int) error {
	args := m.Called(key, value, seconds)
	return args.Error(0)
}

type MDSQSRepositoryMock struct {
	mock.Mock
}

func (m *MDSQSRepositoryMock) SendMessage(messageBody string) (string, error) {
	args := m.Called(messageBody)
	return args.Get(0).(string), args.Error(1)
}
