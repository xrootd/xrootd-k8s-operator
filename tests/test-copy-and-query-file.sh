#!/usr/bin/env sh

DIR=$(cd "$(dirname "$0")"; pwd -P)
. $DIR/base.sh

xrdfs root://base-xrootd-xrootd-redirector:1094/ ls /data

# Copy data to worker
xrdcp $DIR/dummy root://base-xrootd-xrootd-redirector:1094//data/dummy

# Verify copied data exists
xrdfs root://base-xrootd-xrootd-redirector:1094/ ls /data | grep '/data/dummy'
