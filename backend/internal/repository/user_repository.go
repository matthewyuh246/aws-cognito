package repository

import (
	"context"
	"errors"

	"github.com/matthewyuh246/aws-cognito/internal/domain"
	"gorm.io/gorm"
)

type IUserRepository interface {
	createUser(ctx context.Context, user *domain.User) error
	getUserByEmail(ctx context.Context, email string) (*domain.User, error)
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

func (r *userRepository) getUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}