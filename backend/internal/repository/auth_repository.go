package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/matthewyuh246/aws-cognito/internal/domain"
	"github.com/matthewyuh246/aws-cognito/pkg/httpclient"
	"github.com/matthewyuh246/aws-cognito/pkg/logger"
	"github.com/matthewyuh246/aws-cognito/pkg/utils"
)

type IAuthRepository interface {
	ExchangeCodeForTokens(ctx context.Context, authCode string) (*domain.AuthTokens, error)
}

type authRepository struct {
	httpClient       *httpclient.Client
	logger           *logger.Logger
	cognitoDomain    string
	userPoolClientID string
	allowedDomains   []string
}

type AuthConfig struct {
	CognitoDomain    string
	UserPoolClientID string
	AllowedDomains   []string
}

func NewAuthRepository(config AuthConfig) IAuthRepository {
	httpConfig := httpclient.Config{
		Timeout:     30 * time.Second,
		MaxRetries:  3,
		BaseBackoff: 1 * time.Second,
		MaxBackoff:  30 * time.Second,
		JitterMax:   1 * time.Second,
	}

	authLogger := logger.New("AUTH")
	client := httpclient.NewClient(httpConfig, authLogger)

	return &authRepository{
		httpClient:       client,
		logger:           authLogger,
		cognitoDomain:    config.CognitoDomain,
		userPoolClientID: config.UserPoolClientID,
		allowedDomains:   config.AllowedDomains,
	}
}

func (r *authRepository) ExchangeCodeForTokens(ctx context.Context, authCode string) (*domain.AuthTokens, error) {
	// 開発環境のモック処理
	if strings.Contains(r.cognitoDomain, "dummy-domain") {
		return &domain.AuthTokens{
			AccessToken:  "mock_access_token",
			RefreshToken: "mock_refresh_token",
			IdToken:      "mock_id_token",
			ExpiresIn:    3600,
		}, nil
	}

	redirectURI, err := r.buildAndValidateRedirectURI()
	if err != nil {
		return nil, err
	}

	return r.performTokenExchange(ctx, authCode, redirectURI)
}

func (r *authRepository) buildAndValidateRedirectURI() (string, error) {
	feURL := utils.GetEnv("FE_URL", "")
	if feURL == "" {
		return "", domain.NewAuthError(domain.AuthErrorTypeConfig, "FE_URL環境変数が設定されていません", nil)
	}

	// ホワイトリスト検証
	if !r.isAllowedDomain(feURL) {
		return "", domain.NewAuthError(domain.AuthErrorTypeSecurity, "許可されていないフロントエンドURLです", nil)
	}

	// パス正規化
	parsedURL, err := url.Parse(feURL)
	if err != nil {
		return "", domain.NewAuthError(domain.AuthErrorTypeConfig, "FE_URLの形式が正しくありません", err)
	}

	parsedURL.Path = path.Join(parsedURL.Path, "/auth/callback")
	return parsedURL.String(), nil
}

func (r *authRepository) isAllowedDomain(domain string) bool {
	for _, allowed := range r.allowedDomains {
		if domain == allowed {
			return true
		}
	}
	return false
}

func (r *authRepository) performTokenExchange(ctx context.Context, authCode, redirectURI string) (*domain.AuthTokens, error) {
	tokenURL := fmt.Sprintf("%s/oauth2/token", r.cognitoDomain)

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", r.userPoolClientID)
	data.Set("code", authCode)
	data.Set("redirect_uri", url.QueryEscape(redirectURI))

	maskedCode := utils.MaskSensitiveData(authCode, 4, 4)
	r.logger.Debug("トークン交換リクエスト開始", map[string]interface{}{
		"url":         tokenURL,
		"code_masked": maskedCode,
		"client_id":   r.userPoolClientID,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, domain.NewAuthError(domain.AuthErrorTypeRequest, "リクエスト作成に失敗しました", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := r.httpClient.DoWithRetry(ctx, req)
	if err != nil {
		return nil, domain.NewAuthError(domain.AuthErrorTypeNetwork, "ネットワーク接続に失敗しました", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, domain.NewAuthErrorWithCode(
			r.categorizeHTTPError(resp.StatusCode),
			resp.StatusCode,
			"認証サーバーエラー",
		)
	}

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		IdToken      string `json:"id_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, domain.NewAuthError(domain.AuthErrorTypeParse, "レスポンス解析に失敗しました", err)
	}

	if err := r.validateTokenResponse(&tokenResponse); err != nil {
		return nil, err
	}

	r.logger.Info("トークン交換成功", map[string]interface{}{
		"client_id": r.userPoolClientID,
	})

	return &domain.AuthTokens{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		IdToken:      tokenResponse.IdToken,
		ExpiresIn:    tokenResponse.ExpiresIn,
	}, nil
}

func (r *authRepository) categorizeHTTPError(statusCode int) domain.AuthErrorType {
	switch {
	case statusCode >= 400 && statusCode < 500:
		return domain.AuthErrorTypeClient
	case statusCode >= 500 && statusCode < 600:
		return domain.AuthErrorTypeServer
	default:
		return domain.AuthErrorTypeRequest
	}
}

func (r *authRepository) validateTokenResponse(response *struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	ExpiresIn    int    `json:"expires_in"`
}) error {
	if response.AccessToken == "" {
		return domain.NewAuthError(domain.AuthErrorTypeValidation, "アクセストークンが空です", nil)
	}
	if response.IdToken == "" {
		return domain.NewAuthError(domain.AuthErrorTypeValidation, "IDトークンが空です", nil)
	}
	if response.ExpiresIn <= 0 {
		return domain.NewAuthError(domain.AuthErrorTypeValidation, "無効な有効期限です", nil)
	}
	return nil
}