#!/bin/sh
docker rm -f gexservice-srv
docker run -d \
    --name gexservice-srv \
    --link gexservice-postgres:psql.loc \
    --link gexservice-redis:redis.loc \
    -p 3831:3831 \
    -v /data/gexservice/conf:/app/conf \
    -v /data/gexservice/www:/app/www \
    -v /data/gexservice/upload:/app/upload \
    gexservice:$1
