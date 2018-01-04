#!/bin/bash

START=2017-07-30
END=2017-12-31

CURRENT=$START
while true; do
    echo $CURRENT
    daily_reporter -target-date=$CURRENT
    if [ "$CURRENT" = "$END" ]; then
        break
    fi
    CURRENT=`date -d "$CURRENT 1day" +%Y-%m-%d`
done
