package domain

import "fmt"

type AuthError struct {
	Type AuthErrorType
	Code int
	Message string
	Err error
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

func NewAuthErrorWithCode(errorType AuthErrorType, code int, message string) *AuthError {
	return &AuthError{
		Type: errorType,
		Code: code,
		Message: message,
	}
}

func (e *AuthError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}