package main

import (
	"main/pkg/logger"
	mq "main/pkg/mq/nats"
	"main/pkg/utils"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	config, err := utils.InitConfig()
	if err != nil {
		panic(err)
	}

	if err := logger.InitLogger(config); err != nil {
		panic(err)
	}

	cfg, err := loadConfig(config)
	if err != nil {
		panic(err)
	}

	logger.Info().Msgf("NATS URL: %s", cfg.Nats.Url)
	messageQueue, err := mq.NewNats(cfg.Nats.Url)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create NATS client")
	}

	defer messageQueue.Close()

	if err := messageQueue.CreateStream("stream", "subject.>", 3, jetstream.LimitsPolicy); err != nil {
		logger.Fatal().Err(err).Msg("Failed to create stream")
	}

	for range 5000 {
		if err := messageQueue.Publish("stream", "subject.test", "Hello, NATS!"); err != nil {
			logger.Error().Err(err).Msg("Failed to publish message")
		}

		// logger.Info().Msgf("[%d] Published message to subject.test", i)
	}

	for range 5000 {
		if err := messageQueue.Publish("stream", "subject.good", "Hello, NATS!"); err != nil {
			logger.Error().Err(err).Msg("Failed to publish message")
		}

		// logger.Info().Msgf("[%d] Published message to subject.good", i)
	}

	logger.Info().Msg("Published 10000 messages")

	go func() {
		if err := messageQueue.Consume("stream", "subject.good", "subject.good"); err != nil {
			logger.Error().Err(err).Msg("Failed to consume messages")
		}
	}()

	time.Sleep(3 * time.Second)

	if err := messageQueue.PurgeStream("stream", "subject.good"); err != nil {
		logger.Error().Err(err).Msg("Failed to purge stream")
	}

	logger.Info().Msg("Purge completed")

	if err := messageQueue.PurgeStream("stream", "subject.test"); err != nil {
		logger.Error().Err(err).Msg("Failed to purge stream")
	}

	logger.Info().Msg("Purge completed")

	time.Sleep(20 * time.Second)

	// dd := []uint64{}
	// for i := range 10000 {
	// 	dd = append(dd, uint64(i+1))
	// }

	// logger.Info().Msgf("start delete")

	// if err := messageQueue.DeleteMessage("stream", dd); err != nil {
	// 	logger.Error().Err(err).Msg("Failed to delete messages")
	// }
}

type DeletorConfig struct {
	Nats struct {
		Url string `yaml:"url"`
	} `yaml:"nats"`
}

func loadConfig(config *utils.Config) (*DeletorConfig, error) {
	var cfg DeletorConfig
	if err := config.Cfg.UnmarshalKey("config", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

type Agent interface {
	Start(stream, subject string) error
	Stop()
}
