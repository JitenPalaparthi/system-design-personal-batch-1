# HAProxy L4 + L7 Docker Demo — Python + Go

This training demo runs two HTTP services behind one HAProxy container:

- Python/Flask service: `python-service:8080`
- Go/net-http service: `go-service:8080`
- HAProxy L7 listener: host port `8080`
- HAProxy L4 listener: host port `9000`

## Architecture

```text
                         +--------------------+
HTTP :8080 (L7) -------->| HAProxy            |-- /python --> Python :8080
                         | HTTP/path aware     |-- /go     --> Go :8080
                         |                     |-- /       --> round-robin
                         +--------------------+

                         +--------------------+
TCP :9000 (L4) --------->| HAProxy            |-----------> Python :8080
                         | TCP stream aware    |-----------> Go :8080
                         | round-robin         |
                         +--------------------+
```

## Run

```bash
docker compose up --build
```

## L7 tests

HAProxy parses HTTP at Layer 7, so it can make routing decisions using the request path.

```bash
curl http://localhost:8080/python
curl http://localhost:8080/go
curl http://localhost:8080/
curl http://localhost:8080/
```

`/python` always selects Python, `/go` always selects Go, while `/` uses the `mixed_http` round-robin backend.

## L4 tests

Port 9000 is configured with `mode tcp`. HAProxy load-balances TCP connections without using HTTP paths for routing. The backend applications still happen to speak HTTP, so curl can be used as a convenient TCP client:

```bash
curl http://localhost:9000/
curl http://localhost:9000/
curl http://localhost:9000/
curl http://localhost:9000/
```

You should see responses alternate between `python-service` and `go-service` across new connections. Note that persistent/keep-alive connections can remain attached to the same selected backend because an L4 decision is made for the TCP connection.

## Key training point

**L4 (`mode tcp`)**: HAProxy operates on TCP connections. It can balance arbitrary TCP protocols, but cannot route based on HTTP URL paths in this configuration.

**L7 (`mode http`)**: HAProxy understands HTTP and can inspect paths, headers, methods, cookies, etc. This example routes `/python` and `/go` to different services.

## Useful commands

```bash
# Follow HAProxy logs
docker compose logs -f haproxy

# Show containers
docker compose ps

# Rebuild after source changes
docker compose up --build

# Stop and remove demo containers/network
docker compose down
```

## Files

```text
haproxy-l4-l7-demo/
├── docker-compose.yml
├── README.md
├── haproxy/
│   └── haproxy.cfg
├── python-service/
│   ├── app.py
│   ├── requirements.txt
│   └── Dockerfile
└── go-service/
    ├── go.mod
    ├── main.go
    └── Dockerfile
```
