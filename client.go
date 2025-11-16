package main

import (
	"context"
	"log"
	"sync"

	"github.com/falaseade/Market-Data-Consumer-Golang/publisher"
	"github.com/gorilla/websocket"
)

var dialer websocket.Dialer

func StartClient(ctx context.Context, url string, pub publisher.Publisher) {
	var clientWaitGroup sync.WaitGroup
	log.Println("WebSocket client connecting to:", url)
	natsQueue := make(chan []byte, 100)

	var once sync.Once
	closeQueue := func() {
		once.Do(func() {
			close(natsQueue)
		})
	}

	clientWaitGroup.Add(1)
	go func() {
		defer clientWaitGroup.Done()
		for message := range natsQueue {
			if err := pub.Publish(ctx, message); err != nil {
				log.Printf("Error publishing message: %v", err)
			}
		}
	}()

	conn, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		log.Printf("Failed to dial websocket: %v", err)
		closeQueue()
		clientWaitGroup.Wait()
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
			_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			closeQueue()
			clientWaitGroup.Wait()
			return

		case msg, ok := <-msgChan:
			if !ok {
				log.Println("Message channel closed.")
				closeQueue()
				clientWaitGroup.Wait()
				return
			}
			natsQueue <- msg

		case err := <-errChan:
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("Error reading from websocket: %v", err)
			}
			log.Println("WebSocket connection closed.")
			closeQueue()
			clientWaitGroup.Wait()
			return
		}
	}
}
