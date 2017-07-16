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
