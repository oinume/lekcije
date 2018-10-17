#!/bin/bash

set -e

echo $GCP_SERVICE_ACCOUNT_KEY | base64 --decode > ./gcp-service-account-key.json
GOOGLE_APPLICATION_CREDENTIALS=./gcp-service-account-key.json invoke $1
