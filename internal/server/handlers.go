package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Map Server Handlers
func (s *Server) MapHandlers(e *echo.Echo) error {
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	return nil
}