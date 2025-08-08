package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/matthewyuh246/aws-cognito/internal/domain"
)

// LoginResponse - ログイン成功レスポンス
type LoginResponse struct {
	Success     bool                `json:"success"`
	Message     string              `json:"message"`
	Tokens      *domain.AuthTokens  `json:"tokens,omitempty"`
	User        *UserInfo           `json:"user,omitempty"`
}

// UserInfo - ユーザー情報
type UserInfo struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Sub      string `json:"sub"`
}

// ErrorResponse - エラーレスポンス
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// SendLoginSuccess - ログイン成功レスポンスを送信
func SendLoginSuccess(c echo.Context, tokens *domain.AuthTokens, userInfo map[string]interface{}) error {
	user := &UserInfo{
		Email:    getStringFromMap(userInfo, "email"),
		Username: getStringFromMap(userInfo, "username"),
		Name:     getStringFromMap(userInfo, "name"),
		Picture:  getStringFromMap(userInfo, "picture"),
		Sub:      getStringFromMap(userInfo, "sub"),
	}

	response := LoginResponse{
		Success: true,
		Message: "ログインが成功しました",
		Tokens:  tokens,
		User:    user,
	}

	return c.JSON(http.StatusOK, response)
}

// SendError - エラーレスポンスを送信
func SendError(c echo.Context, statusCode int, message string, code ...string) error {
	response := ErrorResponse{
		Success: false,
		Message: message,
	}

	if len(code) > 0 {
		response.Code = code[0]
	}

	return c.JSON(statusCode, response)
}

// SendBadRequest - 400エラーレスポンス
func SendBadRequest(c echo.Context, message string) error {
	return SendError(c, http.StatusBadRequest, message, "BAD_REQUEST")
}

// SendUnauthorized - 401エラーレスポンス
func SendUnauthorized(c echo.Context, message string) error {
	return SendError(c, http.StatusUnauthorized, message, "UNAUTHORIZED")
}

// SendInternalServerError - 500エラーレスポンス
func SendInternalServerError(c echo.Context, message string) error {
	return SendError(c, http.StatusInternalServerError, message, "INTERNAL_SERVER_ERROR")
}

// ヘルパー関数
func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}