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

// Config представляет конфигурацию приложения.
// Содержит настройки для логирования, сервера, базы данных и аутентификации.
type Config struct {
	Logger          LoggerCfg `yaml:"Logger"`
	Server          Server    `yaml:"Server"`
	Postgres        Postgres  `yaml:"Postgres"`
	Auth            Auth      `yaml:"Auth"`
	FileStoragePath string    `yaml:"FileStoragePath"`
	UseDecode       bool      `yaml:"UseDecode"`
}

// Auth содержит конфигурацию, связанную с аутентификацией.
type Auth struct {
	SecretKey  string `yaml:"SecretKey"`
	CookieName string `yaml:"CookieName"`
}

// Server содержит конфигурацию HTTP-сервера.
type Server struct {
	ServerAddress string `yaml:"-"`
	BaseURL       string `yaml:"-"`
	Address       string `yaml:"Address"`
	RTimeout      int64  `yaml:"RTimeout"`
	WTimeout      int64  `yaml:"WTimeout"`
	Port          uint   `yaml:"Port"`
}

// Postgres содержит конфигурацию базы данных PostgreSQL.
type Postgres struct {
	Trace           bool     `yaml:"Trace"`
	MakeMigration   bool     `yaml:"MakeMigration"`
	UsePostgres     bool     `yaml:"UsePostgres"`
	SQLKeyWords     []string `yaml:"SQLKeyWords"`
	ConnString      string   `yaml:"-"`
	DriverName      string   `yaml:"DriverName"`
	Address         string   `yaml:"Address"`
	DBName          string   `yaml:"DBName"`
	User            string   `yaml:"User"`
	Password        string   `yaml:"Password"`
	MaxConnLifeTime int64    `yaml:"MaxConnLifeTime"`
	MaxConn         int      `yaml:"MaxConn"`
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

// GetConfig возвращает указатель на глобальную конфигурацию приложения.
// Инициализирует конфигурацию при первом вызове, загружая данные из файла и переменных окружения.
// Возвращает *Config.
func GetConfig() *Config {
	onceCFG.Do(func() {
		var (
			conf Config
		)
		parseConfig(&conf, cfgPath)
		if conf.UseDecode {
			decodeCFG(&conf)
		}

		flag.StringVar(&conf.Server.ServerAddress, "a", fmt.Sprintf("%s:%d", conf.Server.Address, conf.Server.Port), "HTTP server address")
		flag.StringVar(&conf.Server.BaseURL, "b", fmt.Sprintf("http://%s:%d", conf.Server.Address, conf.Server.Port), "Base URL for short links")
		flag.StringVar(&conf.FileStoragePath, "f", conf.FileStoragePath, "Path to file storage")
		// flag.StringVar(&conf.Postgres.ConnString, "d", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", conf.Postgres.User, conf.Postgres.Password, conf.Postgres.Address, conf.Postgres.DBName), "Database connect string")
		flag.StringVar(&conf.Postgres.ConnString, "d", "", "Database connect string")
		flag.Parse()

		if envAddr, exists := os.LookupEnv("SERVER_ADDRESS"); exists {
			conf.Server.ServerAddress = envAddr
		}

		if envBaseURL, exists := os.LookupEnv("BASE_URL"); exists {
			conf.Server.BaseURL = envBaseURL
		}

		if envPath, exists := os.LookupEnv("FILE_STORAGE_PATH"); exists {
			conf.FileStoragePath = envPath
		}

		if envDB, exists := os.LookupEnv("DATABASE_DSN"); exists {
			conf.Postgres.ConnString = envDB
		}

		//TODO: remove this
		if conf.Postgres.ConnString == "" {
			conf.Postgres.UsePostgres = false
		}

		config = &conf
	})

	if config == nil {
		log.Fatal("nil config")
	}

	return config
}
