# Two-Tier NGINX Gateway + Load Balancer Demo

A training lab with the same two CRUD applications (Go and Python), but with NGINX split into distinct infrastructure roles.

## Architecture

```text
                                Docker network

Client
  |
  | http://localhost:8080
  v
+-----------------------------+
| NGINX #1 - API GATEWAY      |
| gateway                     |
|                             |
| /api/go/*     -> go-lb      |
| /api/python/* -> python-lb  |
| rate limiting               |
| gateway headers             |
| path prefix removal         |
+---------------+-------------+
                |
       +--------+--------+
       |                 |
       v                 v
+---------------+  +------------------+
| NGINX #2      |  | NGINX #3         |
| Go LB/Proxy   |  | Python LB/Proxy  |
| host :8081    |  | host :8082       |
+-------+-------+  +---------+--------+
        |                    |
    round robin          round robin
     /       \            /       \
    v         v          v         v
+-------+ +-------+  +--------+ +--------+
| go1   | | go2   |  |python1 | |python2 |
| :8080 | | :8080 |  | :5000  | | :5000  |
+-------+ +-------+  +--------+ +--------+
```

Only the NGINX containers expose host ports. Application containers remain internal.

## What each NGINX demonstrates

### NGINX #1: API Gateway

Public entry point: `http://localhost:8080`

Responsibilities:
- Single public endpoint
- Path-based API routing
- `/api/go/*` -> Go service tier
- `/api/python/*` -> Python service tier
- Removes the gateway prefix before forwarding
- Rate limiting
- Forwarded client headers
- Request ID propagation
- Gateway response header

### NGINX #2: Go reverse proxy/load balancer

Direct training endpoint: `http://localhost:8081`

Responsibilities:
- Reverse proxy for Go
- Load balancing across `go1` and `go2`
- Default NGINX round-robin algorithm
- Proxy headers and upstream logging

### NGINX #3: Python reverse proxy/load balancer

Direct training endpoint: `http://localhost:8082`

Responsibilities are equivalent to the Go load balancer, targeting `python1` and `python2`.

## Applications

Both applications expose the same API shape:

| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/items` | List items |
| POST | `/items` | Create item |
| GET | `/items/{id}` | Read one item |
| PUT | `/items/{id}` | Update item |
| DELETE | `/items/{id}` | Delete item |
| GET | `/health` | Health check |
| GET | `/whoami` | Show the replica that handled the request |

Data is deliberately in-memory. Each replica has its own independent state. This is useful for demonstrating why stateless services normally use shared persistence when horizontally scaled.

## Prerequisites

- Docker Engine or Docker Desktop
- Docker Compose v2 (`docker compose`)
- curl for command-line testing

## Start

```bash
docker compose up --build -d
docker compose ps
```

Watch all logs:

```bash
docker compose logs -f
```

Stop and remove containers/networks:

```bash
docker compose down
```

## Demo 1 - API Gateway routing

```bash
curl -i http://localhost:8080/gateway-health
curl -i http://localhost:8080/api/go/health
curl -i http://localhost:8080/api/python/health
```

Notice the `X-Gateway: nginx-api-gateway` response header.

The gateway changes the externally visible path:

```text
/api/go/whoami     -> go-lb:80/whoami
/api/python/whoami -> python-lb:80/whoami
```

## Demo 2 - Load balancing

Run multiple requests through the gateway:

```bash
for i in {1..8}; do curl -s http://localhost:8080/api/go/whoami; echo; done
for i in {1..8}; do curl -s http://localhost:8080/api/python/whoami; echo; done
```

You should see requests distributed between `go-1`/`go-2` and `python-1`/`python-2`.

Now bypass the gateway and call the load balancers directly:

```bash
for i in {1..6}; do curl -s http://localhost:8081/whoami; echo; done
for i in {1..6}; do curl -s http://localhost:8082/whoami; echo; done
```

This proves that load balancing is performed by the downstream NGINX instances, independently of the gateway.

## Demo 3 - CRUD through the gateway

### Go

```bash
curl http://localhost:8080/api/go/items

curl -X POST http://localhost:8080/api/go/items \
  -H 'Content-Type: application/json' \
  -d '{"name":"Laptop"}'

curl http://localhost:8080/api/go/items/1

curl -X PUT http://localhost:8080/api/go/items/1 \
  -H 'Content-Type: application/json' \
  -d '{"name":"Updated Go Item"}'

curl -i -X DELETE http://localhost:8080/api/go/items/1
```

### Python

```bash
curl http://localhost:8080/api/python/items

curl -X POST http://localhost:8080/api/python/items \
  -H 'Content-Type: application/json' \
  -d '{"name":"Phone"}'

curl http://localhost:8080/api/python/items/1

curl -X PUT http://localhost:8080/api/python/items/1 \
  -H 'Content-Type: application/json' \
  -d '{"name":"Updated Python Item"}'

curl -i -X DELETE http://localhost:8080/api/python/items/1
```

Because replicas use separate in-memory stores, a write sent to one replica is not guaranteed to be visible to the next request if it reaches the other replica. This behavior is intentional in this infrastructure lab.

## Demo 4 - Observe the request path

Open three terminals:

```bash
# Terminal 1
docker compose logs -f gateway

# Terminal 2
docker compose logs -f go-lb

# Terminal 3
docker compose logs -f go1 go2
```

Then call:

```bash
curl http://localhost:8080/api/go/whoami
```

The logs show the request moving through:

```text
Client -> Gateway -> Go NGINX LB -> Go replica
```

For Python:

```bash
docker compose logs -f gateway python-lb python1 python2
```

## Demo 5 - Failure behavior

Stop one Go replica:

```bash
docker compose stop go1
```

Try several requests:

```bash
for i in {1..6}; do curl -s http://localhost:8080/api/go/whoami; echo; done
```

Restart it:

```bash
docker compose start go1
```

This is a useful demonstration of upstream availability and proxy/load-balancer behavior.

## Automated smoke test

```bash
chmod +x test.sh
./test.sh
```

## Important training distinction

**Reverse proxy** describes forwarding client requests to backend servers while hiding the backend topology.

**Load balancer** distributes requests among multiple backend instances.

**API gateway** is the API-facing entry layer. In this lab it performs API path routing, prefix rewriting, rate limiting, request metadata propagation, and presents a single public API endpoint.

NGINX is capable of all these roles. We use separate NGINX containers so students can see that the *role* comes from configuration and placement in the architecture, not merely from the product name.

## Project structure

```text
nginx-two-tier-gateway-demo/
├── docker-compose.yml
├── README.md
├── test.sh
├── go-service/
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go
├── python-service/
│   ├── Dockerfile
│   ├── requirements.txt
│   └── app.py
├── nginx-gateway/
│   └── nginx.conf
├── nginx-go-lb/
│   └── nginx.conf
└── nginx-python-lb/
    └── nginx.conf
```
