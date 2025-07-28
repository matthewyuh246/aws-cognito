package repository

import (
	"context"

	"github.com/matthewyuh246/aws-cognito/internal/domain"
	"gorm.io/gorm"
)

type IUserRepository interface {
	createUser(ctx context.Context, user *domain.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db:db}
}

func (r *userRepository) createUser(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}