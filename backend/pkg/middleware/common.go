package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupCommonMiddleware(e *echo.Echo, corsConfig *CORSConfig) {
	e.Use(SetupLogging())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	if corsConfig != nil {
		e.Use(SetupCORS(corsConfig))
	}
}