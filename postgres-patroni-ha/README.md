# PostgreSQL 17 + Patroni + etcd + HAProxy HA Lab

A three-node PostgreSQL high-availability lab. Patroni manages PostgreSQL roles/failover, a 3-member etcd cluster is the DCS, HAProxy exposes the current primary on port 5000, and a PostgreSQL client container is included.

## Docker quick start

Requirements: Docker Engine/Desktop with Docker Compose v2.

```bash
cd scripts
./start.sh
./status.sh
./demo.sh
./client.sh
```

From `docker/` you can also connect directly:

```bash
docker compose exec client psql -h haproxy -p 5000 -U postgres
```

Password: `postgres123`.

HAProxy statistics: http://localhost:7000/ (lab only; no authentication).

## Failover demonstration

```bash
cd scripts
./failover-demo.sh
```

The script identifies the current leader, stops it, waits for Patroni to elect/promote another node, and verifies a new SQL connection through HAProxy. Existing TCP sessions can break during failover; production applications need reconnect/retry logic.

Bring the stopped node back with, for example:

```bash
cd docker
docker compose start postgres1 postgres2 postgres3
```

Only the stopped service will need starting; already-running services are unaffected.

## Reset everything

```bash
cd scripts
./cleanup.sh
```

WARNING: cleanup removes all lab database and etcd volumes.

## Architecture

```text
client -> HAProxy:5000 -> current PostgreSQL primary:5432
                    |-> Patroni REST /primary checks:8008

postgres1 ----\
postgres2 -----+-- Patroni + PostgreSQL
postgres3 ----/
     |
3-member etcd DCS/quorum
```

## VM layout

Recommended lab layout:

- 192.168.56.10: HAProxy + psql client
- 192.168.56.11: PostgreSQL + Patroni + etcd member 1
- 192.168.56.12: PostgreSQL + Patroni + etcd member 2
- 192.168.56.13: PostgreSQL + Patroni + etcd member 3

Copy `vm/patroni-template.yml` to `/etc/patroni/patroni.yml` on each DB VM and replace `NODE_NAME` / `NODE_IP`. Configure etcd on each node with a unique member name and the three-member initial cluster. Stop/disable the distro PostgreSQL service before starting Patroni so Patroni exclusively manages the database instance. Copy `vm/haproxy.cfg` to `/etc/haproxy/haproxy.cfg` on the proxy VM and restart HAProxy.

Client test from 192.168.56.10:

```bash
PGPASSWORD=postgres123 psql -h 127.0.0.1 -p 5000 -U postgres -c 'select inet_server_addr(), pg_is_in_recovery();'
```

## Production notes

This is an educational HA lab, not a hardened production deployment. Replace plaintext secrets, enable TLS/authentication for Patroni/etcd, configure firewalls, backups/WAL archiving, monitoring, fencing/watchdogs as appropriate, and redundant HAProxy endpoints/VIP or another highly available ingress. Test RPO/RTO and failure modes before production use.
