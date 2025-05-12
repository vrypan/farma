package natsclient

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type Event interface {
	Key() string
	Json() []byte
	Topic() string
}

var (
	natsEnabled bool
	conn        *nats.Conn
	once        sync.Once
)

func Publish(e Event) error {
	if !natsEnabled {
		return nil
	}
	topic := fmt.Sprintf("farma.%s", e.Topic())
	err := conn.Publish(topic, e.Json())
	if err != nil {
		log.Printf("Error %v", err)
	}
	return err
}

func Connect(url string, options ...nats.Option) error {
	var err error
	natsEnabled = true
	once.Do(func() {
		defaultOptions := []nats.Option{
			nats.Name("GoService"),
			nats.ReconnectWait(2 * time.Second),
			nats.MaxReconnects(10),
			nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
				log.Printf("NATS disconnected: %v", err)
			}),
			nats.ReconnectHandler(func(nc *nats.Conn) {
				log.Printf("NATS reconnected to %v", nc.ConnectedUrl())
			}),
			nats.ClosedHandler(func(nc *nats.Conn) {
				log.Printf("NATS connection closed")
			}),
		}
		defaultOptions = append(defaultOptions, options...)
		conn, err = nats.Connect(url, defaultOptions...)
	})
	return err
}

func Conn() *nats.Conn {
	return conn
}
