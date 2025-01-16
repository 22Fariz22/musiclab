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

// DeleteSongByID удаляет песню по её ID.
// @Summary Удаление песни
// @Description Удаляет песню из базы данных по ID
// @Tags Songs
// @Param id path int true "ID песни"
// @Success 204 "Песня успешно удалена"
// @Failure 400 {object} map[string]string "Некорректный ID"
// @Failure 404 {object} map[string]string "Песня не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /songs/{id} [delete]
func (h lyricsHandlers) DeleteSongByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		h.logger.Debugf("in handler DeleteSongByID")

		ctx := c.Request().Context()

		// Получаем ID песни из параметра маршрута
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid ID format",
			})
		}

		err = h.lyricsUsecase.DeleteSongByID(ctx, uint(id))
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Debugf("song with ID=%d not found", id)
				return c.JSON(http.StatusNotFound, "Track not found")
			}
			h.logger.Debugf("error in handler DeleteSongByID()", err)
			return c.JSON(http.StatusInternalServerError, "Failed to delete song due to internal error")
		}

		h.logger.Debugf("http.StatusOK, song with ID %d is deleted")
		return c.NoContent(http.StatusOK)
	}
}

// UpdateTrackByID обновляет данные песни.
// @Summary Обновление песни
// @Description Обновляет данные песни по ID
// @Tags Songs
// @Accept json
// @Produce json
// @Param body body models.UpdateTrackRequest true "Данные для обновления"
// @Success 200 {object} map[string]string "Песня успешно обновлена"
// @Failure 400 {object} map[string]string "Некорректный запрос"
// @Failure 404 {object} map[string]string "Песня не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /songs [put]
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

// CreateTrack создает новую песню.
// @Summary Создание песни
// @Description Создает новую песню на основе данных запроса
// @Tags Songs
// @Accept json
// @Produce json
// @Param body body models.SongRequest true "Данные новой песни"
// @Success 200 {object} models.SongDetail "Созданная песня"
// @Failure 400 "Некорректные данные"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /songs [post]
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

// GetSongVerseByID получает куплет песни.
// @Summary Получение куплета
// @Description Возвращает куплет песни по ID песни и номеру страницы
// @Tags Songs
// @Param id path int true "ID песни"
// @Param page query int true "Номер страницы"
// @Success 200 {object} map[string]interface{} "Куплет песни"
// @Failure 400 {object} map[string]string "Некорректный ID или номер страницы"
// @Failure 404 {object} map[string]string "Куплет не найден"
// @Router /songs/{id}/verses [get]
func (h lyricsHandlers) GetSongVerseByID() echo.HandlerFunc {
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
		verse, err := h.lyricsUsecase.GetSongVerseByID(ctx, uint(id), page)
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

// GetLibrary возвращает библиотеку песен.
// @Summary Получение библиотеки
// @Description Возвращает список песен на основе фильтров
// @Tags Songs
// @Param group query string false "Фильтр по группе"
// @Param song query string false "Фильтр по названию песни"
// @Param text query string false "Фильтр по тексту"
// @Param release_date query string false "Фильтр по дате выпуска"
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество записей на странице"
// @Success 200 {object} map[string]interface{} "Список песен"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /library [get]
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
