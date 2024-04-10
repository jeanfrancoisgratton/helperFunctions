#!/usr/bin/env sh

GOROOT=/opt/go
OUTPUT=/opt/tests/

#arch=$(uname -m)


if [ "$#" -gt 0 ]; then
    OUTPUT=$1
fi

cd debug
go build -o ${OUTPUT}/helperFunctions .
