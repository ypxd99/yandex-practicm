package util

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	onceCFG sync.Once
	config  *Config
	cfgPath = "configuration/config.yaml"
)

type Config struct {
	Logger   LoggerCfg `yaml:"Logger"`
	Server   Server    `yaml:"Server"`
	Postgres Postgres  `yaml:"Postgres"`
}

type Server struct {
	BaseURL       string
	ServerAddress string
	Address       string `yaml:"Address"`
	Port          uint   `yaml:"Port"`
	RTimeout      int64  `yaml:"RTimeout"`
	WTimeout      int64  `yaml:"WTimeout"`
}

type Postgres struct {
	DriverName      string   `yaml:"DriverName"`
	Address         string   `yaml:"Address"`
	DBName          string   `yaml:"DBName"`
	User            string   `yaml:"User"`
	Password        string   `yaml:"Password"`
	MaxConn         int      `yaml:"MaxConn"`
	MaxConnLifeTime int64    `yaml:"MaxConnLifeTime"`
	Trace           bool     `yaml:"Trace"`
	MakeMigration   bool     `yaml:"MakeMigration"`
	UsePostgres     bool     `yaml:"UsePostgres"`
	SQLKeyWords     []string `yaml:"SQLKeyWords"`
}

func decode(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", errors.WithMessage(err, "error occurred while decoding string(base64)")
	}
	res, err := GetRSA().Decrypt(data)
	if err != nil {
		return "", errors.WithMessage(err, "error occurred while decoding string(RSA)")
	}
	return string(res), err
}

func parseConfig(st interface{}, cfgPath string) {
	f, err := os.Open(cfgPath)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "error occurred while opening cfg file"))
	}

	fi, err := f.Stat()
	if err != nil {
		log.Fatal(errors.WithMessage(err, "error occurred while getting file stats"))
	}

	data := make([]byte, fi.Size())
	_, err = f.Read(data)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "error occurred while reading data"))
	}

	err = yaml.Unmarshal(data, st)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "error occurred while unmashaling data"))
	}
}

func decodeCFG(cfg *Config) error {
	var err error
	cfg.Postgres.Address, err = decode(cfg.Postgres.Address)
	if err != nil {
		return errors.WithMessage(err, "error occurred while decode address")
	}
	cfg.Postgres.User, err = decode(cfg.Postgres.User)
	if err != nil {
		return errors.WithMessage(err, "error occurred while decode user")
	}
	cfg.Postgres.Password, err = decode(cfg.Postgres.Password)
	if err != nil {
		return errors.WithMessage(err, "error occurred while decode password")
	}

	return nil
}

func GetConfig() *Config {
	onceCFG.Do(func() {
		var (
			conf Config
		)
		parseConfig(&conf, cfgPath)
		if conf.Postgres.UsePostgres {
			decodeCFG(&conf)
		}

		flag.StringVar(&conf.Server.ServerAddress, "a", fmt.Sprintf("%s:%d", conf.Server.Address, conf.Server.Port), "HTTP server address")
		flag.StringVar(&conf.Server.BaseURL, "b", fmt.Sprintf("http://%s:%d", conf.Server.Address, conf.Server.Port), "Base URL for short links")
		flag.Parse()

		config = &conf
	})

	if config == nil {
		log.Fatal("nil config")
	}

	return config
}
