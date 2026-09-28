#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../docker"
leader=$(docker compose exec -T postgres1 patronictl -c /etc/patroni/patroni.yml list -f json | python3 -c 'import json,sys; x=json.load(sys.stdin); print(next(n["Member"] for n in x if n["Role"] in ("Leader","Primary")))')
echo "Current leader: $leader"
docker compose stop "$leader"
echo 'Waiting for a new leader...'
for i in {1..30}; do
  sleep 2
  out=$(docker compose exec -T postgres2 patronictl -c /etc/patroni/patroni.yml list 2>/dev/null || true)
  if echo "$out" | grep -q 'Leader'; then
    echo "$out"
    docker compose exec -T client psql -h haproxy -p 5000 -U postgres -c "SELECT inet_server_addr(), pg_is_in_recovery(), now();"
    exit 0
  fi
done
echo 'No new leader elected in expected time.' >&2
exit 1
