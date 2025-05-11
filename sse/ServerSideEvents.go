package sse

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Event interface {
	Key() string
	Json() []byte
	Type() string
}

type Client chan Event

type SSE struct {
	clients      map[Client]bool
	clientsMutex sync.Mutex
}

func Init() *SSE {
	sse := SSE{}
	sse.clients = make(map[Client]bool)
	sse.clientsMutex = sync.Mutex{}
	return &sse
}

// Gin Handler
func (sse *SSE) Handler(c *gin.Context) {
	w := c.Writer
	req := c.Request

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Last-Event-ID")

	flusher, ok := w.(http.Flusher)
	if !ok {
		c.String(http.StatusInternalServerError, "Streaming unsupported")
		return
	}

	ctx := req.Context()

	heartbeatTicker := time.NewTicker(5 * time.Second)
	defer heartbeatTicker.Stop()

	// Create client channel and register it
	clientChan := make(Client, 100)
	sse.clientsMutex.Lock()
	sse.clients[clientChan] = true
	sse.clientsMutex.Unlock()

	defer func() {
		sse.clientsMutex.Lock()
		delete(sse.clients, clientChan)
		sse.clientsMutex.Unlock()
		close(clientChan)
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("Client disconnected")
			return

		// Heartbeats
		case <-heartbeatTicker.C:
			fmt.Fprint(w, "event: heartbeat\n")
			fmt.Fprintf(w, "data: %s\n\n", time.Now().Format(time.RFC3339))
			flusher.Flush()

		case log := <-clientChan:
			fmt.Println(log)
			sendEvent(w, log)
			flusher.Flush()
		}
	}
}

func sendEvent(w http.ResponseWriter, event Event) {
	fmt.Fprintf(w, "event: %s\n", event.Type())
	fmt.Fprintf(w, "id: %s\n", event.Key())
	fmt.Fprintf(w, "data: %s\n\n", event.Json())
}

func (sse *SSE) Broadcast(evt Event) {
	sse.clientsMutex.Lock()
	defer sse.clientsMutex.Unlock()
	for client := range sse.clients {
		select {
		case client <- evt:
		default:
			delete(sse.clients, client)
			close(client)
		}
	}
}
