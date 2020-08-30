#!/usr/bin/env sh

# Script to deploy xrootd operator with given version from releases

set -eu

# VERSION refers to release tag (if IS_RELEASE is empty) else branch name
VERSION=master
IS_RELEASE=false

BASE_URL=https://github.com/xrootd/xrootd-k8s-operator

usage() {
    cat << EOD

Usage: `basename $0` <options>

  Available options:
    -u            Uninstall
    -v            Verbose mode
    -h            Show Usage

  Deploys xrootd-operator@$VERSION in the currently configured cluster.
  The script installs the required CRDs, roles, permission and controller deployment.
EOD
}

DO_UNINSTALL=false

# get the options
while getopts vhu c ; do
    case $c in
        v) set -x ;;
        h) usage; exit ;;
        u) DO_UNINSTALL=true ;;
        \?) usage ; exit 2 ;;
    esac
done

if [ $IS_RELEASE = true ]; then
    # Get the release assets of that version
    version_url=$BASE_URL/releases/download/$VERSION
    if [ $DO_UNINSTALL = true ]; then
        kubectl delete -f $version_url/install.yaml
    else
        kubectl apply -f $version_url/install.yaml
    fi
else
    tmp_dir=$(mktemp -d)
    cd $tmp_dir
    git clone $BASE_URL
    cd $(basename $BASE_URL)
    if [ $DO_UNINSTALL = true ]; then
        make undeploy
    else
        make deploy
    fi
fi
