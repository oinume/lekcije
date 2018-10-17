#!/bin/bash

if [ -z "$HEROKU_LOGIN" -o -z "$HEROKU_API_KEY" ]; then
    echo "Env HEROKU_LOGIN and HEROKU_API_KEY is required."
    exit 1
fi

curl -L -o /tmp/heroku-linux-amd64.tar.gz https://cli-assets.heroku.com/branches/stable/heroku-linux-amd64.tar.gz
tar -xz -C /opt -f /tmp/google-cloud-sdk-$GOOGLE_CLOUD_SDK_VERSION-linux-x86_64.tar.gz
ln -s /opt/heroku/bin/heroku /usr/local/bin/heroku

cat > ~/.netrc << EOF
machine api.heroku.com
login $HEROKU_LOGIN
password $HEROKU_API_KEY
machine git.heroku.com
login $HEROKU_LOGIN
password $HEROKU_API_KEY
EOF

#git remote add heroku https://git.heroku.com/<your-app-name>.git
#wget https://cli-assets.heroku.com/branches/stable/heroku-linux-amd64.tar.gz
#mkdir -p /usr/local/lib /usr/local/bin
#tar -xvzf heroku-linux-amd64.tar.gz -C /usr/local/lib
#ln -s /usr/local/lib/heroku/bin/heroku /usr/local/bin/heroku
