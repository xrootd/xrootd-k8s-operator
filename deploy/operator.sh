#!/bin/bash

# Install and uninstall Xrootd operator

set -e

usage() {
    cat << EOD

Usage: `basename $0` [options] [cmd]

  Available options:
    -d                  Install from local git repository
    -n NAMESPACE        Specify namespace (default: kubeconfig/default)
    -u                  Uninstall Xrootd-operator,
                      and related CustomResourceDefinition/CustomResource
    -v                  Verbose mode
    -h                  Show Usage

  Install Xrootd-operator
EOD
}

DIR=$(cd "$(dirname "$0")"; pwd -P)

echo -n "Check kubeconfig context - "
kubectl config current-context || {
  echo "Set a context (kubectl use-context <context>) out of the following:"
  echo
  kubectl config get-contexts
  exit 1
}
echo ""

DEV_INSTALL=false
UNINSTALL=false
PURGE=true

NAMESPACE=$(kubectl config view --minify --output 'jsonpath={..namespace}')
NAMESPACE=${NAMESPACE:-default}

# get the options
while getopts vhudn: c ; do
    case $c in
        n) NAMESPACE="$OPTARG" ;;
        d) DEV_INSTALL=true ;;
        u) UNINSTALL=true ;;
        v) ux ;;
        h) usage; exit ;;
        \?) usage ; exit 2 ;;
    esac
done

shift $(($OPTIND - 1))

if [ "$UNINSTALL" = true ]; then

  kubectl delete deployment,role,rolebinding,serviceaccount xrootd-operator
  kubectl delete crds xrootds.xrootd.org

  echo -e "\nSuccessfully uninstalled Xrootd-operator"
  exit 0
fi

if [ "$DEV_INSTALL" = true ]; then
  MANIFESTS_DIR=$(dirname "$DIR")
else
  MANIFESTS_DIR="https://raw.githubusercontent.com/xrootd/xrootd-k8s-operator/master"
fi

kapply="kubectl apply -n $NAMESPACE -f "

echo "....... Applying CRDs ......."
$kapply "$MANIFESTS_DIR"/deploy/crds/xrootd.org_xrootds_crd.yaml
echo "....... Applying Rules and Service Account ....."
$kapply "$MANIFESTS_DIR"/deploy/service_account.yaml
$kapply "$MANIFESTS_DIR"/deploy/role.yaml
$kapply "$MANIFESTS_DIR"/deploy/role_binding.yaml
echo "....... Applying Operator ......."
$kapply "$MANIFESTS_DIR"/deploy/operator.yaml

while ! kubectl wait --for=condition=Ready pods -l name=xrootd-operator -n "$NAMESPACE"
do
  echo "Waiting for operator to be ready..."
  kubectl describe pod -l name=xrootd-operator -n "$NAMESPACE"
done

echo -e "\nSuccessfully installed Xrootd operator in '$NAMESPACE' namespace."
