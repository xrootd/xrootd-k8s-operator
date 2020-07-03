#!/usr/bin/env sh

DIR=$(cd "$(dirname "$0")"; pwd -P)
. $DIR/base.sh

xrdfs root://base-xrootd-xrootd-redirector:1094/ ls /data

xrdcp $DIR/dummy root://base-xrootd-xrootd-redirector:1094//data/dummy

xrdfs root://base-xrootd-xrootd-redirector:1094/ ls /data