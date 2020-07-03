#!/usr/bin/env sh

set -eux

usage() {
    cat << EOD

Usage: `basename $0`
EOD
}

DIR=$(cd "$(dirname "$0")"; pwd -P)

NAMESPACE=$(kubectl config view --minify -o='jsonpath={..namespace}')
NAMESPACE=${NAMESPACE:-default}
INSTANCE=$(kubectl get xrootds.xrootd.org -o=jsonpath='{.items[0].metadata.name}' -n "$NAMESPACE")

SHELL_POD="${INSTANCE}-client"

# Run xrootd client pod
kubectl delete pod -l "instance=$INSTANCE,tier=client" -n "$NAMESPACE"
kubectl run "$SHELL_POD" --image="xrootd:latest" --image-pull-policy="IfNotPresent" --restart=Never sleep 3600 -n "$NAMESPACE"
kubectl label pod "$SHELL_POD" "instance=$INSTANCE" "tier=client" -n "$NAMESPACE"

while ! kubectl wait --for=condition=Ready pods -l "instance=$INSTANCE" -n "$NAMESPACE"
do
  echo "Waiting for xrootd pods to be ready..."
  kubectl describe pod -l "instance=$INSTANCE" -n "$NAMESPACE"
done

# Declare Test scripts to use
# Not using Array to make it POSIX-compliant
set -- $(ls tests/test-*.sh)

# Copy all test files
kubectl cp "$DIR/../tests" "$NAMESPACE/$SHELL_POD":"/tmp"
for script in "$@"; do
    kubectl exec "$SHELL_POD" -it "/tmp/$script"
done
