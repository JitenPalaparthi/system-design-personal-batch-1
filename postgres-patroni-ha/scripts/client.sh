#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../docker"
docker compose exec client psql -h haproxy -p 5000 -U postgres
