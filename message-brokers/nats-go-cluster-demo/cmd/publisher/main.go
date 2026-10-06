package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	urls := env("NATS_URLS", "nats://localhost:4222,nats://localhost:4223,nats://localhost:4224")
	subject := env("SUBJECT", "orders.created")

	nc, err := nats.Connect(
		strings.ReplaceAll(urls, ",", " "),
		nats.Name("go-publisher"),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) { log.Printf("DISCONNECTED: %v", err) }),
		nats.ReconnectHandler(func(c *nats.Conn) { log.Printf("RECONNECTED: %s", c.ConnectedUrl()) }),
	)
	if err != nil { log.Fatal(err) }
	defer nc.Drain()

	log.Printf("Connected to %s", nc.ConnectedUrl())
	for i := 1; i <= 100; i++ {
		msg := fmt.Sprintf(`{"order_id":%d,"created_at":%q}`, i, time.Now().Format(time.RFC3339))
		if err := nc.Publish(subject, []byte(msg)); err != nil { log.Fatal(err) }
		log.Printf("PUB %s -> %s", subject, msg)
		time.Sleep(time.Second)
	}
}
func env(k,v string) string { if x:=os.Getenv(k); x!="" { return x }; return v }
