package usecase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/matthewyuh246/aws-cognito/internal/domain"
	"github.com/matthewyuh246/aws-cognito/internal/repository"
)

type IAuthUsecase interface {
	LoginWithSocialProvider(ctx context.Context, provider, authCode string) (*domain.AuthTokens, error)
	exchangeCodeForTokens(authCode string) (*domain.AuthTokens, error)
	parseIDToken(idToken string) (map[string]interface{}, error)
}

type authUsecase struct {
	userRepo repository.IUserRepository
	cognitoClient *cognitoidentityprovider.CognitoIdentityProvider
	userPoolID string
	userPoolClientID string
	jwtSecret string
	httpClient *http.Client
}

func NewAuthUsecase(
	userRepo repository.IUserRepository, 
	awsSession *session.Session, 
	userPoolID, 
	userPoolClientID, 
	jwtSecret string,
) *authUsecase {
	return &authUsecase{
		userRepo: userRepo,
		cognitoClient: cognitoidentityprovider.New(awsSession),
		userPoolID: userPoolID,
		userPoolClientID: userPoolClientID,
		jwtSecret: jwtSecret,
	}
}

func (u *authUsecase) LoginWithSocialProvider(ctx context.Context, provider, authCode string) (*domain.AuthTokens, error) {
	tokens, err := u.exchangeCodeForTokens(authCode)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for tokens: %w", err)
	}

	userInfo, err := u.parseIDToken(tokens.IdToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID token: %w", err)
	}
}

func (u *authUsecase) exchangeCodeForTokens(ctx context.Context, authCode string) (*domain.AuthTokens, error) {
	cognitoDomain := os.Getenv("COGNITO_DOMAIN_URL")
	if cognitoDomain == "" {
		return nil, fmt.Errorf("COGNITO_DOMAIN_URL environment variable is not set")
	}

	if strings.Contains(cognitoDomain, "dummy-domain") {
		return &domain.AuthTokens{
			AccessToken:  "mock_access_token",
			RefreshToken: "mock_refresh_token",
			IdToken:      "mock_id_token",
			ExpiresIn:    3600,
		}, nil
	}

	redirectURI := os.Getenv("FE_URL")
	if redirectURI == "" {
		redirectURI = "http://localhost:5173"
	}
	redirectURI += "/auth/callback"

	return u.exchangeTokenWithRetry(ctx, cognitoDomain, authCode, redirectURI, 3)
}

func (u *authUsecase) exchangeTokenWithRetry(ctx context.Context, cognitoDomain, authCode, redirectURI string, maxRetries int) (*domain.AuthTokens, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(attempt*attempt) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		tokens, err := u.performTokenExchange(ctx, cognitoDomain, authCode, redirectURI)
			if err == nil {
				return tokens, nil
			}

		lastErr = err

		if !isRetriableError(err) {
			break
		}
	}
	return nil, lastErr
}

func (u *authUsecase) performTokenExchange(ctx context.Context, cognitoDomain, authCode, redirectURI string) (*domain.AuthTokens, error) {
	tokenURL := fmt.Sprintf("%s/oauth2/token", cognitoDomain)

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", u.userPoolClientID)
	data.Set("code", authCode)
	data.Set("redirect_uri", redirectURI)

	maskedCode := maskSensitiveData(authCode)
	internalLogger := fmt.Sprintf("Token exchange attempt: URL=%s, code=%s", tokenURL, maskedCode)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成に失敗: %w", err)
	}
	
	resp, err := u.httpClient.Do(req)
	if err != nil {
		fmt.Printf("INTERNAL: %s, error: %v\n", internalLogger, err)
		return nil, fmt.Errorf("ネットワークエラー: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("INTERNAL: Token exchange failed (status: %d): %s\n", resp.StatusCode, string(body))
		return nil, fmt.Errorf("認証サーバーエラー(status: %d)", resp.StatusCode)
	}

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		IdToken      string `json:"id_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("レスポンス解析エラー: %w", err)
	}

	return &domain.AuthTokens{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		IdToken:      tokenResponse.IdToken,
		ExpiresIn:    tokenResponse.ExpiresIn,
	}, nil
}

func maskSensitiveData(data string) string {
	if len(data) <= 8 {
		return "***"
	}
	return data[:4] + "***" + data[len(data)-4:]
}

func isRetriableError(err error) bool {
	return strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "connection") ||
		strings.Contains(err.Error(), "network")
}


func (u *authUsecase) parseIDToken(idToken string) (map[string]interface{}, error) {
	if idToken == "mock_id_token" {
		log.Printf("DEBUG: Using mock ID token for development")
		return map[string]interface{}{
			"email": "test@example.com",
			"username": "testuser",
			"name": "Test User",
			"picture": "https://via.placeholder.com/150",
			"sub": "mock-user-id-123",
		}, nil
	}

	if idToken != "dummy_id_token" {
		parts := strings.Split(idToken, ".")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid JWT token format")
		}

		payload := parts[1]

		if len(payload)%4 != 0 {
			payload += strings.Repeat("=", 4-len(payload)%4)
		}

		decoded, err := base64.URLEncoding.DecodeString(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to decode JWT payload: %w", err)
		}

		var claims map[string]interface{}
		if err := json.Unmarshal(decoded, &claims); err != nil {
			return nil, fmt.Errorf("failed to parse JWT claims: %w", err)
		}

		log.Printf("DEBUG: ID Token Claims: %+v", claims)

		email, _ := claims["email"].(string)
		name, _ := claims["name"].(string)
		picture, _ := claims["picture"].(string)
		sub, _ := claims["sub"].(string)

		if name == "" {
			givenName, _ := claims["given_name"].(string)
			familyName, _ := claims["family_name"].(string)
			if familyName != "" || familyName != "" {
				name = strings.TrimSpace(givenName + " " + familyName)
			}
		}

		log.Printf("DEBUG: Extracted email: %s", email)
		log.Printf("DEBUG: Extracted name: %s", name)
		log.Printf("DEBUG: Extracted picture: %s", picture)
		log.Printf("DEBUG: Extracted sub: %s", sub)

		if picture == "" {
			if avatar, ok := claims["avatar_url"].(string); ok {
				picture = avatar
				log.Printf("DEBUG: Using avatar_url as picture: %s", picture)
			} else if photo, ok := claims["photo"].(string); ok {
				picture = photo
				log.Printf("DEBUG: Using photo as picture: %s", picture)
			} else if profileImage, ok := claims["profile_image_url"].(string); ok {
				picture = profileImage
				log.Printf("DEBUG: Using profile_image_url as picture: %s", picture)
			}
		}

		username := name
		if username == "" {
			username = strings.Split(email, "@")[0]
		}

		userInfo := map[string]interface{}{
			"email": email,
			"username": username,
			"name": name,
			"picture": picture,
			"sub": sub,
		}

		log.Printf("DEBUG: Final userInfo: %+v", userInfo)

		return userInfo, nil
	}

	return nil, fmt.Errorf("invalid ID token: dummy token not allowed")
}