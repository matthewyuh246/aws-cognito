package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/matthewyuh246/aws-cognito/pkg/utils"
)

type CORSConfig struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
}

func NewCORSConifg() *CORSConfig {
	return &CORSConfig{
		AllowOrigins: []string{utils.GetEnv("FE_URL", "http://localhost:5173")},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}
}

func SetupCORS(config *CORSConfig) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: config.AllowOrigins,
		AllowMethods: config.AllowMethods,
		AllowHeaders: config.AllowHeaders,
	})
}