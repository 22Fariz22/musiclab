package http

import (
	"github.com/22Fariz22/musiclab/internal/lyrics"
	"github.com/labstack/echo/v4"
)

// Map lyrics routes
func MapLyricsRoutes(lyricsGroup *echo.Group, h lyrics.Handlers) {
	lyricsGroup.GET("/ping", h.Ping())
	lyricsGroup.DELETE("/delete", h.DeleteSongByGroupAndTrack())
	lyricsGroup.PUT("/update", h.UpdateTrackByID())
	lyricsGroup.POST("/create", h.CreateTrack())
	lyricsGroup.GET("/verses/:id", h.GetSongVerseByID())
	lyricsGroup.GET("/library", h.GetLibrary())
}
