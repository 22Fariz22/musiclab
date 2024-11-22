package http

import (
	"net/http"

	"github.com/22Fariz22/musiclab/config"
	"github.com/22Fariz22/musiclab/internal/lyrics"
	"github.com/22Fariz22/musiclab/pkg/logger"
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
