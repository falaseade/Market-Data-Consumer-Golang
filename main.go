package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/falaseade/Market-Data-Consumer-Golang/config"
	"github.com/falaseade/Market-Data-Consumer-Golang/publisher"
	"github.com/falaseade/Market-Data-Consumer-Golang/transformer"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	

	cfg, configError := config.SetupConfig()
	if configError != nil {
		log.Fatalf("Failed to load configuration: %v", configError)
	}

	webhookUrl := flag.String("url", cfg.WebhookURL, "Websocket URL used")
	flag.Parse()

	nc, natsErr := nats.Connect(cfg.NatsUrl)
	if natsErr != nil{
		log.Fatalf("Error connecting to nats: %v", natsErr)
	}
	
	defer nc.Close()

	js, jetstreamErr := jetstream.New(nc)
	if jetstreamErr != nil {
		log.Fatalf("Error with JetStream: %v", jetstreamErr)
	}

	jetstreamCfg, jetstreamCfgError := config.SetupJetstreamConfig()
	if jetstreamCfgError != nil {
		log.Fatalf("Failed to load jetstream configuration: %v", jetstreamCfgError)
	}

	_, streamError := js.CreateStream(ctx, jetstreamCfg)
	if streamError != nil {
		log.Fatalf("Failed to create JetStream stream: %v", streamError)
	}

	t := transformer.NewBinanceTransformer()
	pub := publisher.NewNatsPublisher(js, t)

	
go StartClient(ctx, *webhookUrl, pub)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Println("WebSocket client started. Press Ctrl+C to exit.")
	
	<-interrupt

	log.Println("Shutdown signal received, telling services to stop.")
    cancel()
	time.Sleep(1 * time.Second)
    log.Println("Application exited.")
}