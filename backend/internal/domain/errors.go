package domain

import (
	"fmt"
)

// 認証に関するドメインエラー
type AuthError struct {
	Type    AuthErrorType
	Code    int
	Message string
	Err     error
}

type AuthErrorType string

const (
	AuthErrorTypeConfig      AuthErrorType = "config_error"
	AuthErrorTypeNetwork     AuthErrorType = "network_error"
	AuthErrorTypeClient      AuthErrorType = "client_error"
	AuthErrorTypeServer      AuthErrorType = "server_error"
	AuthErrorTypeSecurity    AuthErrorType = "security_error"
	AuthErrorTypeValidation  AuthErrorType = "validation_error"
	AuthErrorTypeRequest     AuthErrorType = "request_error"
	AuthErrorTypeParse       AuthErrorType = "parse_error"
)

func NewAuthError(errorType AuthErrorType, message string, err error) *AuthError {
	return &AuthError{
		Type:    errorType,
		Message: message,
		Err:     err,
	}
}

func NewAuthErrorWithCode(errorType AuthErrorType, code int, message string) *AuthError {
	return &AuthError{
		Type:    errorType,
		Code:    code,
		Message: message,
	}
}

func (e *AuthError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AuthError) Unwrap() error {
	return e.Err
}

func (e *AuthError) IsRetriable() bool {
	return (e.Code >= 500 && e.Code < 600) || e.Type == AuthErrorTypeNetwork
}

// ユーザー向けエラーメッセージ
func (e *AuthError) UserMessage() string {
	switch e.Type {
	case AuthErrorTypeConfig:
		return "システム設定エラーが発生しました"
	case AuthErrorTypeNetwork:
		return "ネットワーク接続に失敗しました"
	case AuthErrorTypeClient:
		return "リクエストが正しくありません"
	case AuthErrorTypeServer:
		return "認証サーバーエラーが発生しました"
	case AuthErrorTypeSecurity:
		return "セキュリティエラーが発生しました"
	default:
		return "認証に失敗しました"
	}
}