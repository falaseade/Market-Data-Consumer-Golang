package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	webhookUrl := flag.String("url", "wss://stream.binance.com:9443/ws/btcusdt@trade", "Websocket URL used")
	flag.Parse()

	natsUrl := os.Getenv("NATS_URL")

	nc, natsErr := nats.Connect(natsUrl)
	if natsErr != nil{
		log.Fatal("Error connection to nats", natsErr)
	}

	defer nc.Close()

	js, jetstreamErr := jetstream.New(nc)
	if jetstreamErr != nil {
		log.Fatal("Error with JetStream", jetstreamErr)
	}

	streamName := os.Getenv("JS_STREAM_NAME")
	streamSubjects := os.Getenv("JS_STREAM_SUBJECTS")
	streamRetentionTime, retentionError := strconv.Atoi(os.Getenv("JS_RETENTION_TIME"))
	if retentionError != nil {
		streamRetentionTime = 24
	}

	cfg :=  jetstream.StreamConfig{
		Name: streamName,
		Subjects: []string{streamSubjects},
		MaxAge:  time.Duration(streamRetentionTime) * time.Hour,

	}

	_, streamError := js.CreateStream(ctx, cfg)
	if streamError != nil {
		log.Fatal("Stream Error!!!")
	}

	
go StartClient(*webhookUrl)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Println("WebSocket client started. Press Ctrl+C to exit.")
	
	<-interrupt

	log.Println("Interrupt signal received, shutting down.")
}