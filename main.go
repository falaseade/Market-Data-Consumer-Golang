package main

import (
	"log"

	"github.com/gorilla/websocket"
)

var dialer = websocket.Dialer{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

var url = "wss://stream.binance.com:9443/ws/btcusdt@trade"

func main(){
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		log.Printf("Message Type: %d recv: %s", messageType, p)
		if err != nil{
			log.Printf("Error: %s", err)
		}
	}
}