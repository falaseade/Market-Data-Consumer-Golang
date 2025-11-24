package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/falaseade/Market-Data-Consumer-Golang/publisher"
	"github.com/gorilla/websocket"
)

func StartClient(ctx context.Context, url string, pub publisher.Publisher) error {
	log.Println("WebSocket client connecting to:", url)

	natsQueue := make(chan []byte, 100)

	var wg sync.WaitGroup

	wg.Go(func() {
		for msg := range natsQueue {
			if err := pub.Publish(ctx, msg); err != nil {
				log.Printf("error publishing message: %v", err)
			}
		}
	})

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		close(natsQueue)
		wg.Wait()
		return fmt.Errorf("dial websocket: %w", err)
	}
	defer conn.Close()
	log.Println("WebSocket client connected.")

	msgChan := make(chan []byte)
	errChan := make(chan error, 1)

	go func() {
		defer close(msgChan)

		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}

			select {
			case msgChan <- p:
			case <-ctx.Done():
				return
			}
		}
	}()

	defer func() {
		close(natsQueue)
		wg.Wait()
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("shutdown signal received, closing WebSocket client")

			_ = conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)

			return ctx.Err()

		case msg, ok := <-msgChan:
			if !ok {
				log.Println("message channel closed")
				return nil
			}

			select {
			case natsQueue <- msg:
			case <-ctx.Done():
				return ctx.Err()
			}

		case err := <-errChan:
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("error reading from websocket: %v", err)
			}
			log.Println("websocket connection closed")
			return err
		}
	}
}