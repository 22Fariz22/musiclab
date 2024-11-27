package mock

import (
	"context"

	"github.com/22Fariz22/musiclab/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUseCase struct {
	mock.Mock
}

func (m *MockUseCase) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUseCase) DeleteSongByGroupAndTrack(ctx context.Context, groupName, trackName string) error {
	args := m.Called(ctx, groupName, trackName)
	return args.Error(0)
}

func (m *MockUseCase) UpdateTrackByID(ctx context.Context, updateData models.UpdateTrackRequest) error {
	args := m.Called(ctx, updateData)
	return args.Error(0)
}

func (m *MockUseCase) CreateTrack(ctx context.Context, song models.SongRequest) error {
	args := m.Called(ctx, song)
	return args.Error(0)
}

func (m *MockUseCase) GetSongVerseByPage(ctx context.Context, id uint, page int) (string, error) {
	args := m.Called(ctx, id, page)
	return args.String(0), args.Error(1)
}

func (m *MockUseCase) GetLibrary(ctx context.Context, group, song, releaseDate string, page, limit int) ([]models.Song, int, error) {
	args := m.Called(ctx, group, song, releaseDate, page, limit)
	return args.Get(0).([]models.Song), args.Int(1), args.Error(2)
}
