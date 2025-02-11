// Code generated by MockGen. DO NOT EDIT.
// Source: pg_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	models "github.com/22Fariz22/musiclab/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateTrack mocks base method.
func (m *MockRepository) CreateTrack(ctx context.Context, song models.SongRequest, songDetail models.SongDetail) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTrack", ctx, song, songDetail)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTrack indicates an expected call of CreateTrack.
func (mr *MockRepositoryMockRecorder) CreateTrack(ctx, song, songDetail interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTrack", reflect.TypeOf((*MockRepository)(nil).CreateTrack), ctx, song, songDetail)
}

// DeleteSongByGroupAndTrack mocks base method.
func (m *MockRepository) DeleteSongByGroupAndTrack(ctx context.Context, groupName, trackName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSongByGroupAndTrack", ctx, groupName, trackName)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSongByGroupAndTrack indicates an expected call of DeleteSongByGroupAndTrack.
func (mr *MockRepositoryMockRecorder) DeleteSongByGroupAndTrack(ctx, groupName, trackName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSongByGroupAndTrack", reflect.TypeOf((*MockRepository)(nil).DeleteSongByGroupAndTrack), ctx, groupName, trackName)
}

// GetLibrary mocks base method.
func (m *MockRepository) GetLibrary(ctx context.Context, group, song, releaseDate string, offset, limit int) ([]models.Song, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLibrary", ctx, group, song, releaseDate, offset, limit)
	ret0, _ := ret[0].([]models.Song)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetLibrary indicates an expected call of GetLibrary.
func (mr *MockRepositoryMockRecorder) GetLibrary(ctx, group, song, releaseDate, offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLibrary", reflect.TypeOf((*MockRepository)(nil).GetLibrary), ctx, group, song, releaseDate, offset, limit)
}

// GetSongByID mocks base method.
func (m *MockRepository) GetSongByID(ctx context.Context, id uint) (models.Song, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSongByID", ctx, id)
	ret0, _ := ret[0].(models.Song)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSongByID indicates an expected call of GetSongByID.
func (mr *MockRepositoryMockRecorder) GetSongByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSongByID", reflect.TypeOf((*MockRepository)(nil).GetSongByID), ctx, id)
}

// Ping mocks base method.
func (m *MockRepository) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockRepositoryMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockRepository)(nil).Ping))
}

// UpdateTrackByID mocks base method.
func (m *MockRepository) UpdateTrackByID(ctx context.Context, updateData models.UpdateTrackRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTrackByID", ctx, updateData)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTrackByID indicates an expected call of UpdateTrackByID.
func (mr *MockRepositoryMockRecorder) UpdateTrackByID(ctx, updateData interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTrackByID", reflect.TypeOf((*MockRepository)(nil).UpdateTrackByID), ctx, updateData)
}
