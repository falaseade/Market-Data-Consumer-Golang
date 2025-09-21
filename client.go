package main

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
)

var dialer = websocket.Dialer{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartClient(ctx context.Context, url string) {
	log.Println("WebSocket client connecting to:", url)
	conn, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		log.Printf("Failed to dial websocket: %v", err)
		return
	}
	defer conn.Close()
	log.Println("WebSocket client connected.")

	msgChan := make(chan []byte)
    errChan := make(chan error, 1)

    go func() {
		defer close(msgChan)
        defer close(errChan)
        for {
            _, p, err := conn.ReadMessage()
            if err != nil {
                errChan <- err
                return
            }
            msgChan <- p
        }
    }()
for {
        select {
			case <-ctx.Done():
            log.Println("Shutdown signal received, closing WebSocket client.")
            err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
            if err != nil {
                log.Println("Error during websocket close:", err)
            }
            return
			
			case msg := <-msgChan:
            log.Printf("Received message: %s", msg)

			case err := <-errChan:
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
                log.Printf("Error reading from websocket: %v", err)
            }
            log.Println("WebSocket connection closed.")
            return
        }
    }
}