package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Broker struct {
	Notifier       chan []byte
	newClients     chan chan []byte
	closingClients chan chan []byte
	clients        map[chan []byte]bool
	mu             sync.Mutex
}

func NewBroker() *Broker {
	b := &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}
	go b.listen()
	return b
}

func (b *Broker) listen() {
	for {
		select {
		case s := <-b.newClients:
			b.mu.Lock()
			b.clients[s] = true
			b.mu.Unlock()
		case s := <-b.closingClients:
			b.mu.Lock()
			delete(b.clients, s)
			b.mu.Unlock()
		case event := <-b.Notifier:
			b.mu.Lock()
			for clientChan := range b.clients {
				select {
				case clientChan <- event:
					// Success
				default:
					// Client too slow, drop message
					log.Println("Skipping slow client")
				}
			}
			b.mu.Unlock()
		}
	}
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // For Nginx/Render proxies

	messageChan := make(chan []byte)
	b.newClients <- messageChan
	defer func() {
		b.closingClients <- messageChan
	}()

	for {
		select {
		case msg := <-messageChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (b *Broker) BroadcastUpdate(data interface{}) {
	payload, _ := json.Marshal(data)
	b.Notifier <- payload
}
