package config

import (
	"errors"
	"os"
)

type Config struct {
	WebhookURL string
	NatsUrl string
}

func SetupConfig()(Config, error){
	webhookUrl := os.Getenv("WEBHOOK_URL")
	natsUrl := os.Getenv("NATS_URL")
	if webhookUrl == "" {
		return Config{}, errors.New("environment variable WEBHOOK_URL is not set, must be set")
	}
	if natsUrl == "" {
		return Config{}, errors.New("environment variable NATS_URL is not set, must be set")
	}
	config := Config{
		WebhookURL: webhookUrl,
		NatsUrl: natsUrl,
	}
	return config, nil
}