package usecase

import (
	"context"

	"github.com/22Fariz22/musiclab/config"
	"github.com/22Fariz22/musiclab/internal/lyrics"
	"github.com/22Fariz22/musiclab/internal/models"
	"github.com/22Fariz22/musiclab/pkg/logger"
)

type lyricsUseCase struct {
	cfg        *config.Config
	lyricsRepo lyrics.Repository
	logger     logger.Logger
}

func NewLyricsUseCase(cfg *config.Config, lyricsRepo lyrics.Repository, logger logger.Logger) lyrics.UseCase {
	return &lyricsUseCase{cfg: cfg, lyricsRepo: lyricsRepo, logger: logger}
}

// Ping check
func (u lyricsUseCase) Ping() error {
	u.logger.Debug("Call UseCase Ping()")

	err := u.lyricsRepo.Ping()
	if err != nil {
		u.logger.Debug("error in lyricsUseCase Ping()")
		return err
	}

	u.logger.Debug("Pong in lyricsUseCase Ping()")
	return nil
}

func (u lyricsUseCase) DeleteSongByGroupAndTrack(ctx context.Context, groupName string, trackName string) error {
	u.logger.Debugf("in usecase DeleteSongByGroupAndTrack. Deleting song. Group: %s, Track: %s", groupName, trackName)
	return u.lyricsRepo.DeleteSongByGroupAndTrack(ctx, groupName, trackName)
}

func (u lyricsUseCase)UpdateTrackByID(ctx context.Context, updateData models.UpdateTrackRequest) error{
	u.logger.Debugf("in usecase UpdateTrackByID() ID:%d", updateData)
	return u.lyricsRepo.UpdateTrackByID(ctx, updateData)
}
