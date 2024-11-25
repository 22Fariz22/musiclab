package lyrics

import (
	"github.com/labstack/echo/v4"
)

type Handlers interface {
	Ping() echo.HandlerFunc
	DeleteSongByGroupAndTrack() echo.HandlerFunc
	UpdateTrackByID() echo.HandlerFunc
	CreateTrack() echo.HandlerFunc
	GetSongVerseByPage() echo.HandlerFunc
}
