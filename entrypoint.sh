#!/bin/sh
set -xe
# copy default configure
if [ ! -f /app/conf/gexservice.properties ];then
    mkdir -p /app/conf
    cp -f /app/gexservice/conf/gexservice.properties /app/conf/gexservice.properties
fi

cd /app/gexservice/
if [ "$1" != "" ];then
    /app/gexservice/service $1
else
    /app/gexservice/service /app/conf/gexservice.properties
fi
echo "service is done"
