package request

import (
	"github.com/labstack/echo/v4"
)

// LoginRequest - ソーシャルログインリクエスト
type LoginRequest struct {
	Provider string `json:"provider" validate:"required,oneof=google github facebook"`
	Code     string `json:"code" validate:"required"`
	State    string `json:"state,omitempty"`
}

// BindAndValidate - リクエストをバインドして検証
func (r *LoginRequest) BindAndValidate(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	
	// 基本的な検証
	if r.Provider == "" {
		return echo.NewHTTPError(400, "provider is required")
	}
	
	if r.Code == "" {
		return echo.NewHTTPError(400, "code is required")
	}
	
	// プロバイダーの検証
	allowedProviders := map[string]bool{
		"google":   true,
		"github":   true,
		"facebook": true,
	}
	
	if !allowedProviders[r.Provider] {
		return echo.NewHTTPError(400, "invalid provider")
	}
	
	return nil
}