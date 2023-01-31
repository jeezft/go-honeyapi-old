package cfg

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Database *db
	Captcha  *captcha
	Jwt      *jwtauth
	Port     int
}

type db struct {
	Hostname string
	Port     int
	Name     string
	Username string
	Password string
}

type captcha struct {
	Sitekey   string
	SecretKey string
}

type jwtauth struct {
	SecretKey string
}

var Cfg *Config
var confPath = "./config.toml"

func Load() {

	f, err := os.Open(confPath)
	switch {
	case errors.Is(err, os.ErrNotExist):
		f.Close()
		log.Warn("No file found, creating config...")
		if err := Write(SetDefaults()); err != nil {
			log.WithFields(log.Fields{"error": err}).Fatalln("Can't Write to file.")
		}
		log.Info("Config generated. Please, edit before continue.\n\nPress any key to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		Load()
	case err != nil:
		{
			log.WithFields(log.Fields{"error": err}).Fatalln("Unknown error")
		}
	}

	r, err := io.ReadAll(f)
	f.Close()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatalln("Error while reading file.")
	}

	var conf Config

	_, err = toml.Decode(string(r), &conf)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatalln("Can't decode from filereader.")
	}
	Cfg = &conf

	log.Info("Config loaded.")
}

func Write(c *Config) error {
	f, e := os.Create(confPath)
	if e != nil {
		return e
	}

	if e = toml.NewEncoder(f).Encode(c); e != nil {
		return e
	}

	if e = f.Close(); e != nil {
		return e
	}
	// f.Write()

	return nil
}

func SetDefaults() *Config {
	return &Config{
		Database: &db{
			Hostname: "localhost",
			Port:     5432,
			Name:     "prosellers",
			Username: "postgres",
			Password: " ",
		},
		Captcha: &captcha{
			Sitekey:   "",
			SecretKey: "",
		},
		Jwt: &jwtauth{
			SecretKey: "insecure",
		},
	}
}
