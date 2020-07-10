#!/usr/bin/env sh

# Script to generate the CSV manifest for the operator.

set -eux

DIR=$(cd "$(dirname "$0")"; pwd -P)
. ${DIR}/env.sh

# Generate CSV 
echo "Generating CSV for version ${XROOTD_OPERATOR_VERSION}"
operator-sdk generate csv \
    --operator-name ${XROOTD_OPERATOR_NAME} \
    --csv-version ${XROOTD_OPERATOR_VERSION} \
    --make-manifests=false \
    --update-crds
