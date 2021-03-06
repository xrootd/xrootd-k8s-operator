set -eu

usage() {
    cat << EOD

Usage: `basename $0` <options>

  Available options:
    -i           Instance name
    -v           Verbose mode
    -h           Show Usage

  Run provided test
EOD
}

echo "Parsing options..."

instance=""

# get the options
while getopts vhi: c ; do
    case $c in
        i) instance="$OPTARG" ;;
        v) set -x ;;
        h) usage; exit ;;
        \?) usage ; exit 2 ;;
    esac
done

shift $(($OPTIND - 1))

echo "Starting test for '$instance'..."
