
case $1 in
loc)
    srv_ver=`git rev-parse --abbrev-ref HEAD`
    docker run --rm --name gex-loc -p 3831:3831 -it \
        --link postgres:psql.loc --link redis:redis.loc \
        gexservice:$srv_ver /app/gexservice/conf/local.properties
;;
esac