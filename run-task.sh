#!/bin/bash

set -ex

echo $GCLOUD_SERVICE_KEY | base64 --decode > ./gcloud-service-key.json
GOOGLE_APPLICATION_CREDENTIALS=./gcloud-service-key.json invoke $1
