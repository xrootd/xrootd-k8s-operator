#!/usr/bin/env sh

# Script to generate the OLM bundle for the operator.

set -eu

usage() {
    cat << EOD

Usage: `basename $0` <options>

  Available options:
    -t <token>    quay.io token for basic auth
    -r <registry> Registry namespace to use
    -v            Verbose mode
    -h            Show Usage

  Pushes the operator bundle information to quay.io appregistry
EOD
}

token=""
registry=""

# get the options
while getopts vht:r: c ; do
    case $c in
        t) token="$OPTARG" ;;
        r) registry="$OPTARG" ;;
        v) set -x ;;
        h) usage; exit ;;
        \?) usage ; exit 2 ;;
    esac
done

DIR=$(cd "$(dirname "$0")"; pwd -P)
. ${DIR}/env.sh

ROOT_DIR="$(dirname $DIR)"

if [[ -z "$registry" ]]; then
  echo -n "Quay Namespace (ex. johndoe): "
  read registry
fi

if [[ -z "$token" ]]; then
  token=$($DIR/get-quay-token.sh)
  if [[ -z "$token" ]]; then
    echo "[error] Provide quay.io token with -t option."
    echo "To get the token, run ./scripts/get-quay-token.sh"
    exit 1
  fi
fi

# Check operator-courier is present
if ! hash operator-courier; then
  pip install operator-courier || ( \
    echo "[error] Install operator-courier using pip and ensure the binary is in PATH env variable"; \
    exit 2 \
  )
fi

# Push using operator-courier
operator-courier push \
  "deploy/olm-catalog/${XROOTD_OPERATOR_NAME}" \
  "$registry" \
  "$XROOTD_OPERATOR_NAME" \
  "$XROOTD_OPERATOR_VERSION" \
  "$token"