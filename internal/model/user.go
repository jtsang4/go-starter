package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"size:32;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"size:128;not null" json:"-"`
	Email     string         `gorm:"size:128;uniqueIndex;not null" json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
