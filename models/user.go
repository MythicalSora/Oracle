package models

import (
	"github.com/jinzhu/gorm"
)

// User ...
type User struct {
	gorm.Model
	DiscordID string `gorm:"NOT_NULL"`
	Balance   float64
	LastDaily string
	Houses    []House `gorm:"foreignkey:Owner"`
}
