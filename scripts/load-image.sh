#!/usr/bin/env sh

usage() {
    cat << EOD

Usage: `basename $0` <options> [image name]

  Available options:
    -p <provider> Cluster Provider (kind/k3s)
    -c <cluster> Cluster Name
    -v           Verbose mode
    -h           Show Usage

  Load provided image name in the provided cluster.
EOD
}

RE_IMAGE_WORD="[a-z0-9._-]+"
RE_IMAGE_TAG_WORD="[a-zA-Z0-9_.-]+"

provider="kind"
image=""
cluster=""

# get the options
while getopts vhp:c: c ; do
    case $c in
        p) provider="$OPTARG" ;;
        c) cluster="$OPTARG" ;;
        v) set -eux ;;
        h) usage; exit ;;
        \?) usage ; exit 2 ;;
    esac
done

if [[ -z "$provider" ]] || egrep -qv "^(kind|k3s)$" <<< "$provider"; then
    echo "[error] Provide valid Cluster provider name!"
    usage
    exit 2
fi

shift $(($OPTIND - 1))

if [[ $# -eq 0 ]]; then
    usage
    exit 2
fi

image="$1"

if [[ -z "$image" ]] || egrep -qv "^($RE_IMAGE_WORD/)?$RE_IMAGE_WORD(:$RE_IMAGE_TAG_WORD)?$" <<< "$image" ; then
    echo "[error] Provide valid image name!"
    usage
    exit 2
fi

cmd=""
case "$provider" in
    kind)
        cmd="kind load docker-image $image"
        if [[ -n "$cluster" ]]; then
            cmd="$cmd --name $cluster"
        fi
        ;;
    k3s)
        cmd="k3d load image $image"
        if [[ -n "$cluster" ]]; then
            cmd="$cmd --cluster $cluster"
        fi
        ;;
esac

if [[ -z "$cmd" ]]; then
    echo "[error] Not able to generate the command string"
    exit 1
fi

echo "Executing '$cmd'..."
$cmd
