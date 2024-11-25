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
	// GetVerses(ctx context.Context, ID uint) error
	GetSongByID(ctx context.Context, id uint) (models.Song, error)
}
