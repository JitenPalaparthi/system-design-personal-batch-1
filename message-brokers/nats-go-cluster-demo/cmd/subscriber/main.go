package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	var name, queue string
	flag.StringVar(&name, "name", "consumer-1", "consumer name")
	flag.StringVar(&queue, "queue", "", "queue group; empty = normal pub/sub")
	flag.Parse()

	urls := env("NATS_URLS", "nats://localhost:4222,nats://localhost:4223,nats://localhost:4224")
	subject := env("SUBJECT", "orders.created")
	nc, err := nats.Connect(strings.ReplaceAll(urls, ",", " "), nats.Name(name),
		nats.MaxReconnects(-1), nats.ReconnectWait(time.Second),
		nats.ReconnectHandler(func(c *nats.Conn){ log.Printf("[%s] RECONNECTED -> %s",name,c.ConnectedUrl()) }))
	if err != nil { log.Fatal(err) }
	defer nc.Drain()

	handler := func(m *nats.Msg){ log.Printf("[%s] received subject=%s data=%s", name,m.Subject,string(m.Data)) }
	if queue == "" {
		_, err = nc.Subscribe(subject, handler)
		log.Printf("[%s] PUB/SUB on %s via %s",name,subject,nc.ConnectedUrl())
	} else {
		_, err = nc.QueueSubscribe(subject, queue, handler)
		log.Printf("[%s] QUEUE=%s on %s via %s",name,queue,subject,nc.ConnectedUrl())
	}
	if err != nil { log.Fatal(err) }
	if err=nc.Flush(); err!=nil { log.Fatal(err) }
	select {}
}
func env(k,v string) string { if x:=os.Getenv(k); x!="" { return x }; return v }
