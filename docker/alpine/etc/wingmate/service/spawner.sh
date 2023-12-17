#!/bin/sh

export WINGMATE_ONESHOT_PATH=/usr/local/bin/wmoneshot
export WINGMATE_DUMMY_PATH=/usr/local/bin/wmdummy
exec /usr/local/bin/wmexec --user 1001 -- /usr/local/bin/wmspawner