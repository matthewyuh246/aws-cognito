package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/matthewyuh246/aws-cognito/internal/controller/request"
	"github.com/matthewyuh246/aws-cognito/internal/controller/response"
	"github.com/matthewyuh246/aws-cognito/internal/usecase"
	"github.com/matthewyuh246/aws-cognito/pkg/logger"
)

type AuthController struct {
	authUsecase usecase.IAuthUsecase
	logger      *logger.Logger
}

func NewAuthController(authUsecase usecase.IAuthUsecase) *AuthController {
	return &AuthController{
		authUsecase: authUsecase,
		logger:      logger.New("AUTH_CONTROLLER"),
	}
}

// LoginWithSocialProvider - ソーシャルプロバイダーでのログイン
func (ac *AuthController) LoginWithSocialProvider(c echo.Context) error {
	var req request.LoginRequest
	
	// リクエストのバインドと検証
	if err := req.BindAndValidate(c); err != nil {
		ac.logger.Error("リクエストバインドエラー", map[string]interface{}{
			"error": err.Error(),
		})
		return response.SendBadRequest(c, "無効なリクエストです")
	}

	ac.logger.Info("ソーシャルログイン開始", map[string]interface{}{
		"provider":    req.Provider,
		"code_masked": maskCode(req.Code),
	})

	// ビジネスロジックの実行
	tokens, err := ac.authUsecase.LoginWithSocialProvider(c.Request().Context(), req.Provider, req.Code)
	if err != nil {
		ac.logger.Error("認証エラー", map[string]interface{}{
			"provider": req.Provider,
			"error":    err.Error(),
		})
		return response.SendUnauthorized(c, err.Error())
	}

	// 成功レスポンス（実際の実装では、usecaseからユーザー情報も取得）
	userInfo := map[string]interface{}{
		"email":    "user@example.com", // 実際の実装では動的に取得
		"username": "user",
		"name":     "User Name",
		"picture":  "",
		"sub":      "user-sub-id",
	}

	ac.logger.Info("ソーシャルログイン成功", map[string]interface{}{
		"provider": req.Provider,
		"user_id":  userInfo["sub"],
	})

	return response.SendLoginSuccess(c, tokens, userInfo)
}

// HealthCheck - ヘルスチェック
func (ac *AuthController) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"message": "Auth service is healthy",
	})
}

// maskCode - 認可コードのマスク
func maskCode(code string) string {
	if len(code) <= 8 {
		return "***"
	}
	return code[:4] + "***" + code[len(code)-4:]
}
