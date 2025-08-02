package usecase

import (
	"context"
	"fmt"
	"os"

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
		return nil, fmt.Errorf("COGNITO_DOMAIN_URL")
	}
}