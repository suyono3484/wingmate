#!/bin/sh

if [ $# -gt 0 ]; then
    exec "$@"
else
    exec /usr/local/bin/wingmate
fi