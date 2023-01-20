package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/ProSellers/go-honeyapi/internal/database/models"
	"github.com/ProSellers/go-honeyapi/utils/cfg"
	"github.com/golang-jwt/jwt"
)

// func (d *Db) CreateSession(user *models.User) string {
// 	authKey := jwt.New(jwt.SigningMethodEdDSA)
// 	claims := authKey.Claims.(jwt.MapClaims)

// 	claims["authorized"] = true
// 	claims["uID"] = user.ID
// 	key, err := authKey.SignedString(cfg.Cfg.Jwt.SecretKey)
// 	if err != nil {
// 		return ""
// 	}
// 	return key
// }

func (d *Db) CreateSession(user *models.User, ttl time.Duration) (string, error) {
	authKey := jwt.New(jwt.SigningMethodHS256)
	claims := authKey.Claims.(jwt.MapClaims)

	tbd := time.Now().Add(ttl)

	claims["authorized"] = true
	claims["uID"] = user.ID
	claims["tbd"] = tbd.Format(time.RFC1123)

	key, err := authKey.SignedString([]byte(cfg.Cfg.Jwt.SecretKey))
	if err != nil {
		fmt.Println(12123)
		return "", err
	}
	session := models.Session{
		UserID:     user.ID,
		SessionKey: key,
		TBD:        tbd,
	}
	tx := d.DB.Create(&session)
	if tx.Error != nil {
		fmt.Println(123123)
		return "", err
	}
	return key, nil
}

// func (d *Db) CheckSession(session string) (*models.Session, error) {
// 	claims := jwt.MapClaims{}
// 	_, err := jwt.ParseWithClaims(session, claims, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(cfg.Cfg.Jwt.SecretKey), nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	fmt.Println(claims["uID"])
// 	fmt.Println(claims["tbd"].(time.Time))
// 	return nil, nil
// }

func (d *Db) FindSession(session string) (*models.Session, error) {
	s := &models.Session{}
	tx := d.DB.Where(&models.Session{SessionKey: session}).First(s)
	if tx.Error != nil || tx.RowsAffected < 1 {
		return nil, errors.New("not found")
	}
	return s, nil
}

func (d *Db) CheckSession(session string) (*models.Session, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(session, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Cfg.Jwt.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	_, ok := claims["uID"]
	if !ok {
		return nil, fmt.Errorf("claims['uID'] does not exist")
	}

	tbd, ok := claims["tbd"]
	if !ok {
		return nil, fmt.Errorf("claims['tbd'] does not exist")
	}

	t, err := time.Parse(time.RFC1123, tbd.(string))
	if err != nil {
		return nil, err
	}

	if time.Now().After(t) {
		return nil, fmt.Errorf("token is expired")
	}

	s, e := d.FindSession(session)

	return s, e
}

func (d *Db) FindUserByID(ID uint) (*models.User, error) {
	u := &models.User{}
	tx := d.DB.Where("id = ?", ID).First(u)
	if tx.Error != nil || tx.RowsAffected < 1 {
		return nil, errors.New("not found")
	}
	return u, nil
}
