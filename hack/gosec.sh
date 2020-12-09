#!/bin/sh

set -u

result=$(./tools/gosec -quiet ./...)

if [ "$result" != "" ];then
    echo "$result"
    exit 1
fi
