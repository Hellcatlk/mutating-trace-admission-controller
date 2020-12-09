#!/bin/sh

set -u

result=$(./tools/golint ./...)

if [ "$result" != "" ];then
    echo "$result"
    exit 1
fi
