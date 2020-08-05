#!/usr/bin/env bash

# Install and uninstall Xrootd operator

set -eu

usage() {
    cat << EOD

Usage: `basename $0` [options] [cmd]

  Available options:
    -d                  Install from local git repository
    -c CATALOG_DIR      Install XrootdVersion CRs for version catalog
    -n NAMESPACE        Specify namespace (default: kubeconfig/default)
    -u                  Uninstall Xrootd-operator,
                      and related CustomResourceDefinition/CustomResource
    -v                  Verbose mode
    -h                  Show Usage

  Install Xrootd-operator
EOD
}

DIR=$(cd "$(dirname "$0")"; pwd -P)
ROOT_DIR=$(dirname "$DIR")

. $ROOT_DIR/scripts/env.sh

if [[ -z $KUBE_CONTEXT ]]; then
  echo "Set a context (kubectl use-context <context>) out of the following:"
  echo
  kubectl config get-contexts
  exit 1
else
  echo "Current K8S context - $KUBE_CONTEXT"
fi

IS_OPENSHIFT=$(if [ $KUBE_CLUSTER_PROVIDER == minishift ]; then echo -n true; else echo -n false; fi)
DEV_INSTALL=false
UNINSTALL=false
PURGE=true

NAMESPACE=$(kubectl config view --minify --output 'jsonpath={..namespace}')
NAMESPACE=${NAMESPACE:-default}

CATALOG_DIR="$ROOT_DIR/manifests/xrootd-catalog"

# get the options
while getopts vhudn:c: c ; do
    case $c in
        n) NAMESPACE="$OPTARG" ;;
        c) CATALOG_DIR="$OPTARG" ;;
        d) DEV_INSTALL=true ;;
        u) UNINSTALL=true ;;
        v) set -x ;;
        h) usage; exit ;;
        \?) usage ; exit 2 ;;
    esac
done

shift $(($OPTIND - 1))

if [ $UNINSTALL = true ]; then
  kubectl delete deployment,role,rolebinding,serviceaccount xrootd-operator
  kubectl delete crds xrootds.xrootd.org

  echo -e "\nSuccessfully uninstalled Xrootd-operator"
  exit 0
fi

if [ $DEV_INSTALL = true ]; then
  MANIFESTS_DIR="$ROOT_DIR"
else
  MANIFESTS_DIR="https://raw.githubusercontent.com/xrootd/xrootd-k8s-operator/master"
fi

if [ $IS_OPENSHIFT = true ]; then
  echo "PS: Ensure 'oc' binary is in path"
  kapply="oc apply -f"
else
  kapply="kubectl apply -n $NAMESPACE -f"
fi

echo "....... Applying CRDs ......."
apply_crd() {
  crd_yml="$MANIFESTS_DIR/deploy/crds/$1"
  if [ $DEV_INSTALL = false ]; then
    down_path="/tmp/$1"
    wget $crd_yml -O $down_path
    crd_yml=$down_path
  fi
  crd_code=$(if [ $IS_OPENSHIFT = true ]; then sed 's|apiextensions.k8s.io/v1|apiextensions.k8s.io/v1beta1|' $crd_yml; else cat $crd_yml; fi)
  $kapply - << EOF
$crd_code
EOF
}

apply_crd catalog.xrootd.org_xrootdversions_crd.yaml
apply_crd xrootd.org_xrootds_crd.yaml

echo "....... Applying Rules and Service Account ....."
$kapply "$MANIFESTS_DIR"/deploy/service_account.yaml
$kapply "$MANIFESTS_DIR"/deploy/role.yaml
$kapply "$MANIFESTS_DIR"/deploy/role_binding.yaml

echo "....... Applying Operator ......."
op_code=$(sed "s|REPLACE_IMAGE|$XROOTD_OPERATOR_FULL_IMAGE|g" "$MANIFESTS_DIR/deploy/operator.yaml")
$kapply - << EOF
$op_code
EOF

if [ -d "$CATALOG_DIR" ]; then
  echo "....... Applying Xrootd Version Catalogs ......."
  kubectl apply -f $CATALOG_DIR/*
else
  echo "....... Skipping Xrootd Version Catalogs ......."
fi

while ! kubectl wait --for=condition=Ready pods -l name=xrootd-operator -n "$NAMESPACE"
do
  echo "Waiting for operator to be ready..."
  kubectl describe pod -l name=xrootd-operator -n "$NAMESPACE"
done

echo -e "\nSuccessfully installed Xrootd operator in '$NAMESPACE' namespace."
