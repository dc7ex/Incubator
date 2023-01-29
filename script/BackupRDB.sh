#!/usr/bin/env bash

mysqlRootPass=2abfd8d3c3120%^d7e3d5
timestamp=$(date +"%Y-%m-%d.%H-%M-%S.%z")
scriptPath=$(cd `dirname $0`; pwd)

list=(365ex_collection 365ex_passport 365ex_pay bluesea_passport bluesea_pay fugu fugu_financial_admin fugu_passport fugu_pay goods_admin lion lion_goods)
for i in ${list[*]}; do

dbName="$i"
outputFile=${scriptPath}/${timestamp}.${dbName}.sql.gz

mysqldump --add-drop-database --set-gtid-purged=off --column-statistics=0 -h rm-uf68bi0zqg0g9212i.mysql.rds.aliyuncs.com -u root -p${mysqlRootPass} -B ${dbName} | gzip > $outputFile
echo "Data backup to complete ${outputFile}"
find ${scriptPath}/ -mtime +7 -name "*.${dbName}.sql.gz" -exec rm -f {} \;

done
