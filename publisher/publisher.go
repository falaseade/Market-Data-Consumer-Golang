package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/falaseade/Market-Data-Consumer-Golang/transformer"
	"github.com/nats-io/nats.go/jetstream"
)

type Publisher interface {
	Publish(ctx context.Context, payload []byte) error
}

type NatsPublisher struct {
	js jetstream.JetStream
	transformer transformer.Transformer
}

func NewNatsPublisher(js jetstream.JetStream, transformer transformer.Transformer) *NatsPublisher {
	return &NatsPublisher{
		js: js,
		transformer: transformer,
	}
}

func (p *NatsPublisher) Publish(ctx context.Context, payload []byte) error {
	canonicalTick, err := p.transformer.Transform(payload)
	if err != nil {
		return fmt.Errorf("failed to transform message: %w", err)
	}

	msgToSend, err := json.Marshal(canonicalTick)
	if err != nil {
		return fmt.Errorf("failed to marshal canonical tick: %w", err)
	}

	subject := os.Getenv("STREAM_SUBJECT") + canonicalTick.Symbol
	_, err = p.js.Publish(ctx, subject, msgToSend)
	if err != nil {
		return fmt.Errorf("failed to publish to NATS: %w", err)
	}

	log.Printf("Published to %s", subject)
	return nil
}