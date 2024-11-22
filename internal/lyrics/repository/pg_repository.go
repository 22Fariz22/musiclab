package repository

import (
	"github.com/22Fariz22/musiclab/internal/lyrics"
	"github.com/22Fariz22/musiclab/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type lyricsRepo struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewLyricsRepository(db *sqlx.DB, logger logger.Logger) lyrics.Repository {
	return &lyricsRepo{db: db, logger: logger}
}

// Ping check
func (r lyricsRepo) Ping() error {
	r.logger.Debug("Call repo Ping()")
	err := r.db.Ping()
	if err != nil {
		r.logger.Debug("error in repo Ping():", err)
		return err
	}
	r.logger.Debug("Pong in repo Ping()")
	return nil
}
