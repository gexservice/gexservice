
case $1 in
loc)
    docker rm -f gex-loc
    docker run --name gex-loc -p 3831:3831 -d \
        --link postgres:psql.loc --link redis:redis.loc \
        gex.loc/gexservice:v1.0.0 /app/gexservice/conf/local.properties
;;
esac