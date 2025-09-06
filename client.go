package main

import (
	"log"

	"github.com/gorilla/websocket"
)

var dialer = websocket.Dialer{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartClient(url string) {
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		log.Printf("Message Type: %d recv: %s", messageType, p)
	}
}