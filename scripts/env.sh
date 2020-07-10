#!/usr/bin/env sh

# Environment variables used in development. This file is sourced in the other scripts.

export XROOTD_OPERATOR_NAME="xrootd-operator"
export XROOTD_OPERATOR_VERSION=$(awk '$1 == "Version" {gsub(/"/, "", $3); print $3}' version/version.go)
