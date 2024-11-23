package lyrics

import "context"

type Repository interface {
	Ping() error
	DeleteSongByGroupAndTrack(ctx context.Context, groupName string, trackName string) error
}
