package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Username  string         `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Password  string         `gorm:"not null" json:"-"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
}
