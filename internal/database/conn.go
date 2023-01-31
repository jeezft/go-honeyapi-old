package database

import (
	"strconv"

	"github.com/ProSellers/go-honeyapi/internal/database/models"
	"github.com/ProSellers/go-honeyapi/utils"
	"github.com/ProSellers/go-honeyapi/utils/cfg"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Latest = &Db{}

type Db struct {
	DB *gorm.DB
}

func Init() (*Db, error) {
	dsn := "host=" + cfg.Cfg.Database.Hostname + " user=" + cfg.Cfg.Database.Username + " password=" + cfg.Cfg.Database.Password + " dbname=" + cfg.Cfg.Database.Name + " port=" + strconv.Itoa(cfg.Cfg.Database.Port) + " sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	Latest.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	Latest.DB.AutoMigrate(&models.User{}, &models.Session{})

	var usr models.User
	var src models.User

	src.ID = 0

	tx := Latest.DB.Where(&src).First(&usr)

	if tx.RowsAffected < 1 {
		pp := utils.RandStringBytesMaskImprSrcSB(18)
		rt, err := models.User_Frompass("root", "root@root.pp", pp)
		if err != nil {
			log.Fatalln(err)
		}

		rt.Roles = models.ROLE_ADMIN

		Latest.createUser(rt)
		log.WithFields(log.Fields{"login": "root", "password": pp}).Info("User created")
	}

	// latest.DB.AutoMigrate() // migrage all models into DB
	return Latest, err
}

// func (d *Db) GetSession(tkn string)
