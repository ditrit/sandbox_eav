package models

import (
	"time"
)

// Base for GORM, see https://gorm.io/docs/models.html#Declaring-Models
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
