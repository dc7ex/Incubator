#!/bin/bash

## 文件路径
FILES="/data/logs/nginx/365ex.art/collection.access.log"
DATE=`date -d '8 hours ago 1 minutes ago' +%Y:%H:%M`

# blackIps=`grep ${DATE} ${FILES}| grep -v "OPTIONS" | awk '{print $1}'|sort -n|uniq -c |sort -nr | awk '{if($1>100)print $2}'`

for blackIP in `grep ${DATE} ${FILES}| grep -v "OPTIONS" | awk '{print $1}'|sort -n|uniq -c |sort -nr | awk '{if($1>80)print $2}'`
do
  /usr/sbin/ipset add nginx_blacklist $blackIP
done
