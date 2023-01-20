package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Username string `gorm:"unique"`
	Password []byte
	Roles    int8
	gorm.Model
}

var salt = "Euhcmx4@ed8X2H4"

var DEFAULT_ROLES = 0

func User_Frompass(username string, password string) (*User, error) {
	p, e := GetPassHash(password)
	if e != nil {
		return nil, e
	}

	return &User{
		Username: username,
		Password: p,
		Roles:    int8(DEFAULT_ROLES),
	}, nil
}

func GetPassHash(password string) ([]byte, error) {
	p, e := bcrypt.GenerateFromPassword([]byte(salt+password), bcrypt.DefaultCost)
	return p, e
}

func ComparePass(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(salt+password))
}
