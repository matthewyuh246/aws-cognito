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
	Provider string `json:"provider" gorm:"not null"`
	SubjectID string `json:"subject_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}