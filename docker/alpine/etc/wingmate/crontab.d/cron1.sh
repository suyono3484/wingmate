#!/bin/sh

export WINGMATE_DUMMY_PATH=/usr/local/bin/wmdummy
export WINGMATE_LOG=/var/log/cron1.log
export WINGMATE_LOG_MESSAGE="cron executed in minute 5,10,15,20,25,30,35,40,45,50,55"

echo "I'm runnig with dummy=$WINGMATE_DUMMY_PATH, log=$WINGMATE_LOG and mesage=$WINGMATE_LOG_MESSAGE" >> /var/log/debug-cron.log

exec /usr/local/bin/wmoneshot