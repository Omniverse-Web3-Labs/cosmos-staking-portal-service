#!/bin/bash

if [ -z $WORK_DIR ];then
    DIR=$(dirname "$0")
    export WORK_DIR=$DIR/..
    export WORK_DIR=$(cd $WORK_DIR && pwd)
fi

cd $WORK_DIR

. scripts/env $1

sh scripts/build.sh $APPLICATION && $OUTPUT_DIR/server
