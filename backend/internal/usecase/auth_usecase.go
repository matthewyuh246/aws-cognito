package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
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
	userRepo repository.IUserRepository
	cognitoClient *cognitoidentityprovider.CognitoIdentityProvider
	userPoolID string
	userPoolClientID string
	jwtSecret string
}

func NewAuthUsecase(
	userRepo repository.IUserRepository, 
	awsSession *session.Session, 
	userPoolID, 
	userPoolClientID, 
	jwtSecret string,
) authUsecase {
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
}

func (u *authUsecase) exchangeCodeforTokens(authCode string) (*domain.AuthTokens, error) {
	cognitoDomain := os.Getenv("COGNITO_DOMAIN_URL")
	if cognitoDomain == "" {
		return nil, fmt.Errorf("COGNITO_DOMAIN_URL environment variable is not set")
	}

	if strings.Contains(cognitoDomain, "dummy-domain") {
		log.Printf("DEBUG: Using mock authentication for development")
		return &domain.AuthTokens{
			AccessToken: "mock_access_token",
			RefreshToken: "mock_refresh_token",
			IdToken: "mock_id_token",
			ExpiresIn: 3600,
		}, nil
	}

	tokenURL := fmt.Sprintf("%s/oauth2/token", cognitoDomain)

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", u.userPoolClientID)
	data.Set("code", authCode)
	data.Set("redirect_uri", "http://localhost:5173/auth/callback")

	log.Printf("DEBUG: Sending token request to: %s", tokenURL)
	log.Printf("DEBUG: Request data: %+v", data)

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		log.Printf("ERROR: Failed to send token request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("ERROR: Token exchange failed (status: %d): %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("token exchange failed (status: %d): %s", resp.StatusCode, string(body))
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		IdToken string `json:"id_token"`
		ExpiresIn int `json:"expires_in`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &domain.AuthTokens{
		AccessToken: tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		IdToken: tokenResponse.IdToken,
		ExpiresIn: tokenResponse.ExpiresIn,
	}, nil
}