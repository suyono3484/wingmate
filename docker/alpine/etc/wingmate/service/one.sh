#!/bin/sh

export DUMMY_PATH=/usr/local/bin/wmdummy
exec /usr/local/bin/wmexec --setsid --user user1:user1 -- /usr/local/bin/wmstarter