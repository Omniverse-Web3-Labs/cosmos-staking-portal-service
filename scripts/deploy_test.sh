#/bin/bash

if [ -z $WORK_DIR ];then
    DIR=$(dirname "$0")
    export WORK_DIR=$DIR/..
    export WORK_DIR=$(cd $WORK_DIR && pwd)
fi

. scripts/env $1

sh scripts/build.sh

scp bin/release/api/server ubuntu@54.183.118.180:/home/ubuntu

ssh ubuntu@54.183.118.180 'mv server /home/ubuntu/bin/release/api/server && docker restart api'