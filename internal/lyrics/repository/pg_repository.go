package repository

import (
	"context"
	"database/sql"

	"github.com/22Fariz22/musiclab/internal/lyrics"
	"github.com/22Fariz22/musiclab/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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

func (r lyricsRepo) DeleteSongByGroupAndTrack(ctx context.Context, groupName string, trackName string) error {
	r.logger.Debugf("In repo. Deleting song. Group: %s, Track: %s", groupName, trackName)

	query := "DELETE FROM songs WHERE group_name = $1 AND song_name = $2"

	result, err := r.db.ExecContext(ctx, query, groupName, trackName)
	if err != nil {
		r.logger.Debugf("error in repo r.db.ExecContext(): ", err)
		return errors.Wrap(err, "SongRepository.DeleteSongByGroupAndTrack.ExecContext")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Debugf("error in repo result.RowsAffected(): ", err)
		return errors.Wrap(err, "SongRepository.DeleteSongByGroupAndName.RowsAffected")
	}
	if rowsAffected == 0 {
		r.logger.Debugf("in repo rowsAffected == 0")
		return errors.Wrap(sql.ErrNoRows, "SongRepository.DeleteSongByGroupAndName.rowsAffected")
	}

	return nil
}
