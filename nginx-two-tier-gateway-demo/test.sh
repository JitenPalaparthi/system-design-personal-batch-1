#!/usr/bin/env sh
set -eu
BASE=${BASE:-http://localhost:8080}

echo "== Gateway health =="
curl -fsS "$BASE/gateway-health"; echo

echo "\n== Go through gateway (observe go-1/go-2) =="
for i in 1 2 3 4 5 6; do curl -fsS "$BASE/api/go/whoami"; echo; done

echo "\n== Python through gateway (observe python-1/python-2) =="
for i in 1 2 3 4 5 6; do curl -fsS "$BASE/api/python/whoami"; echo; done

echo "\n== Go CRUD =="
curl -fsS "$BASE/api/go/items"; echo
curl -fsS -X POST "$BASE/api/go/items" -H 'Content-Type: application/json' -d '{"name":"created-through-gateway"}'; echo

echo "\n== Python CRUD =="
curl -fsS "$BASE/api/python/items"; echo
curl -fsS -X POST "$BASE/api/python/items" -H 'Content-Type: application/json' -d '{"name":"created-through-gateway"}'; echo

echo "\n== Directly test downstream load balancers =="
curl -fsS http://localhost:8081/whoami; echo
curl -fsS http://localhost:8082/whoami; echo

echo "\nSmoke test complete."
