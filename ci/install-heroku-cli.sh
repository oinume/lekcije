#!/bin/bash

curl -L -o /tmp/heroku-linux-amd64.tar.gz https://cli-assets.heroku.com/branches/stable/heroku-linux-amd64.tar.gz
tar -xz -C /opt -f /tmp/heroku-linux-amd64.tar.gz
#ln -s /opt/heroku/bin/heroku /usr/local/bin/heroku

#git remote add heroku https://git.heroku.com/<your-app-name>.git
#wget https://cli-assets.heroku.com/branches/stable/heroku-linux-amd64.tar.gz
#mkdir -p /usr/local/lib /usr/local/bin
#tar -xvzf heroku-linux-amd64.tar.gz -C /usr/local/lib
#ln -s /usr/local/lib/heroku/bin/heroku /usr/local/bin/heroku
