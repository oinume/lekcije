from __future__ import print_function
from invoke import task
from google.cloud import storage

@task
def hello(ctx):
    client = storage.Client()
    print(client.get_bucket('lekcije'))
