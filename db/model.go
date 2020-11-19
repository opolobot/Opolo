package db

import (
	"time"

	"gorm.io/gorm"
)

// Model is the base model for our database.
type Model struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
