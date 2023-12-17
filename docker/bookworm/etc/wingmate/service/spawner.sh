#!/usr/bin/bash

export WINGMATE_ONESHOT_PATH=/usr/local/bin/wmoneshot
export WINGMATE_DUMMY_PATH=/usr/local/bin/wmdummy
exec /usr/local/bin/wmexec --user 1200 -- /usr/local/bin/wmspawner