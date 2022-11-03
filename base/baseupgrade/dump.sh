#!/bin/bash

dump_clear(){
    ofile=$1
    used=sed
    if [ `uname` == "Darwin" ];then
        used=gsed
    fi
    $used -i 's/public\.//g' $ofile
    $used -i 's/ Owner: .*//g' $ofile
    $used -i 's/ALTER TABLE .* OWNER TO .*;//g' $ofile
    $used -i 's/-- Dumped from database.*//g' $ofile
    $used -i 's/-- Dumped by pg_dump.*//g' $ofile   
    $used -i ':a;N;$!ba;s/\n\n\n\n/\n\n/g' $ofile
    $used -i ':a;N;$!ba;s/START WITH 1\n/START WITH 1000\n/g' $ofile
    $used -i ':a;N;$!ba;s/SET[^=]*=[^\n]*\n//g' $ofile
    $used -i ':a;N;$!ba;s/ALTER TABLE ONLY [^\n]* DROP [^\n]*\n//g' $ofile
    $used -i ':a;N;$!ba;s/SELECT pg_catalog[^\n]*\n//g' $ofile
    $used -i 's/ALTER TABLE/ALTER TABLE IF EXISTS/g' $ofile
    $used -i 's/DROP INDEX/DROP INDEX IF EXISTS/g' $ofile
    $used -i 's/DROP SEQUENCE/DROP SEQUENCE IF EXISTS/g' $ofile
    $used -i 's/DROP TABLE/DROP TABLE IF EXISTS/g' $ofile
}

dump_all(){
    tmpfile=latest.sql
    ssh psql.loc "docker exec postgres pg_dump -s -c -U dev -d base -f /tmp/$tmpfile"
    ssh psql.loc "docker cp postgres:/tmp/$tmpfile /tmp/"
    scp psql.loc:/tmp/$tmpfile ./
    dump_clear $tmpfile


    cat > latest.go  << EOF
package baseupgrade

var LATEST = \`
EOF
    cat $tmpfile | grep -v DROP >> latest.go
    cat >> latest.go  << EOF
\` + INIT

EOF

    cat $tmpfile | grep DROP > clear.sql

    cat >> latest.go  << EOF
var DROP = \`
EOF

    cat clear.sql >> latest.go 
    cat >> latest.go  << EOF
\`

EOF

    cat >> latest.go  << EOF
const CLEAR = \`
EOF
    cat $tmpfile | grep 'DROP TABLE' | sed 's/DROP TABLE IF EXISTS/DELETE FROM/' >> latest.go

    cat >> latest.go  << EOF
\`
EOF

}

dump_upgrade(){
    tmpfile=emall_v$1.sql
    ver=v$1
    verup=V$1
    table=$2
    go run clear/clear.go postgresql://dev:123@psql.loc:5432/emservice
    ssh psql.loc "docker exec postgres pg_dump -s -c -U dev -d emservice -f /tmp/$tmpfile $table"
    ssh psql.loc "docker cp postgres:/tmp/$tmpfile /tmp/"
    scp psql.loc:/tmp/$tmpfile ./
    dump_clear $tmpfile
    cat > $ver.go  << EOF
package emsupgrade

const $verup = \`
EOF
    cat $tmpfile | grep -v DROP >> $ver.go
    cat >> $ver.go  << EOF
\`

EOF
}

case $1 in
upgrade)
     dump_upgrade $2 "$3"
;;
clear)
    dump_clear $2
;;
*)
    dump_all
;;
esac

