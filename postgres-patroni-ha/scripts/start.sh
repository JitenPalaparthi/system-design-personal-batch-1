#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../docker"
docker compose up -d --build
echo 'Waiting for Patroni cluster...'
for i in {1..60}; do
  if docker compose exec -T postgres1 patronictl -c /etc/patroni/patroni.yml list 2>/dev/null | grep -q 'Leader'; then
    docker compose exec -T postgres1 patronictl -c /etc/patroni/patroni.yml list
    echo 'Connect: docker compose exec client psql -h haproxy -p 5000 -U postgres'
    exit 0
  fi
  sleep 2
done
echo 'Cluster did not become ready in time.' >&2
docker compose logs --tail=100 postgres1 postgres2 postgres3
exit 1
