# NATS 3-Node Cluster + Go Demo

This lab demonstrates Core NATS clustering with Go:

- 3 NATS servers
- Publisher / Subscriber
- Pub/Sub fan-out
- Queue groups / load balancing
- Request / Reply
- Go client reconnection when a NATS node is stopped

## 1. Start the cluster

```bash
docker compose up -d
docker compose ps
```

Client ports:

- nats1: localhost:4222
- nats2: localhost:4223
- nats3: localhost:4224

Monitoring:

- http://localhost:8222
- http://localhost:8223
- http://localhost:8224

## 2. Download Go dependencies

```bash
go mod tidy
```

## 3. Normal Pub/Sub

Open three terminals:

```bash
go run ./cmd/subscriber -name consumer-1
go run ./cmd/subscriber -name consumer-2
go run ./cmd/subscriber -name consumer-3
```

Then:

```bash
go run ./cmd/publisher
```

All three subscribers should receive each publication.

## 4. Queue Group

Stop the previous subscribers and start:

```bash
go run ./cmd/subscriber -name worker-1 -queue order-workers
go run ./cmd/subscriber -name worker-2 -queue order-workers
go run ./cmd/subscriber -name worker-3 -queue order-workers
```

Then:

```bash
go run ./cmd/publisher
```

Only one member of `order-workers` receives each message.

## 5. Request / Reply

Start three service instances:

```bash
go run ./cmd/service -name service-1
go run ./cmd/service -name service-2
go run ./cmd/service -name service-3
```

Then run:

```bash
go run ./cmd/requester
```

Requests are load-balanced among the service instances. Each service replies using the reply subject carried by the request.

## 6. Cluster Failover Demo

Keep a publisher and consumers running. In another terminal:

```bash
docker compose stop nats1
```

Watch the Go clients reconnect to another discovered/configured server.

Bring it back:

```bash
docker compose start nats1
```

You can also stop nats2 or nats3 to demonstrate continued service while other cluster members remain available.

## 7. Stop Everything

```bash
docker compose down
```

## Important Teaching Point

This is Core NATS, not JetStream. The cluster provides routing and connection resiliency, but this demo does not provide durable message storage. If a subscriber is offline, Core NATS does not retain its messages for later replay. JetStream is the NATS persistence layer.
