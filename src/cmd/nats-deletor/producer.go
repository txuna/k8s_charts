package main

var _ Agent = (*Producer)(nil)

type Producer struct {
}

func NewProducer() *Producer {
	return &Producer{}
}

func (p *Producer) Start(stream, subject string) error {
	return nil
}

func (p *Producer) Stop() {
}
