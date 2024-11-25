package lyrics

import (
	"context"

	"github.com/22Fariz22/musiclab/internal/models"
)

type UseCase interface {
	// Search()
	// GetVerses(ctx context.Context, ID uint) error
	DeleteSongByGroupAndTrack(ctx context.Context, groupName string, trackName string) error
	UpdateTrackByID(ctx context.Context, updateData models.UpdateTrackRequest) error
	CreateTrack(ctx context.Context, song models.SongRequest) error
	Ping() error
	GetSongVerseByPage(ctx context.Context, id uint, page int) (string, error)
}
