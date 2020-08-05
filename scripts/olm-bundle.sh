#!/usr/bin/env sh

# Script to generate the OLM bundle for the operator.

set -eux

DIR=$(cd "$(dirname "$0")"; pwd -P)
. ${DIR}/env.sh

# Create required bundle build and manifest directory
mkdir -p $XROOTD_OPERATOR_BUNDLE_BUILD_DIR/manifests
cp -r $XROOTD_OPERATOR_BUNDLE_DIR/* $XROOTD_OPERATOR_BUNDLE_BUILD_DIR
cp -r $XROOTD_OPERATOR_BUNDLE_MANIFEST_VERSION_DIR/* $XROOTD_OPERATOR_BUNDLE_BUILD_DIR/manifests/

# Build the bundle container image
echo "Building image ${XROOTD_OPERATOR_BUNDLE_IMAGE}"
docker build -t ${XROOTD_OPERATOR_BUNDLE_IMAGE} ${XROOTD_OPERATOR_BUNDLE_BUILD_DIR}
