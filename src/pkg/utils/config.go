package utils

import (
	"flag"

	"github.com/spf13/viper"
)

type Config struct {
	Cfg *viper.Viper
}

func InitConfig() (*Config, error) {
	configPath := flag.String("c", "config.yaml", "config file path")
	flag.Parse()

	v := viper.New()
	v.SetConfigFile(*configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{Cfg: v}, nil
}
