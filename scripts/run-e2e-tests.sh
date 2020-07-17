#!/usr/bin/env sh

set -e

usage() {
    cat << EOD

Usage: `basename $0`

  Available options:
    -v           Verbose mode
    -h           Show Usage
  
  Run all e2e test scripts in tests/ folder
EOD
}

VERBOSE=false

# get the options
while getopts vh: c ; do
    case $c in
        v) VERBOSE=true; set -ux ;;
        h) usage; exit ;;
        \?) usage ; exit 2 ;;
    esac
done

shift $(($OPTIND - 1))

DIR=$(cd "$(dirname "$0")"; pwd -P)

NAMESPACE=$(kubectl config view --minify -o='jsonpath={..namespace}')
NAMESPACE=${NAMESPACE:-default}
INSTANCE=$(kubectl get xrootds.xrootd.org -o=jsonpath='{.items[0].metadata.name}' -n "$NAMESPACE") || (echo "No xrootd instance deployed!"; exit 1;)

SHELL_POD="${INSTANCE}-client"

# Run xrootd client pod
kubectl delete pod "$SHELL_POD" -n "$NAMESPACE" || echo "No existing shell pod found."
kubectl run "$SHELL_POD" --image="qserv/xrootd:latest" --image-pull-policy="IfNotPresent" --restart=Never sleep 3600 -n "$NAMESPACE"
kubectl label pod "$SHELL_POD" "instance=$INSTANCE" "tier=client" -n "$NAMESPACE"

while ! kubectl wait --for=condition=ContainersReady pods -l "instance=$INSTANCE" -n "$NAMESPACE"; do
  echo "Waiting for xrootd pods to be ready..."
  kubectl describe pod -l "instance=$INSTANCE" -n "$NAMESPACE"
done

# Declare Test scripts to use
# Not using Array to make it POSIX-compliant
set -- $(ls $DIR/../tests/e2e/test-*.sh)

# Copy all test files
kubectl cp "$DIR/../tests/e2e/" "$NAMESPACE/$SHELL_POD":"/tmp"

# Wait for cluster to run fine!
sleep 30s

for script in "$@"; do
  if ! kubectl exec "$SHELL_POD" -it -- "/tmp/e2e/$(basename $script)" -i "$INSTANCE" $(if $VERBOSE; then echo -n "-v"; fi); then
    echo "Xrootd Worker - xrootd logs"
    kubectl logs -l "component=xrootd-worker,instance=$INSTANCE" -n "$NAMESPACE" -c xrootd
    echo "Xrootd Redirector - cmsd logs"
    kubectl logs -l "component=xrootd-redirector,instance=$INSTANCE" -n "$NAMESPACE" -c cmsd
    exit 3
  fi
done
