#!/bin/bash

set -eu

echo -n "Username: " > /dev/tty
read USERNAME
echo -n "Password: " > /dev/tty
read -s PASSWORD

# Get the token from quay.io API
json=$(curl -H "Content-Type: application/json" -XPOST https://quay.io/cnr/api/v1/users/login -d '
{
    "user": {
        "username": "'"${USERNAME}"'",
        "password": "'"${PASSWORD}"'"
    }
}')

token=$(echo "$json" | jq '.token')
if [ "$token" = "null" ]; then
    echo "[error] failed to get token - $json" > /dev/tty
    exit 2
fi

# Strip the leading and ending quotes
token=$(echo $token | sed -e 's/^"//' -e 's/"$//')

echo -n $token
