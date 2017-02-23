#!/bin/bash

set -e

CWD=`pwd`
DIR=`dirname $0`
MYSQLDUMP="$CWD/$DIR/mysqldump"
#echo $MYSQLDUMP

TIMESTAMP=`date "+%Y%m%d_%H%M%S"`
$MYSQLDUMP -u$MYSQL_USER -p$MYSQL_PASSWORD -h$MYSQL_HOST $MYSQL_DATABASE | bzip2 -9 > lekcije_$TIMESTAMP.dump.bz2
aws s3 cp lekcije_$TIMESTAMP.dump.bz2 s3://lekcije.com/backup/
