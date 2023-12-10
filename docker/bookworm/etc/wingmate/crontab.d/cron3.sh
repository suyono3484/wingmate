#!/usr/bin/bash

export WINGMATE_DUMMY_PATH=/usr/local/bin/wmdummy
export WINGMATE_LOG=/var/log/cron3.log
export WINGMATE_LOG_MESSAGE="cron entry: 7,19,23,47 22 * * * /etc/wingmate/crontab.d/cron3.sh"

exec /usr/local/bin/wmoneshot