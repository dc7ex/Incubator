#!/bin/bash
 
getdatestring()
{
   TZ='Asia/Chongqing' date "+%Y%m%d%H%M"
}
datestring=$(getdatestring)

LOGS_PATH=/var/log/nginx/7ex

list=(api open passport stars www)
for i in ${list[*]}; do

name="$i"

find ${LOGS_PATH}/ -mtime +7 -name "${name}.access.20*.log" -exec rm -f {} \;
find ${LOGS_PATH}/ -mtime +7 -name "${name}.error.20*.log" -exec rm -f {} \;
 
mv ${LOGS_PATH}/${name}.access.log ${LOGS_PATH}/${name}.access.${datestring}.log
mv ${LOGS_PATH}/${name}.error.log ${LOGS_PATH}/${name}.error.${datestring}.log

done

kill -USR1 `cat /var/run/nginx.pid` 
