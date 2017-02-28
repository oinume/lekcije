#!/bin/bash

set -e

CWD=`pwd`
DIR=`dirname $0`
MYSQLDUMP="$CWD/$DIR/mysqldump"
#echo $MYSQLDUMP

DATE=`date "+%Y%m%d"`
$MYSQLDUMP -u$MYSQL_USER -p$MYSQL_PASSWORD -h$MYSQL_HOST $MYSQL_DATABASE | bzip2 -9 > lekcije_$DATE.dump.bz2
aws s3 cp lekcije_$DATE.dump.bz2 s3://lekcije.com/backup/

DELETE_DATE=`date "+%Y%m%d" --date "7 days ago"`
aws s3 ls s3://lekcije.com/backup/lekcije_$DELETE_DATE.dump.bz2 && aws s3 rm s3://lekcije.com/backup/lekcije_$DELETE_DATE.dump.bz2
