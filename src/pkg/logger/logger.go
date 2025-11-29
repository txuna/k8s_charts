package logger

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type LoggerConfig struct {
	Level string   `yaml:"level"`
	Out   []string `yaml:"out"`
	Json  bool     `yaml:"json"`
}

func loadConfig() (*LoggerConfig, error) {
	configPath := flag.String("c", "config.yaml", "config file path")
	flag.Parse()

	v := viper.New()
	v.SetConfigFile(*configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	var logCfg LoggerConfig
	if err := v.UnmarshalKey("log", &logCfg); err != nil {
		return nil, fmt.Errorf("unmarshal log config: %w", err)
	}

	return &logCfg, nil
}

func InitLogger() error {
	logCfg, err := loadConfig()
	if err != nil {
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
