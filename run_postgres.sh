#!/bin/sh
docker rm -f gexservice-postgres
docker run -d \
    --name gexservice-postgres \
    -e POSTGRES_DB=gexservice \
    -e POSTGRES_USER=gexservice \
    -e POSTGRES_PASSWORD=123 \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -e POSTGRES_HOST_AUTH_METHOD=md5 \
    -v /data/gexservice/postgres/:/var/lib/postgresql/data \
    postgres:latest
