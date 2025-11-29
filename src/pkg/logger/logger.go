package logger

import (
	"fmt"
	"main/pkg/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LoggerConfig struct {
	Level string   `yaml:"level"`
	Out   []string `yaml:"out"`
	Json  bool     `yaml:"json"`
}

func InitLogger(config *utils.Config) error {
	var logCfg LoggerConfig
	if err := config.Cfg.UnmarshalKey(utils.LogKey, &logCfg); err != nil {
		return err
	}

	fmt.Println(logCfg)

	setLogLevel(logCfg.Level)
	zerolog.TimestampFieldName = "logtime"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	setcallerMarshFunc()

	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	return nil
}

func setcallerMarshFunc() {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		f := filepath.Base(file)
		return f + ":" + strconv.Itoa(line)
	}
}

func setLogLevel(level string) {
	if logLevel, err := zerolog.ParseLevel(strings.ToLower(level)); err == nil {
		zerolog.SetGlobalLevel(logLevel)
	}
}

func Panic() *zerolog.Event {
	return log.Panic()
}

func Fatal() *zerolog.Event {
	return log.Fatal()
}

func Error() *zerolog.Event {
	return log.Error()
}

func Warn() *zerolog.Event {
	return log.Warn()
}

func Debug() *zerolog.Event {
	return log.Debug()
}

func Info() *zerolog.Event {
	return log.Info()
}

func Trace() *zerolog.Event {
	return log.Trace()
}
