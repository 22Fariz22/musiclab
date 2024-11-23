package lyrics

import (
	"context"

	"github.com/22Fariz22/musiclab/internal/models"
)

type UseCase interface {
	// Search()
	// GetSongVerses()
	DeleteSongByGroupAndTrack(ctx context.Context, groupName string, trackName string) error
	UpdateTrackByID(ctx context.Context, updateData models.UpdateTrackRequest) error
	// CreateTrack()
	Ping() error
}
