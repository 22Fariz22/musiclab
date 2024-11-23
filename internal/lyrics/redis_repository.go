package lyrics

import (
	"time"
)

type SongRedisRepository interface {
	CacheSong(id string, text string, ttl time.Duration) error
	GetCachedSong(id string) (string, error)
}
