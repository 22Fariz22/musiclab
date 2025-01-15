package http

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/22Fariz22/musiclab/config"
	"github.com/22Fariz22/musiclab/internal/lyrics"
	"github.com/22Fariz22/musiclab/internal/models"
	"github.com/22Fariz22/musiclab/pkg/logger"
	"github.com/22Fariz22/musiclab/pkg/utils"
	"github.com/labstack/echo/v4"
)

type lyricsHandlers struct {
	cfg           *config.Config
	lyricsUsecase lyrics.UseCase
	logger        logger.Logger
}

func NewLyricsHandler(cfg *config.Config, lyricsUsecase lyrics.UseCase, logger logger.Logger) lyrics.Handlers {
	return &lyricsHandlers{cfg: cfg, lyricsUsecase: lyricsUsecase, logger: logger}
}

// Ping godoc
// @Summary Проверка доступности базы данных
// @Description Проверяет доступность базы данных, возвращает "pong"
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {string} string "pong"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /lyrics/ping [get]
func (h lyricsHandlers) Ping() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.logger.Debug("Call Handler Ping()")

		err := h.lyricsUsecase.Ping()
		if err != nil {
			h.logger.Debug("error in handlers Ping()")
			return c.JSON(echo.ErrInternalServerError.Code, "")
		}
		return c.JSON(http.StatusOK, "pong")
	}
}

// DeleteSongByGroupAndTrack godoc
// @Summary Удалить песню
// @Description Удаляет песню из библиотеки по названию группы и трека
// @Tags Lyrics
// @Accept json
// @Produce json
// @Param group query string true "Название группы"
// @Param track query string true "Название трека"
// @Success 200 {string} string "Track is deleted"
// @Failure 400 {object} map[string]string "Group and song name are required"
// @Failure 404 {object} map[string]string "Track not found"
// @Failure 500 {object} map[string]string "Failed to delete song"
// @Router /lyrics/delete [delete]
func (h lyricsHandlers) DeleteSongByGroupAndTrack() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.logger.Debugf("in handler DeleteSongByGroupAndTrack")

		ctx := c.Request().Context()

		groupName := c.QueryParam("group")
		trackName := c.QueryParam("track")
		h.logger.Debugf("in handler queryparams group=%s track=%s", groupName, trackName)

		if groupName == "" || trackName == "" {
			h.logger.Debug("Group and song name are required")
			return c.JSON(http.StatusBadRequest, "Group and song name are required")
		}

		err := h.lyricsUsecase.DeleteSongByGroupAndTrack(ctx, groupName, trackName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Debug("Track not found")
				return c.JSON(http.StatusNotFound, "Track not found")
			}
			h.logger.Debugf("error in handler DeleteSongByGroupAndTrack()", err)
			return c.JSON(http.StatusInternalServerError, "Failed to delete song")
		}

		h.logger.Debugf("http.StatusOK, Track %s is deleted", trackName)
		return c.JSON(http.StatusOK, "Track is deleted")
	}
}

// UpdateTrackByID godoc
// @Summary Обновить данные песни
// @Description Обновляет данные песни по идентификатору
// @Tags Lyrics
// @Accept json
// @Produce json
// @Param song body models.UpdateTrackRequest true "Данные для обновления песни"
// @Success 200 {object} map[string]string "Track updated successfully"
// @Failure 400 {object} map[string]string "Invalid JSON format"
// @Failure 404 {object} map[string]string "Song not found"
// @Failure 500 {object} map[string]string "Failed to update song"
// @Router /lyrics/update [put]

