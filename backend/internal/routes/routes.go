package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/matthewyuh246/aws-cognito/internal/controller"
	"github.com/matthewyuh246/aws-cognito/pkg/middleware"
)

// SetupRoutes - APIルートを設定
func SetupRoutes(e *echo.Echo, authController *controller.AuthController) {
	// CORS設定
	corsConfig := middleware.NewCORSConifg()
	middleware.SetupCommonMiddleware(e, corsConfig)

	// API v1 グループ
	v1 := e.Group("/api/v1")

	// ヘルスチェック
	v1.GET("/health", authController.HealthCheck)

	// 認証関連のルート
	auth := v1.Group("/auth")
	{
		// ソーシャルログイン
		auth.POST("/login", authController.LoginWithSocialProvider)
	}
}