#!/bin/sh
docker rm -f gexservice-redis
docker run -d \
    --name gexservice-redis \
    redis
