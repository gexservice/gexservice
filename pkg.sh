#!/bin/bash
##############################
#####Setting Environments#####
echo "Setting Environments"
set -xe
export cpwd=`pwd`
export LD_LIBRARY_PATH=/usr/local/lib:/usr/lib
output=$cpwd/build
#### Package ####
echo "Setting Environments"
set -xe
export cpwd=`pwd`
output=$cpwd/build
#### Package ####
srv_ver=`git rev-parse --abbrev-ref HEAD`
srv_name=gexservice

cat header.md > header_.md
item2md "Api Common Code" base/define/define.go >>  header_.md
item2md "Base Struct Define" base/basedb/items.go >>  header_.md
item2md "GEX Struct Define" gexdb/auto_models.go >>  header_.md
item2md "GEX External Define" gexdb/auto_external.go >>  header_.md
apidoc -c apidoc.json -i base/basedb -i base/baseapi -i gexdb -i matcher -i gexapi -o www/apidoc

docker build --build-arg=https_proxy=$HTTPS_PROXY --build-arg=http_proxy=$HTTPS_PROXY -t gexservice:$srv_ver .
if [ "$1" != "" ];then
    docker tag gexservice:$srv_ver $1/gexservice:$srv_ver
    docker push $1/gexservice:$srv_ver
fi

