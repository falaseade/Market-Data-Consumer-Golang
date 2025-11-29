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
	js          jetstream.JetStream
	transformer transformer.Transformer
}

func NewNatsPublisher(js jetstream.JetStream, t transformer.Transformer) (*NatsPublisher, error) {
	if js == nil {
		return nil, fmt.Errorf("jetstream instance is nil")
	}
	if t == nil {
		return nil, fmt.Errorf("transformer is nil")
	}
	return &NatsPublisher{
		js:          js,
		transformer: t,
	}, nil
}

func (p *NatsPublisher) Publish(ctx context.Context, payload []byte) error {
	if p == nil {
		return fmt.Errorf("publisher is nil")
	}
	if p.transformer == nil {
		return fmt.Errorf("transformer is nil")
	}
	if p.js == nil {
		return fmt.Errorf("jetstream client is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	evt, err := p.transformer.Transform(payload)
	if err != nil {
		return fmt.Errorf("transform error: %w", err)
	}

	if evt == nil {
		return nil
	}

	subject := fmt.Sprintf("%s.%s", os.Getenv("JS_STREAM_SUBJECTS"), evt.Symbol)

	msg, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("marshal canonical event: %w", err)
	}

	_, err = p.js.Publish(ctx, subject, msg, jetstream.WithMsgID(evt.MsgID))
	if err != nil {
		return fmt.Errorf("nats publish error: %w", err)
	}

	log.Printf("Published message to subject=%s msg_id=%s", subject, evt.MsgID)
	return nil
}
