package cache

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type cacheMock struct {
	mock.Mock
}

func NewRedisCacheMock() *cacheMock {
	return &cacheMock{}
}

func (m *cacheMock) Get(ctx context.Context, key string) (string, error) {
	agrs := m.Called(ctx, key)
	return agrs.String(0), agrs.Error(1)
}

func (m *cacheMock) Set(ctx context.Context, key string, value string, expiration int64) error {
	agrs := m.Called(ctx, key, value, expiration)
	return agrs.Error(0)
}
