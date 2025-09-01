package main

import (
	"flag"
	"log"

	"github.com/gorilla/websocket"
)

var dialer = websocket.Dialer{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func main(){
	
	webhookUrl := flag.String("url", "wss://stream.binance.com:9443/ws/btcusdt@trade", "Websocket URL used")
	flag.Parse()

	conn, _, err := dialer.Dial(*webhookUrl, nil)
	if err != nil {
		log.Fatal("dial: ", err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil{
			log.Printf("Error: %s", err)
		}
		log.Printf("Message Type: %d recv: %s", messageType, p)
		
	}
}