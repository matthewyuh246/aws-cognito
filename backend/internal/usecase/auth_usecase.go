package usecase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/matthewyuh246/aws-cognito/internal/domain"
	"github.com/matthewyuh246/aws-cognito/internal/repository"
)

type IAuthUsecase interface {
	LoginWithSocialProvider(ctx context.Context, provider, authCode string) (*domain.AuthTokens, error)
}

type authUsecase struct {
	userRepo      repository.IUserRepository
	authRepo      repository.IAuthRepository
	cognitoClient *cognitoidentityprovider.CognitoIdentityProvider
	userPoolID    string
	jwtSecret     string
}

func NewAuthUsecase(
	userRepo repository.IUserRepository,
	authRepo repository.IAuthRepository,
	awsSession *session.Session,
	userPoolID,
	jwtSecret string,
) *authUsecase {
	return &authUsecase{
		userRepo:      userRepo,
		authRepo:      authRepo,
		cognitoClient: cognitoidentityprovider.New(awsSession),
		userPoolID:    userPoolID,
		jwtSecret:     jwtSecret,
	}
}

func (u *authUsecase) LoginWithSocialProvider(ctx context.Context, provider, authCode string) (*domain.AuthTokens, error) {
	// 外部認証システムとの統合をリポジトリに委譲
	tokens, err := u.authRepo.ExchangeCodeForTokens(ctx, authCode)
	if err != nil {
		// ドメインエラーをユーザー向けメッセージに変換
		if authErr, ok := err.(*domain.AuthError); ok {
			return nil, fmt.Errorf(authErr.UserMessage())
		}
		return nil, fmt.Errorf("認証に失敗しました")
	}

	// IDトークンの解析（ビジネスロジック）
	userInfo, err := u.parseIDToken(tokens.IdToken)
	if err != nil {
		return nil, fmt.Errorf("ユーザー情報の取得に失敗しました")
	}

	// ユーザー処理のビジネスロジック（今後実装予定）
	log.Printf("DEBUG: User authenticated: %s", userInfo["email"])

	return tokens, nil
}

// parseIDToken - JWTのパースはビジネスロジックのため、usecaseに残す


func (u *authUsecase) parseIDToken(idToken string) (map[string]interface{}, error) {
	if idToken == "mock_id_token" {
		log.Printf("DEBUG: Using mock ID token for development")
		return map[string]interface{}{
			"email":    "test@example.com",
			"username": "testuser",
			"name":     "Test User",
			"picture":  "https://via.placeholder.com/150",
			"sub":      "mock-user-id-123",
		}, nil
	}

	if idToken == "dummy_id_token" {
		return nil, fmt.Errorf("無効なIDトークンです")
	}

	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("無効なJWTトークン形式です")
	}

	payload := parts[1]
	if len(payload)%4 != 0 {
		payload += strings.Repeat("=", 4-len(payload)%4)
	}

	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("JWTペイロードのデコードに失敗しました: %w", err)
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, fmt.Errorf("JWTクレームの解析に失敗しました: %w", err)
	}

	return u.extractUserInfo(claims), nil
}

func (u *authUsecase) extractUserInfo(claims map[string]interface{}) map[string]interface{} {
	email, _ := claims["email"].(string)
	name, _ := claims["name"].(string)
	picture, _ := claims["picture"].(string)
	sub, _ := claims["sub"].(string)

	// 名前が空の場合、given_nameとfamily_nameから構築
	if name == "" {
		givenName, _ := claims["given_name"].(string)
		familyName, _ := claims["family_name"].(string)
		if givenName != "" || familyName != "" {
			name = strings.TrimSpace(givenName + " " + familyName)
		}
	}

	// 画像URLの代替フィールドをチェック
	if picture == "" {
		if avatar, ok := claims["avatar_url"].(string); ok {
			picture = avatar
		} else if photo, ok := claims["photo"].(string); ok {
			picture = photo
		} else if profileImage, ok := claims["profile_image_url"].(string); ok {
			picture = profileImage
		}
	}

	username := name
	if username == "" {
		username = strings.Split(email, "@")[0]
	}

	return map[string]interface{}{
		"email":    email,
		"username": username,
		"name":     name,
		"picture":  picture,
		"sub":      sub,
	}
}