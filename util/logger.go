package util

import (
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type LoggerCfg struct {
	Level  log.Level `yaml:"Level"`
	File   LogFile   `yaml:"File"`
	SysLog SysLog    `yaml:"Syslog"`
}

type LogFile struct {
	Enabled    bool   `yaml:"Enabled"`
	FileName   string `yaml:"FileName"`
	MaxSize    int    `yaml:"MaxSize"`
	MaxBackups int    `yaml:"MaxBackups"`
	MaxAge     int    `yaml:"MaxAge"`
}

type SysLog struct {
	Enabled bool   `yaml:"Enabled"`
	Address string `yaml:"Address"`
	Network string `yaml:"Network"`
	Tag     string `yaml:"Tag"`
}

var logger *log.Logger

func InitLogger(conf LoggerCfg) {
	logger = &log.Logger{
		Out:          os.Stdout,
		Hooks:        nil,
		Formatter:    &log.JSONFormatter{},
		ReportCaller: false,
		Level:        conf.Level,
		ExitFunc:     os.Exit,
	}
	if conf.File.Enabled {
		if len(conf.File.FileName) == 0 {
			panic(errors.New("parameter FileName not defined"))
		}
		rotateHook, err := NewRotateFileHook(RotateFileConfig{
			Filename:   conf.File.FileName,
			MaxSize:    conf.File.MaxSize,
			MaxBackups: conf.File.MaxBackups,
			MaxAge:     conf.File.MaxAge,
			Level:      conf.Level,
			Formatter:  &log.JSONFormatter{},
		})
		if err != nil {
			logger.Info("Failed to log to file")
			return
		}
		logger.Hooks = log.LevelHooks{}
		logger.AddHook(rotateHook)
	}
	if conf.SysLog.Enabled {
		syslog, err := NewSyslogHook(conf)
		if err != nil {
			logger.Info("Failed to log to syslog")
			return
		}
		if logger.Hooks == nil {
			logger.Hooks = log.LevelHooks{}
		}
		logger.AddHook(syslog)
	}
}

func GetLogger() *log.Logger {
	return logger
}
