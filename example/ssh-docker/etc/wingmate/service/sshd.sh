#!/bin/sh

exec /usr/local/bin/wmpidproxy --pid-file /var/run/sshd.pid -- /usr/sbin/sshd
