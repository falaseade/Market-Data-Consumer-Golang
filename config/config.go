package config

import (
	"os"
	"strings"
)

type Config struct {
	WebhookURL string
	NatsUrl string
	Symbols []string
}

func SetupConfig()(Config, error){
	return Config{
		WebhookURL: os.Getenv("WEBHOOK_URL"),
		NatsUrl: os.Getenv("NATS_URL"),
		Symbols: createSymbolString(os.Getenv("SYMBOLS")),
	}, nil
}

func createSymbolString(symbol string) []string {
	var symbolString []string
	for s := range strings.SplitSeq(symbol, ",") {
        symbolString = append(symbolString, strings.ToUpper(strings.TrimSpace(s)))
	}
	return symbolString
}