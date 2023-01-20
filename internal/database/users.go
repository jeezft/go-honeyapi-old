package database

import (
	"errors"

	"github.com/ProSellers/go-honeyapi/internal/database/models"
)

func (d *Db) createUser(u *models.User) (*models.User, error) {
	tx := d.DB.Create(u)

	if tx.RowsAffected < 1 {
		return nil, errors.New("user already exists")
	}

	return u, nil
}

func (d *Db) CreateUser(username string, password string) (*models.User, error) {
	u, e := models.User_Frompass(username, password)
	if e != nil {
		return nil, e
	}

	return d.createUser(u)
}

func (d *Db) UpdateUser(u *models.User) (*models.User, error) {
	tx := d.DB.Where(&models.User{Username: u.Username}).Save(u)
	if tx.Error != nil || tx.RowsAffected < 1 {
		return nil, errors.New("not found")
	}

	return u, nil
}

func (d *Db) FindUser(username string) (*models.User, error) {
	var user models.User

	tx := d.DB.Where(&models.User{Username: username}).First(&user)
	if tx.Error != nil || tx.RowsAffected < 1 {
		return nil, errors.New("not found")
	}

	return &user, nil
}

func (d *Db) CheckPassword(user *models.User, password string) error {
	return models.ComparePass(user.Password, password)
}
