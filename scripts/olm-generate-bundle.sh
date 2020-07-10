#!/usr/bin/env sh

# Script to generate the OLM bundle for the operator.

set -eux

DIR=$(cd "$(dirname "$0")"; pwd -P)
. ${DIR}/env.sh

# Generate Bundle
echo "Generating Bundle for version ${XROOTD_OPERATOR_VERSION}"
operator-sdk bundle create \
    --generate-only