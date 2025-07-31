package usecase

import (
	"context"

	"github.com/matthewyuh246/aws-cognito/internal/domain"
)

type IAuthUsecase interface {
	LoginWithSocialProvider(ctx context.Context, provider, authCode string) (*domain.AuthTokens, error)
}

