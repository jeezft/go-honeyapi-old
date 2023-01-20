package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	UserID     uint
	SessionKey string `gorm:"unique"`
	TBD        time.Time
	gorm.Model
}
