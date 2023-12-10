#!/usr/bin/bash

export WINGMATE_DUMMY_PATH=/usr/local/bin/wmdummy
export WINGMATE_LOG=/var/log/cron2.log
export WINGMATE_LOG_MESSAGE="cron scheduled using 17,42 */2 * * *"

exec /usr/local/bin/wmoneshot