package util

import (
	"fmt"
	"io"
	"log/syslog"

	"github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type RotateFileConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Level      logrus.Level
	Formatter  logrus.Formatter
}

type RotateFileHook struct {
	Config    RotateFileConfig
	logWriter io.Writer
}

var syslogLevelMap = map[logrus.Level]syslog.Priority{
	logrus.PanicLevel: syslog.LOG_CRIT,
	logrus.FatalLevel: syslog.LOG_CRIT,
	logrus.ErrorLevel: syslog.LOG_ERR,
	logrus.WarnLevel:  syslog.LOG_WARNING,
	logrus.InfoLevel:  syslog.LOG_INFO,
	logrus.DebugLevel: syslog.LOG_DEBUG,
	logrus.TraceLevel: syslog.LOG_DEBUG,
}

func NewSyslogHook(conf LoggerCfg) (logrus.Hook, error) {
	var level syslog.Priority
	ok := false
	if level, ok = syslogLevelMap[conf.Level]; !ok {
		panic(fmt.Errorf("unknown level %s", conf.Level))
	}
	sysLogHook, err := lSyslog.NewSyslogHook(conf.SysLog.Network,
		conf.SysLog.Address, level, conf.SysLog.Tag)
	if err != nil {
		return nil, err
	}
	return sysLogHook, nil
}

func (hook *RotateFileHook) Levels() []logrus.Level {
	return logrus.AllLevels[:hook.Config.Level+1]
}

func (hook *RotateFileHook) Fire(entry *logrus.Entry) error {
	b, err := hook.Config.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.logWriter.Write(b)
	return err
}

func NewRotateFileHook(config RotateFileConfig) (logrus.Hook, error) {
	hook := RotateFileHook{
		Config: config,
	}
	hook.logWriter = &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
	}
	return &hook, nil
}
