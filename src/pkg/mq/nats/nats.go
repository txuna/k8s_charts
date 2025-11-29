package mq

import (
	"context"
	"main/pkg/logger"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsMQ struct {
	nc *nats.Conn
	js jetstream.JetStream
}

func NewNats(addr string) (*NatsMQ, error) {
	n := &NatsMQ{}

	nc, err := nats.Connect(addr)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	n.nc = nc
	n.js = js

	return n, nil
}

func (n *NatsMQ) Close() {
	n.nc.Close()
}

func (n *NatsMQ) PurgeStream(stream, subject string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := n.js.Stream(ctx, stream)
	if err != nil {
		return err
	}

	return res.Purge(ctx, jetstream.WithPurgeSubject(subject))
}

func (n *NatsMQ) CreateStream(stream, subject string, replicas int, policy jetstream.RetentionPolicy) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := n.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:      stream,
		Subjects:  []string{subject},
		Replicas:  replicas,
		Retention: policy,
	})

	return err
}

func (n *NatsMQ) Publish(stream, subject, message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := n.js.Publish(ctx, subject, []byte(message), jetstream.WithExpectStream(stream))
	return err
}

func (n *NatsMQ) DeleteMessageLegacy(stream string, seq []uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, s := range seq {
		stream, err := n.js.Stream(ctx, stream)
		if err != nil {
			return err
		}

		if err := stream.DeleteMsg(ctx, s); err != nil {
			return err
		}
	}

	return nil
}

func (n *NatsMQ) DeleteMessage(stream string, seq []uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := n.js.Stream(ctx, stream)
	if err != nil {
		return err
	}

	for _, s := range seq {
		if err := res.DeleteMsg(ctx, s); err != nil {
			return err
		}
	}

	return nil
}

var count atomic.Int64

func (n *NatsMQ) Consume(stream, subject, durable string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cctx, err := n.js.CreateOrUpdateConsumer(ctx, stream, jetstream.ConsumerConfig{
		FilterSubject: subject,
		Durable:       "consumer",
		Name:          "consumer",
	})

	if err != nil {
		return err
	}

	cctx.Consume(func(msg jetstream.Msg) {
		msg.Ack()
		// logger.Info().Msgf("Received message: %s", string(msg.Data()))
		a := count.Add(1)
		logger.Info().Msgf("Total processed messages: %d", a)
	})

	return nil
}
