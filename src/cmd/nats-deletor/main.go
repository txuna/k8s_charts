package main

import "main/pkg/logger"

func main() {
	logger.InitLogger()
	logger.Info().Msg("Hello World")
}
