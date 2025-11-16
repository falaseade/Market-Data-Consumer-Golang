package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/falaseade/Market-Data-Consumer-Golang/config"
	"github.com/falaseade/Market-Data-Consumer-Golang/publisher"
	"github.com/falaseade/Market-Data-Consumer-Golang/transformer"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	cfg, err := config.SetupConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	if len(cfg.Symbols) == 0 {
		log.Fatalf("Config error: no symbols configured")
	}

	streamURL := cfg.WebhookURL + buildStreamSuffix(cfg.Symbols)
	fmt.Println(streamURL)
	webhookURLFlag := flag.String("url", streamURL, "WebSocket URL used")
	flag.Parse()
	webhookURL := *webhookURLFlag

	log.Printf("Using WebSocket URL: %s", webhookURL)
	log.Printf("NATS URL: %s", cfg.NatsUrl)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	nc, err := nats.Connect(cfg.NatsUrl)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer func() {
		log.Println("Draining NATS connection...")
		nc.Drain()
	}()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatalf("Error creating JetStream context: %v", err)
	}

	jsCfg, err := config.SetupJetstreamConfig()
	if err != nil {
		log.Fatalf("Failed to load JetStream configuration: %v", err)
	}
	log.Printf("Ensuring JetStream stream %q exists", jsCfg.Name)

	if _, err := js.Stream(ctx, jsCfg.Name); err != nil {
		log.Printf("Stream %q not found, creating it...", jsCfg.Name)
		if _, err := js.CreateStream(ctx, jsCfg); err != nil {
			log.Fatalf("Failed to create JetStream stream %q: %v", jsCfg.Name, err)
		}
	} else {
		log.Printf("Stream %q already exists", jsCfg.Name)
	}

	t := transformer.NewBinanceTransformer(cfg.Symbols)
	pub, err := publisher.NewNatsPublisher(js, t)
	if err != nil {
		log.Fatalf("Failed to create NATS publisher: %v", err)
	}

	go func() {
		log.Println("Starting WebSocket client...")
		StartClient(ctx, webhookURL, pub)
		log.Println("WebSocket client goroutine exited")
	}()

	log.Println("Service started. Press Ctrl+C to exit.")

	<-ctx.Done()
	log.Println("Shutdown signal received, cancelling context...")

	time.Sleep(1 * time.Second)
	log.Println("Application exited.")
}

func buildStreamSuffix(symbols []string) string {
	parts := make([]string, len(symbols))
	for i, s := range symbols {
		parts[i] = strings.ToLower(s) + "@trade"
	}
	return strings.Join(parts, "/")
}
