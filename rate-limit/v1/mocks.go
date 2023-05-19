package main

import (
	"github.com/stretchr/testify/mock"
)

type CacheRepositoryMock struct {
	mock.Mock
}

func (m *CacheRepositoryMock) get(compositeKey string) (RateLimitCounter, error) {
	args := m.Called(compositeKey)
	return args.Get(0).(RateLimitCounter), args.Error(1)
}

func (m *CacheRepositoryMock) set(rlc RateLimitCounter) error {
	args := m.Called(rlc)
	return args.Error(0)
}

type VerifyUCMock struct {
	mock.Mock
}

func (m *VerifyUCMock) handle(typeNotification string, rateLimitCounter *RateLimitCounter) (bool, error) {
	args := m.Called(typeNotification, rateLimitCounter)
	return args.Get(0).(bool), args.Error(1)
}

type ConfigRepositoryMock struct {
	mock.Mock
}

func (c *ConfigRepositoryMock) getNotificationLimitsByType(typeNotification string) (int, int, error) {
	args := c.Called(typeNotification)
	return args.Get(0).(int), args.Get(1).(int), args.Error(2)
}
