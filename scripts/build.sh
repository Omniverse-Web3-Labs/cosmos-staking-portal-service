if [ -z $WORK_DIR ];then
    DIR=$(dirname "$0")
    export WORK_DIR=$DIR/..
    export WORK_DIR=$(cd $WORK_DIR && pwd)
fi

. scripts/env $1

echo "build $APPLICATION to $OUTPUT_DIR/server"
GOARCH=amd64 GOOS=linux go build -tags=jsoniter -o $OUTPUT_DIR/server cmd/$APPLICATION/main.go && echo "build success"