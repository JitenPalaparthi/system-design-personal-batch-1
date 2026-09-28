#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../docker"
docker compose exec -T postgres1 patronictl -c /etc/patroni/patroni.yml list
