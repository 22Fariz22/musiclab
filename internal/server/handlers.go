package server

import (
	"net/http"
	"strings"

	lyricsHTTP "github.com/22Fariz22/musiclab/internal/lyrics/delivery/http"
	lyricsRepository "github.com/22Fariz22/musiclab/internal/lyrics/repository"
	lyricsUseCase "github.com/22Fariz22/musiclab/internal/lyrics/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
		_ "github.com/22Fariz22/musiclab/docs"
)

// Map Server Handlers
func (s *Server) MapHandlers(e *echo.Echo) error {
	//check server
	e.GET("/ping", func(c echo.Context) error {
		s.logger.Debug("Pong in MapHandlers().logger Debag level")
		return c.String(http.StatusOK, "pong")
	})

	// Init repositories
	lyricsRepo := lyricsRepository.NewLyricsRepository(s.db, s.logger)

	// Init useCases
	lyricsUC := lyricsUseCase.NewLyricsUseCase(s.cfg, lyricsRepo, s.redisClient, s.logger)

	// Init handlers
	lyricsHandler := lyricsHTTP.NewLyricsHandler(s.cfg, lyricsUC, s.logger)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Static("/swagger", "./docs")

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("1000M"))

	v1 := e.Group("/api/v1")

	lyricsGroup := v1.Group("/lyrics")

	lyricsHTTP.MapLyricsRoutes(lyricsGroup, lyricsHandler)

	return nil
}
