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

func (h lyricsHandlers)UpdateTrackByID()echo.HandlerFunc{
	return func(c echo.Context) error {
	h.logger.Debugf("in handler UpdateTrackByID")

	var updateData models.UpdateTrackRequest

	ctx:=c.Request().Context()

	if err := c.Bind(&updateData); err != nil {
		h.logger.Debug("in handler UpdateTrackByID() Bind() return error: ", err)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid JSON format",
		})
	}

	err := h.lyricsUsecase.UpdateTrackByID(ctx, updateData)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Debug("in handler h.lyricsUsecase.UpdateTrackByID() return error: ", err)
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "song not found",
			})
		}
		h.logger.Debug("in handler h.lyricsUsecase.UpdateTrackByID() return error: ", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update song",
		})
	}

	h.logger.Debug("in handler UpdateTrackByID() return statusOk")
	return c.JSON(http.StatusOK, map[string]string{
		"message": "track updated successfully",
	})
	}
}

func (h lyricsHandlers)CreateTrack()echo.HandlerFunc{
	return func(c echo.Context) error {
		h.logger.Debugf("in handler CreateTrack")
		
		ctx:=c.Request().Context()

		var songRequest models.SongRequest

		if err := c.Bind(&songRequest); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid JSON for SongRequest",
		})
	}

	if err := utils.ValidateStruct(ctx, &songRequest); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Invalid JSON fields",
				"details": err.Error(), 
			})
		}

	err := h.lyricsUsecase.CreateTrack(ctx, songRequest)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to create track",
			})
		}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Track created successfully",
	})
	}
}

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

		// Возвращаем куплет клиенту
		return c.JSON(http.StatusOK, map[string]interface{}{
			"page":  page,
			"verse": verse,
		})
	}
}