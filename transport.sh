#!/bin/bash
cur_ver=`git branch --show-current`
cur_ver=${cur_ver//./}
srv_pro=wss
srv_addr=$cur_ver.exservice.test.com
transport :5435 $srv_pro://exservice:123@$srv_addr/transport/pg?skip_verify=1 :6389 $srv_pro://exservice:123@$srv_addr/transport/redis?skip_verify=1
