package http

import (
	"github.com/22Fariz22/musiclab/internal/lyrics"
	"github.com/labstack/echo/v4"
)

// Map lyrics routes
func MapLyricsRoutes(lyricsGroup *echo.Group, h lyrics.Handlers) {
	lyricsGroup.GET("/ping", h.Ping())
	lyricsGroup.DELETE("/delete", h.DeleteSongByGroupAndTrack())
}