func (h lyricsHandlers) UpdateTrackByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.logger.Debugf("in handler UpdateTrackByID")

		var updateData models.UpdateTrackRequest

		// Привязка данных из запроса
		if err := c.Bind(&updateData); err != nil {
			h.logger.Debug("in handler UpdateTrackByID() Bind() return error: ", err)
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid JSON format",
			})
		}

		// Валидация данных
		if err := c.Validate(&updateData); err != nil {
			h.logger.Debug("in handler UpdateTrackByID() Validate() return error: ", err)
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "validation failed",
				"details": err.Error(),
			})
		}

		// Логика обновления данных
		err := h.lyricsUsecase.UpdateTrackByID(c.Request().Context(), updateData)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Debug("song not found")
				return c.JSON(http.StatusNotFound, map[string]string{
					"error": "song not found",
				})
			}
			h.logger.Debug("failed to update song: ", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "failed to update song",
			})
		}

		h.logger.Debug("track updated successfully")
		return c.JSON(http.StatusOK, map[string]string{
			"message": "track updated successfully",
		})
	}
}

func (h lyricsHandlers) CreateTrack() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.logger.Debugf("in handler CreateTrack")

		ctx := c.Request().Context()

		var songRequest models.SongRequest

		if err := c.Bind(&songRequest); err != nil {
			h.logger.Debugf("Failed to Bind() in CreateTrack(): %v", err)
			return c.NoContent(http.StatusBadRequest)
		}

		if err := utils.ValidateStruct(ctx, &songRequest); err != nil {
			h.logger.Debugf("Failed to ValidateStruct() in CreateTrack(): %v", err)
			return c.NoContent(http.StatusBadRequest)
		}

		songDetail, err := h.lyricsUsecase.CreateTrack(ctx, songRequest)
		if err != nil {
			h.logger.Debugf("Failed to create track: %v", err)
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, songDetail)
	}
}

// GetSongVerseByPage godoc
// @Summary Получить конкретный куплет песни
// @Description Возвращает конкретный куплет песни по идентификатору песни и номеру страницы
// @Tags Lyrics
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param page query int true "Номер страницы (куплета)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /lyrics/verses/{id} [get]
func (h lyricsHandlers) GetSongVerseByPage() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.logger.Debug("In handler GetSongVerseByPage")

		ctx := c.Request().Context()

		// Получаем ID песни из параметра маршрута
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid song ID",
			})
		}

		// Получаем номер страницы из query-параметров
		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil || page <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid page number",
			})
		}

		// Вызываем usecase для получения куплета
		verse, err := h.lyricsUsecase.GetSongVerseByPage(ctx, uint(id), page)
		if err != nil {
			h.logger.Errorf("Error fetching verse: %v", err)
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}

		// Проверяем на пустой текст для первой страницы
		if page == 1 && verse == "" {
			h.logger.Debugf("Page %d has an empty verse", page)
			return c.JSON(http.StatusOK, map[string]interface{}{
				"page":  page,
				"verse": "",
			})
		}

		// Возвращаем куплет клиенту
		return c.JSON(http.StatusOK, map[string]interface{}{
			"page":  page,
			"verse": verse,
		})
	}
}

// GetLibrary godoc
// @Summary      Получить библиотеку
// @Description  Получить список песен с фильтрацией по полям и пагинацией
// @Tags         Lyrics
// @Param        group        query    string  false  "Название группы"
// @Param        song         query    string  false  "Название песни"
// @Param        releaseDate  query    string  false  "Дата релиза"
// @Param        page         query    int     false  "Номер страницы"
// @Param        limit        query    int     false  "Количество элементов на странице"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /lyrics/library [get]
func (h lyricsHandlers) GetLibrary() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		group := c.QueryParam("group")
		song := c.QueryParam("song")
		text := c.QueryParam("text")
		releaseDate := c.QueryParam("release_date")

		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil || page <= 0 {
			page = 1
		}

		limit, err := strconv.Atoi(c.QueryParam("limit"))
		if err != nil || limit <= 0 {
			limit = 10
		}

		songs, total, err := h.lyricsUsecase.GetLibrary(ctx, group, song, text, releaseDate, page, limit)
		if err != nil {
			h.logger.Errorf("Error in GetLibrary: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to fetch library",
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
			"data":  songs,
		})
	}
}
