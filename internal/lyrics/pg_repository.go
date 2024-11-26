package lyrics

import (
	"context"

	"github.com/22Fariz22/musiclab/internal/models"
)

type Repository interface {
	Ping() error
	DeleteSongByGroupAndTrack(ctx context.Context, groupName string, trackName string) error
	UpdateTrackByID(ctx context.Context, updateData models.UpdateTrackRequest) error
	CreateTrack(ctx context.Context, song models.SongRequest, songDetail models.SongDetail) error
	GetSongByID(ctx context.Context, id uint) (models.Song, error)
	GetLibrary(ctx context.Context, group, song, releaseDate string, offset, limit int) ([]models.Song, int, error)
}
