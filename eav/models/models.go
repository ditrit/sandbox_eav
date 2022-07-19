package models

import (
	"strings"
	"time"
)

// Base for GORM, see https://gorm.io/docs/models.html#Declaring-Models
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func buildJson(pairs []string) string {
	var b strings.Builder
	// starting the json string
	b.WriteString("{")
	b.WriteString(strings.Join(pairs, ","))
	//ending the json string
	b.WriteString("}")
	return b.String()
}
