package util

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	onceCFG sync.Once
	config  *Config
	cfgPath = "configuration/config.yaml"
)

// ConfigTask представляет конфигурацию приложения требуемое в задание.
type ConfigTask struct {
	ServerAddres    string `json:"server_address"`    // аналог переменной окружения SERVER_ADDRESS или флага -a
	BaseURL         string `json"base_url"`           // аналог переменной окружения BASE_URL или флага -b
	FileStoragePath string `json:"file_storage_path"` // аналог переменной окружения FILE_STORAGE_PATH или флага -f
	DatabaseDNS     string `json:"database_dsn"`      // аналог переменной окружения DATABASE_DSN или флага -d
	EnableHTTPS     bool   `json:"enable_https"`      // аналог переменной окружения ENABLE_HTTPS или флага -s
}

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
	EnableHTTPS   bool   `yaml:"EnableHTTPS"`
	TLSCertPath   string `yaml:"TLSCertPath"`
	TLSKeyPath    string `yaml:"TLSKeyPath"`
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

func parseConfigTask(st interface{}, path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.WithMessage(err, "error occurred while opening JSON config file")
	}
	if err := json.Unmarshal(data, st); err != nil {
		return errors.WithMessage(err, "error occurred while unmarshaling JSON config")
	}
	return nil
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
			conf                                                 Config
			configTask                                           ConfigTask
			configPath                                           string
			serverAddress, baseURL, fileStoragePath, databaseDNS, enableHTTPS string
		)
		parseConfig(&conf, cfgPath)
		if conf.UseDecode {
			decodeCFG(&conf)
		}

		flag.StringVar(&configPath, "c", "", "Path to JSON config file")
		flag.StringVar(&configPath, "config", "", "Path to JSON config file (long)")
		flag.StringVar(&serverAddress, "a", "", "HTTP server address")
		flag.StringVar(&baseURL, "b", "", "Base URL for short links")
		flag.StringVar(&fileStoragePath, "f", "", "Path to file storage")
		// flag.StringVar(&conf.Postgres.ConnString, "d", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", conf.Postgres.User, conf.Postgres.Password, conf.Postgres.Address, conf.Postgres.DBName), "Database connect string")
		flag.StringVar(&databaseDNS, "d", "", "Database connect string")
		flag.StringVar(&enableHTTPS, "s", "", "Enable HTTPS mode")
		flag.StringVar(&conf.Server.TLSCertPath, "tls-cert", conf.Server.TLSCertPath, "Path to TLS certificate file")
		flag.StringVar(&conf.Server.TLSKeyPath, "tls-key", conf.Server.TLSKeyPath, "Path to TLS key file")
		flag.Parse()

		if envConfig, exists := os.LookupEnv("CONFIG"); exists {
			configPath = envConfig
		}

		if configPath != "" {
			if err := parseConfigTask(&configTask, configPath); err != nil {
				log.Println("error occurred while parsing config task", err)
			}
		}

		if envAddr, exists := os.LookupEnv("SERVER_ADDRESS"); exists {
			conf.Server.ServerAddress = envAddr
		} else if serverAddress != "" {
			conf.Server.ServerAddress = serverAddress
		} else if configTask.ServerAddres != "" {
			conf.Server.ServerAddress = configTask.ServerAddres
		} else {
			conf.Server.ServerAddress = fmt.Sprintf("%s:%d", conf.Server.Address, conf.Server.Port)
		}

		if envBaseURL, exists := os.LookupEnv("BASE_URL"); exists {
			conf.Server.BaseURL = envBaseURL
		} else if baseURL != "" {
			conf.Server.BaseURL = baseURL
		} else if configTask.BaseURL != "" {
			conf.Server.BaseURL = configTask.BaseURL
		} else {
			conf.Server.BaseURL = fmt.Sprintf("http://%s:%d", conf.Server.Address, conf.Server.Port)
		}

		if envPath, exists := os.LookupEnv("FILE_STORAGE_PATH"); exists {
			conf.FileStoragePath = envPath
		} else if fileStoragePath != "" {
			conf.FileStoragePath = fileStoragePath
		} else if configTask.FileStoragePath != "" {
			conf.FileStoragePath = configTask.FileStoragePath
		}

		if envDB, exists := os.LookupEnv("DATABASE_DSN"); exists {
			conf.Postgres.ConnString = envDB
		} else if databaseDNS != "" {
			conf.Postgres.ConnString = databaseDNS
		} else if configTask.DatabaseDNS != "" {
			conf.Postgres.ConnString = configTask.DatabaseDNS
		}

		if envEnableHTTPS, exists := os.LookupEnv("ENABLE_HTTPS"); exists {
			if envEnableHTTPS == "1" || strings.ToLower(envEnableHTTPS) == "true" {
				conf.Server.EnableHTTPS = true
			}
		} else if enableHTTPS != "" {
			if enableHTTPS == "1" || strings.ToLower(enableHTTPS) == "true" {
				conf.Server.EnableHTTPS = true
			}
		} else if !configTask.EnableHTTPS {
			conf.Server.EnableHTTPS = configTask.EnableHTTPS
		}

		if envCert, exists := os.LookupEnv("TLS_CERT_PATH"); exists {
			conf.Server.TLSCertPath = envCert
		}

		if envKey, exists := os.LookupEnv("TLS_KEY_PATH"); exists {
			conf.Server.TLSKeyPath = envKey
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
