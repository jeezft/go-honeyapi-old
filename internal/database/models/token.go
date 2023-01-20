package models

import "gorm.io/gorm"

type Token struct {
	UserID     uint
	Token      string
	gorm.Model `json:"-"`
}
