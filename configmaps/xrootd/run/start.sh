#!/usr/bin/env sh

set -eux

usage() {
    cat << EOD

Usage: `basename $0` [options]

  Available options:
    -s <service> Service to start (cmsd/xrootd)
    -h           Show Usage

  Run cmsd/xrootd (ulimit setup).
EOD
}

service=""

# get the options
while getopts hs: c ; do
    case $c in
        s) service="$OPTARG" ;;
        \? | h) usage ; exit 2 ;;
    esac
done

if [[ $# -eq 0 ]] || [[ -z "$service" ]] || (printf '%s\n' "$service" | egrep -qv "^(xrootd|cmsd)$"); then
    usage
    exit 2
fi

shift $(($OPTIND - 1))

CONFIG_DIR="/config-etc"
XROOTD_CONFIG="$CONFIG_DIR/xrootd.cf"

XROOTD_RDR_DN="{{.XrootdRedirectorDn}}"

if hostname | egrep "^${XROOTD_RDR_DN}-[0-9]+$"; then
    COMPONENT_NAME='manager'
else
    COMPONENT_NAME='worker'
fi
export COMPONENT_NAME

# Start service
#
echo "Start $service as $(whoami) user"
cmd="$service -c $XROOTD_CONFIG -n $COMPONENT_NAME"
$cmd
