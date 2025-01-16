package lyrics

import (
	"github.com/labstack/echo/v4"
)

type Handlers interface {
	Ping() echo.HandlerFunc
	DeleteSongByID() echo.HandlerFunc
	UpdateTrackByID() echo.HandlerFunc
	CreateTrack() echo.HandlerFunc
	GetSongVerseByID() echo.HandlerFunc
	GetLibrary() echo.HandlerFunc
}
