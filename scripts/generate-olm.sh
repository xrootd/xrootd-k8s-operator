#!/usr/bin/env sh

# Script to generate the OLM artifacts for the operator.

DIR=$(cd "$(dirname "$0")"; pwd -P)
source ${DIR}/env.sh

# Generate CSV 
echo "Generating CSV for version ${XROOTD_OPERATOR_VERSION}"
operator-sdk generate csv \
    --operator-name ${XROOTD_OPERATOR_NAME} \
    --csv-version ${XROOTD_OPERATOR_VERSION} \
    --make-manifests=false \
    --update-crds
