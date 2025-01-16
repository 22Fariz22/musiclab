package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

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

// DeleteSongByGroupAndTrack Удаление песни
func (r lyricsRepo) DeleteSongByID(ctx context.Context, ID uint) error {
	r.logger.Debugf("In repo. Deleting ID %d", ID)

	query := `
        DELETE FROM songs
        WHERE id = $1
    `

	result, err := r.db.ExecContext(ctx, query, ID)
	if err != nil {
		r.logger.Debugf("error in repo r.db.ExecContext(): %v", err)
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Debugf("error in repo result.RowsAffected(): %v", err)
		return fmt.Errorf("failed to fetch rows affected: %w", err)
	}
	if rowsAffected == 0 {
		r.logger.Debugf("no rows were deleted")
		return sql.ErrNoRows
	}

	return nil
}

func (r lyricsRepo) UpdateTrackByID(ctx context.Context, updateData models.UpdateTrackRequest) error {
	r.logger.Debugf("in repo UpdateTrackByID() updateData: %+v", updateData)

	// Проверяем, что GroupName и SongName присутствуют (на всякий случай)
	if updateData.GroupName == nil || updateData.SongName == nil {
		r.logger.Debug("GroupName or SongName is nil, but should not be")
		return errors.New("group name and song name cannot be nil")
	}

	// Проверяем существование или создаем группу
	var groupID uint
	err := r.db.QueryRowContext(ctx, "SELECT id FROM groups WHERE name = $1", *updateData.GroupName).Scan(&groupID)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Debugf("group not found, creating new group with name: %s", *updateData.GroupName)
			err = r.db.QueryRowContext(ctx, "INSERT INTO groups (name,created_at,updated_at) VALUES ($1,$2,$2) RETURNING id", *updateData.GroupName, time.Now()).Scan(&groupID)
			if err != nil {
				r.logger.Debugf("error creating new group: %v", err)
				return errors.Wrap(err, "failed to create new group")
			}
		} else {
			r.logger.Debugf("error fetching group ID for name %s: %v", *updateData.GroupName, err)
			return errors.Wrap(err, "error fetching group ID")
		}
	}

	// Начинаем формировать запрос для обновления песни
	query := "UPDATE songs SET updated_at = NOW(), group_id = $1, song_name = $2"
	params := []interface{}{groupID, *updateData.SongName}
	paramCount := 3

	// Добавляем опциональные поля
	if updateData.ReleaseDate != nil {
		query += fmt.Sprintf(", release_date = $%d", paramCount)
		params = append(params, *updateData.ReleaseDate)
		paramCount++
	}

	if updateData.Text != nil {
		query += fmt.Sprintf(", text = $%d", paramCount)
		params = append(params, *updateData.Text)
		paramCount++
	}

	if updateData.Link != nil {
		query += fmt.Sprintf(", link = $%d", paramCount)
		params = append(params, *updateData.Link)
		paramCount++
	}

	// Добавляем условие WHERE
	query += fmt.Sprintf(" WHERE id = $%d", paramCount)
	params = append(params, updateData.ID)

	// Выполняем запрос
	r.logger.Debugf("executing query: %s with params: %+v", query, params)
	result, err := r.db.ExecContext(ctx, query, params...)
	if err != nil {
		r.logger.Debugf("error in UpdateTrackByID() r.db.ExecContext: %v", err)
		return errors.Wrap(err, "LyricsRepository.UpdateTrackByID.ExecContext")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Debugf("error in UpdateTrackByID() result.RowsAffected(): %v", err)
		return errors.Wrap(err, "LyricsRepository.UpdateTrackByID.RowsAffected")
	}

	if rowsAffected == 0 {
		r.logger.Debug("no rows affected, song not found")
		return errors.Wrap(sql.ErrNoRows, "song not found")
	}

	r.logger.Debug("successfully updated track")
	return nil
}

