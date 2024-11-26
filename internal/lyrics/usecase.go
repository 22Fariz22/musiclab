package lyrics

import (
	"context"

	"github.com/22Fariz22/musiclab/internal/models"
)

type UseCase interface {
	DeleteSongByGroupAndTrack(ctx context.Context, groupName string, trackName string) error
	UpdateTrackByID(ctx context.Context, updateData models.UpdateTrackRequest) error
	CreateTrack(ctx context.Context, song models.SongRequest) error
	Ping() error
	GetSongVerseByPage(ctx context.Context, id uint, page int) (string, error)
	GetLibrary(ctx context.Context, group, song, releaseDate string, page, limit int) ([]models.Song, int, error)
}
