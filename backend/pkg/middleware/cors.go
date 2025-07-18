package middleware

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

func SetupOCRS(config *CORSConfig) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: config.AllowOrigins,
		AllowMethods: config.AllowMethods,
		AllowHeaders: config.AllowHeaders,
	})
}