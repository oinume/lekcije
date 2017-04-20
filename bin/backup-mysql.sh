#!/bin/bash

set -ex

CWD=`pwd`
DIR=`dirname $0`
MYSQLDUMP="$CWD/$DIR/mysqldump"
#echo $MYSQLDUMP
DATE=`date "+%Y%m%d"`
$MYSQLDUMP -u$MYSQL_USER -p$MYSQL_PASSWORD -h$MYSQL_HOST $MYSQL_DATABASE | bzip2 -9 > lekcije_$DATE.dump.bz2

if [ ! -d ./google-cloud-sdk ]; then
  curl https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-151.0.1-linux-x86_64.tar.gz \
    -o ./google-cloud-sdk.tar.gz
  tar xzf ./google-cloud-sdk.tar.gz
fi

echo $GCLOUD_SERVICE_KEY | base64 --decode > ./gcloud-service-key.json
./google-cloud-sdk/bin/gcloud \
  auth activate-service-account \
  --key-file ./gcloud-service-key.json

DELETE_DATE=`date "+%Y%m%d" --date "7 days ago"`
#aws s3 cp lekcije_$DATE.dump.bz2 s3://lekcije.com/backup/
#aws s3 ls s3://lekcije.com/backup/lekcije_$DELETE_DATE.dump.bz2 && aws s3 rm s3://lekcije.com/backup/lekcije_$DELETE_DATE.dump.bz2

gsutil cp lekcije_$DATE.dump.bz2 gs://lekcije/backup/
gsutil ls gs://lekcije/backup/lekcije_$DELETE_DATE.dump.bz2 && gsutil ls gs://lekcije/backup/lekcije_$DELETE_DATE.dump.bz2
