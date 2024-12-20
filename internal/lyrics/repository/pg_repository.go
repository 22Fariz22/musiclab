package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/22Fariz22/musiclab/internal/lyrics"
	"github.com/22Fariz22/musiclab/internal/models"
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

func (r lyricsRepo)UpdateTrackByID(ctx context.Context, updateData models.UpdateTrackRequest) error{
	r.logger.Debugf("in repo UpdateTrackByID() updateData:%+v\n", updateData)

	query := `
	UPDATE songs
	SET 
		group_name = COALESCE($1, group_name),
		song_name = COALESCE($2, song_name),
		release_date = COALESCE($3, release_date),
		text = COALESCE($4, text),
		link = COALESCE($5, link),
		updated_at = NOW()
		WHERE id = $6
		`

	result, err := r.db.ExecContext(
		ctx,
		query,
		updateData.GroupName,  
		updateData.SongName,   
		updateData.ReleaseDate, 
		updateData.Text,        
		updateData.Link,        
		updateData.ID,                    
	)
	if err != nil {
		r.logger.Debugf("in repo UpdateTrackByID() r.db.ExecContext return error: ",err)
		return errors.Wrap(err, "LyricsRepository.UpdateTrackByID.ExecContext")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Debugf("in repo UpdateTrackByID() result.RowsAffected() return error: ",err)
		return errors.Wrap(err, "LyricsRepository.UpdateTrackByID.RowsAffected")
	}
	if rowsAffected == 0 {
		r.logger.Debug("in repo UpdateTrackByID() rowsAffected == 0")
		return errors.Wrap(sql.ErrNoRows, "LyricsRepository.UpdateTrackByID.rowsAffected")
	}

	r.logger.Debug("in repo UpdateTrackByID() return nil")
	return nil
}

func (r lyricsRepo)CreateTrack(ctx context.Context, song models.SongRequest, songDetail models.SongDetail) error{
	r.logger.Debugf("in repo CreateTrack() song:%+v, songDetail:%+v\n", song, songDetail)
	
	query := `
		INSERT INTO songs (group_name, song_name, release_date, text, link, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		song.Group,             
		song.Song,              
		songDetail.ReleaseDate, 
		songDetail.Text,        
		songDetail.Link,      
	)

	if err != nil {
		r.logger.Errorf("error in repo CreateTrack() in r.db.ExecContext: %v", err)
		return errors.Wrap(err, "lyricsRepo.CreateTrack.ExecContext")
	}
	
	r.logger.Debug("error in repo CreateTrack() return nil")
	return nil
}

//GetSongByID получаем песню по ID
func (r lyricsRepo) GetSongByID(ctx context.Context, id uint) (models.Song, error) {
	var song models.Song
	query := `SELECT text FROM songs WHERE id = $1`

	err := r.db.GetContext(ctx, &song, query, id)
	if err != nil {
		return models.Song{}, fmt.Errorf("failed to fetch song: %w", err)
	}

	return song, nil
}

func (r lyricsRepo) GetLibrary(ctx context.Context, group, song, releaseDate string, offset, limit int) ([]models.Song, int, error) {
    var songs []models.Song
    var total int

    baseQuery := `SELECT id, group_name, song_name, text, release_date, link FROM songs`
    baseCountQuery := `SELECT COUNT(*) FROM songs`
    conditions := []string{}
    args := []interface{}{}

    if group != "" {
        conditions = append(conditions, "group_name ILIKE $"+strconv.Itoa(len(args)+1))
        args = append(args, "%"+group+"%")
    }
    if song != "" {
        conditions = append(conditions, "song_name ILIKE $"+strconv.Itoa(len(args)+1))
        args = append(args, "%"+song+"%")
    }
    if releaseDate != "" {
        conditions = append(conditions, "release_date = $"+strconv.Itoa(len(args)+1))
        args = append(args, releaseDate)
    }

    if len(conditions) > 0 {
        conditionString := " WHERE " + strings.Join(conditions, " AND ")
        baseQuery += conditionString
        baseCountQuery += conditionString
    }

    baseQuery += " ORDER BY id LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
    args = append(args, limit, offset)

    if err := r.db.SelectContext(ctx, &songs, baseQuery, args...); err != nil {
        return nil, 0, fmt.Errorf("failed to fetch songs: %w", err)
    }

    if err := r.db.GetContext(ctx, &total, baseCountQuery, args[:len(args)-2]...); err != nil {
        return nil, 0, fmt.Errorf("failed to fetch total count: %w", err)
    }

    return songs, total, nil
}
