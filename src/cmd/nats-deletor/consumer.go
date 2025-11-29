package main

var _ Agent = (*Consumer)(nil)

type Consumer struct {
}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Start(stream, subject string) error {
	return nil
}

func (c *Consumer) Stop() {
}
