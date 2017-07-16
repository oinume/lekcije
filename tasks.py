from __future__ import print_function
import datetime
from google.cloud import storage
from invoke import task
import os

BUCKET_NAME = 'lekcije'

@task
def backup_mysql(ctx):
    mysqldump = os.getenv('MYSQLDUMP')
    if not mysqldump:
        mysqldump = './bin/mysqldump'
    user = os.getenv('MYSQL_USER')
    password = os.getenv('MYSQL_PASSWORD')
    host = os.getenv('MYSQL_HOST')
    port = os.getenv('MYSQL_PORT')
    database = os.getenv('MYSQL_DATABASE')
    dump_file = 'lekcije_' + datetime.datetime.now().strftime('%Y%m%d') + '.dump.bz2'
    ctx.run('{mysqldump} -u{user} -p{password} -h{host} -P{port} {database} | bzip2 -9 > {dump_file}'.format(**locals()))

    client = storage.Client()
    bucket = client.get_bucket(BUCKET_NAME)
    new_blob = bucket.blob('backup/' + dump_file)
    new_blob.upload_from_filename(dump_file)

    delete_date = (datetime.datetime.now() - datetime.timedelta(days=7)).strftime('%Y%m%d')
    delete_blob_name = 'backup/lekcije_' + delete_date + '.dump.bz2'
    delete_blob = bucket.get_blob(delete_blob_name)
    if delete_blob:
        delete_blob.delete()


# MYSQLDUMP="$CWD/$DIR/mysqldump"
# #echo $MYSQLDUMP
# DATE=`date "+%Y%m%d"`
# $MYSQLDUMP -u$MYSQL_USER -p$MYSQL_PASSWORD -h$MYSQL_HOST $MYSQL_DATABASE | bzip2 -9 > lekcije_$DATE.dump.bz2
# echo $GCLOUD_SERVICE_KEY | base64 --decode > ./gcloud-service-key.json
# ./google-cloud-sdk/bin/gcloud components update -q
# ./google-cloud-sdk/bin/gcloud \
#     auth activate-service-account \
#          --key-file ./gcloud-service-key.json
#
# DELETE_DATE=`date "+%Y%m%d" --date "7 days ago"`
#
# ./google-cloud-sdk/bin/gsutil cp lekcije_$DATE.dump.bz2 gs://lekcije/backup/
# ./google-cloud-sdk/bin/gsutil ls gs://lekcije/backup/lekcije_$DELETE_DATE.dump.bz2 && ./google-cloud-sdk/bin/gsutil rm gs://lekcije/backup/lekcije_$DELETE_DATE.dump.bz2
