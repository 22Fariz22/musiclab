package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockDB - мок для базы данных
type MockDB struct {
	mock.Mock
}

func (m *MockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (int64, error) {
	arguments := m.Called(ctx, query, args)
	return arguments.Get(0).(int64), arguments.Error(1)
}