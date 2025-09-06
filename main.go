package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
)

func main() {
	webhookUrl := flag.String("url", "wss://stream.binance.com:9443/ws/btcusdt@trade", "Websocket URL used")
	flag.Parse()

	go StartClient(*webhookUrl)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Println("WebSocket client started. Press Ctrl+C to exit.")
	
	<-interrupt

	log.Println("Interrupt signal received, shutting down.")
}