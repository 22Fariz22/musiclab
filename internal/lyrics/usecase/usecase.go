package usecase

import (
	"github.com/22Fariz22/musiclab/config"
	"github.com/22Fariz22/musiclab/internal/lyrics"
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
