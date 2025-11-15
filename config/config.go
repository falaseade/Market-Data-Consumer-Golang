package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	WebhookURL string
	NatsUrl string
	Symbols []string
}

func SetupConfig()(Config, error){
	webhookUrl := os.Getenv("WEBHOOK_URL")
	natsUrl := os.Getenv("NATS_URL")
	symbolsString := os.Getenv("SYMBOLS")
	if symbolsString == "" {
		return Config{}, errors.New("environment variable SYMBOLS is not set, must be set")
	}
	var symbols []string
	for s := range strings.SplitSeq(symbolsString, ",") {
        symbols = append(symbols, strings.ToUpper(strings.TrimSpace(s)))
	}
	
	
	if webhookUrl == "" {
		return Config{}, errors.New("environment variable WEBHOOK_URL is not set, must be set")
	}
	if natsUrl == "" {
		return Config{}, errors.New("environment variable NATS_URL is not set, must be set")
	}
	config := Config{
		WebhookURL: webhookUrl,
		NatsUrl: natsUrl,
		Symbols: symbols,
	}
	return config, nil
}