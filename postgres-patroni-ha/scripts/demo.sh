#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../docker"
docker compose exec -T client psql -h haproxy -p 5000 -U postgres -v ON_ERROR_STOP=1 <<'SQL'
CREATE TABLE IF NOT EXISTS ha_demo(id bigserial primary key, created_at timestamptz default now(), server inet default inet_server_addr());
INSERT INTO ha_demo DEFAULT VALUES;
SELECT inet_server_addr() AS server, pg_is_in_recovery() AS is_replica;
SELECT * FROM ha_demo ORDER BY id DESC LIMIT 5;
SQL
