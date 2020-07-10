#!/usr/bin/env sh

# Environment variables used in development. This file is sourced in the other scripts.

DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd -P)
. $DIR/release-support.sh

## General vars
export XROOTD_OPERATOR_NAME="xrootd-operator"
export XROOTD_OPERATOR_VERSION=$(awk '$1 == "Version" {gsub(/"/, "", $3); print $3}' version/version.go)
export XROOTD_OPERATOR_IMAGE_VERSION="$(getVersion)"

## Container Image vars
export XROOTD_OPERATOR_IMAGE_REPO="qserv/$XROOTD_OPERATOR_NAME"
export XROOTD_OPERATOR_IMAGE_TAG="$XROOTD_OPERATOR_IMAGE_VERSION"
export XROOTD_OPERATOR_IMAGE="$XROOTD_OPERATOR_IMAGE_REPO:$XROOTD_OPERATOR_IMAGE_TAG"

## Operator Bundle vars
export XROOTD_OPERATOR_BUNDLE_DIR="deploy/bundle"
export XROOTD_OPERATOR_BUNDLE_BUILD_DIR="build/_output/bundle"
export XROOTD_OPERATOR_BUNDLE_MANIFEST_DIR="deploy/olm-catalog/${XROOTD_OPERATOR_NAME}/${XROOTD_OPERATOR_VERSION}"
export XROOTD_OPERATOR_BUNDLE_METADATA_DIR="$XROOTD_OPERATOR_BUNDLE_DIR/metadata"
export XROOTD_OPERATOR_BUNDLE_IMAGE_REPO="qserv/$XROOTD_OPERATOR_NAME-bundle"
export XROOTD_OPERATOR_BUNDLE_IMAGE_TAG="$XROOTD_OPERATOR_VERSION"
export XROOTD_OPERATOR_BUNDLE_IMAGE="$XROOTD_OPERATOR_BUNDLE_IMAGE_REPO:$XROOTD_OPERATOR_BUNDLE_IMAGE_TAG"

## Ensure go module support is enabled
export GO111MODULE=on
