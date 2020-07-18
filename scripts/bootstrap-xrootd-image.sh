#!/usr/bin/env sh

set -eu

usage() {
    cat << EOD

Usage: `basename $0`

  Available options:
    -p <provider> Cluster Provider (kind/k3s)
    -c <cluster> Cluster Name
    -v           Verbose mode
    -h           Show Usage

  Script to bootstrap provided cluster with xrootd image
  [TEMPORARY] It's an hack until we can get an official xrootd image hosted at Dockerhub.
EOD
}

provider=""
cluster=""
image="xrootd"
scripts="$(dirname $0)"

# get the options
while getopts vhp:c: c ; do
    case $c in
        p) provider="$OPTARG" ;;
        c) cluster="$OPTARG" ;;
        v) set -x ;;
        h) usage; exit ;;
        \?) usage ; exit 2 ;;
    esac
done

echo "Building '$image' image..."
docker build -t "$image" https://raw.githubusercontent.com/lnielsen/xrootd-docker/master/Dockerfile
cmd="$scripts/load-image.sh"
if [ -n "$provider" ]; then
    cmd+=" -p $provider"
fi
if [ -n "$cluster" ]; then
    cmd+=" -c $cluster"
fi
cmd+=" $image"
$cmd
