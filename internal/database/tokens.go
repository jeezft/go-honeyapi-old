package database

import (
	"errors"

	"github.com/ProSellers/go-honeyapi/internal/database/models"
)

func (d *Db) GetTokens(u *models.User) (*[]models.Token, error) {
	var tokens []models.Token

	tx := d.DB.Where(&models.Token{UserID: u.ID}).Find(&tokens)

	if tx.RowsAffected < 1 || tx.Error != nil {
		return nil, errors.New("not found")
	}

	return &tokens, nil
}

func (d *Db) FindTokens(token string) (*[]models.Token, error) {
	var tokens []models.Token

	tx := d.DB.Where(&models.Token{Token: token}).Find(&tokens)

	if tx.RowsAffected < 1 || tx.Error != nil {
		return nil, errors.New("not found")
	}

	return &tokens, nil
}

func (d *Db) FindDupTokens(u *models.User, token string) bool {
	tx := d.DB.Where(&models.Token{Token: token, UserID: u.ID})

	if tx.RowsAffected < 1 || tx.Error != nil {
		return false
	}

	return true
}

func (d *Db) AddToken(u *models.User, tk string) (*models.Token, error) {
	if duplicate := d.FindDupTokens(u, tk); duplicate {
		return nil, errors.New("duplicated")
	}

	var token = models.Token{
		UserID: u.ID,
		Token:  tk,
	}

	tx := d.DB.Create(&token)

	if tx.RowsAffected < 1 || tx.Error != nil {
		return nil, errors.New("not found")
	}

	return &token, nil
}

func (d *Db) UpdateToken(t *models.Token) (*models.Token, error) {
	tx := d.DB.Where("id = ?", t.ID).Save(t)

	if tx.Error != nil || tx.RowsAffected < 1 {
		return nil, errors.New("not found")
	}

	return t, nil
}

func (d *Db) DeleteToken(t *models.Token) error {
	tx := d.DB.Delete(t)
	if tx.Error != nil || tx.RowsAffected < 1 {
		return errors.New("not found")
	}

	return nil
}
