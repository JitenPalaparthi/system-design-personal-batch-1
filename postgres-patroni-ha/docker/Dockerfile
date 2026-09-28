FROM postgres:17-bookworm
USER root
RUN apt-get update && apt-get install -y --no-install-recommends python3 python3-pip python3-psycopg2 curl ca-certificates && \
    pip3 install --break-system-packages --no-cache-dir 'patroni[etcd3]' && rm -rf /var/lib/apt/lists/*
USER postgres
ENTRYPOINT ["patroni", "/etc/patroni/patroni.yml"]
