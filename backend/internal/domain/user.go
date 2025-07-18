package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID uint `json:"id" gorm:"primarykey"`
	Email string `json:"email" gorm:"uniqueIndex;not null"`
	Username string `json:"username" gorm:"uniqueIndex;not null"`
	Name string `json:"name"`
	Picture string `json:"picture"`
	Provider string `json:"provider" gorm:"not null"`
	SubjectID string `json:"subject_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type AuthTokens struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken string `json:"id_token"`
	ExpiresIn int `json:"expires_in"`
}

type UserClaims struct {
	UserID uint `json:"user_id"`
	Email string `json:"email"`
	Username string `json:"username"`
	Name string `json:"name"`
	Picture string `json:"picture"`
	Provider string `json:"provider"`
	Exp int64 `json:"exp"`
}