func (r lyricsRepo) CreateTrack(ctx context.Context, songRequest models.SongRequest, songDetail models.SongDetail) error {
	// Начинаем транзакцию
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "lyricsRepo.CreateTrack.BeginTx")
	}
	defer tx.Rollback()

	// Получаем или создаем группу
	var groupID uint
	queryGroup := `
        WITH ins AS (
            INSERT INTO groups (name, created_at, updated_at)
            VALUES ($1, NOW(), NOW())
            ON CONFLICT (name) DO NOTHING
            RETURNING id
        )
        SELECT id FROM ins
        UNION ALL
        SELECT id FROM groups WHERE name = $1
        LIMIT 1
    `
	err = tx.QueryRowContext(ctx, queryGroup, songRequest.Group).Scan(&groupID)
	if err != nil {
		r.logger.Errorf("error getting/creating group: %v", err)
		return errors.Wrap(err, "lyricsRepo.CreateTrack.QueryGroup")
	}

	// Проверяем существование песни у этой группы
	var exists bool
	queryCheck := `
        SELECT EXISTS (
            SELECT 1 FROM songs 
            WHERE group_id = $1 AND song_name = $2
        )
    `
	err = tx.QueryRowContext(ctx, queryCheck, groupID, songRequest.Song).Scan(&exists)
	if err != nil {
		r.logger.Errorf("error checking song existence: %v", err)
		return errors.Wrap(err, "lyricsRepo.CreateTrack.CheckExistence")
	}

	if exists {
		return nil
	}

	// Добавляем песню
	queryInsert := `
        INSERT INTO songs (group_id, song_name, release_date, text, link, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
    `
	_, err = tx.ExecContext(
		ctx,
		queryInsert,
		groupID,
		songRequest.Song,
		songDetail.ReleaseDate,
		songDetail.Text,
		songDetail.Link,
	)
	if err != nil {
		r.logger.Errorf("error inserting song: %v", err)
		return errors.Wrap(err, "lyricsRepo.CreateTrack.InsertSong")
	}

	// Подтверждаем транзакцию
	if err = tx.Commit(); err != nil {
		r.logger.Errorf("error committing transaction: %v", err)
		return errors.Wrap(err, "lyricsRepo.CreateTrack.Commit")
	}

	r.logger.Debug("successfully created track")
	return nil
}

// GetSongByID получаем песню по ID
func (r lyricsRepo) GetSongByID(ctx context.Context, id uint) (models.Song, error) {
	var song models.Song
	query := `SELECT text FROM songs WHERE id = $1`

	err := r.db.GetContext(ctx, &song, query, id)
	if err != nil {
		return models.Song{}, fmt.Errorf("failed to fetch song: %w", err)
	}

	return song, nil
}

// GetLibrary Получение данных библиотеки с фильтрацией по всем полям и пагинацией
func (r lyricsRepo) GetLibrary(ctx context.Context, group, song, text, releaseDate string, offset, limit int) ([]models.Song, int, error) {
	var songs []models.Song
	var total int

	baseQuery := `SELECT s.id, s.group_id, g.name AS group_name, s.song_name, s.text, s.release_date, s.link
                  FROM songs s
                  INNER JOIN groups g ON s.group_id = g.id`
	baseCountQuery := `SELECT COUNT(*) FROM songs s INNER JOIN groups g ON s.group_id = g.id`
	conditions := []string{}
	args := []interface{}{}

	if group != "" {
		conditions = append(conditions, "g.name ILIKE $"+strconv.Itoa(len(args)+1)) // Фильтрация по названию группы
		args = append(args, "%"+group+"%")
	}
	if song != "" {
		conditions = append(conditions, "s.song_name ILIKE $"+strconv.Itoa(len(args)+1)) // Фильтрация по названию песни
		args = append(args, "%"+song+"%")
	}
	if releaseDate != "" {
		conditions = append(conditions, "s.release_date = $"+strconv.Itoa(len(args)+1)) // Фильтрация по дате релиза
		args = append(args, releaseDate)
	}
	if text != "" {
		conditions = append(conditions, "s.text ILIKE $"+strconv.Itoa(len(args)+1)) // Фильтрация по тексту песни
		args = append(args, "%"+text+"%")
	}

	// Добавляем условия, если они есть
	if len(conditions) > 0 {
		conditionString := " WHERE " + strings.Join(conditions, " AND ")
		baseQuery += conditionString
		baseCountQuery += conditionString
	}

	// Добавляем сортировку и пагинацию
	baseQuery += " ORDER BY s.id LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	// Выполняем запрос на выборку данных
	if err := r.db.SelectContext(ctx, &songs, baseQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to fetch songs: %w", err)
	}

	// Выполняем запрос на подсчет общего количества записей
	if err := r.db.GetContext(ctx, &total, baseCountQuery, args[:len(args)-2]...); err != nil {
		return nil, 0, fmt.Errorf("failed to fetch total count: %w", err)
	}

	return songs, total, nil
}
