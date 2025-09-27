#!/bin/bash

if [ -z $WORK_DIR ];then
    DIR=$(dirname "$0")
    export WORK_DIR=$DIR/..
    export WORK_DIR=$(cd $WORK_DIR && pwd)
fi

. scripts/env $1

echo "build $APPLICATION to $OUTPUT_DIR/bootstrap"
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags "jsoniter lambda.norpc" -o $OUTPUT_DIR/bootstrap cmd/$APPLICATION/main.go
cd $OUTPUT_DIR
zip $APPLICATION.zip bootstrap
