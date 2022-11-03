#!/bin/sh
set -xe
# copy default configure
if [ ! -f /app/conf/gexservice.properties ];then
    mkdir -p /app/conf
    cp -f /app/gexservice/conf/gexservice.properties /app/conf/gexservice.properties
fi

cd /app/gexservice/
/app/gexservice/service /app/conf/gexservice.properties
echo "service is done"
