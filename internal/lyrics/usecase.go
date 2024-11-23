package lyrics

import "context"

type UseCase interface {
	// Search()
	// GetSongVerses()
	DeleteSongByGroupAndTrack(ctx context.Context, groupName string, trackName string) error
	// UpdateTrack()
	// CreateTrack()
	Ping() error
}
