#/bin/bash
rootdir=$(cd "$(dirname $0)";pwd)
#mkdir build
go build -o $rootdir/build/api $rootdir/cmd/api/
go build -o $rootdir/build/user $rootdir/cmd/user/
go build -o $rootdir/build/video $rootdir/cmd/video/
go build -o $rootdir/build/comment $rootdir/cmd/comment/

#docker build -t douyin-demo-micro .

export JAEGER_DISABLED=false
export JAEGER_SAMPLER_TYPE="const"
export JAEGER_SAMPLER_PARAM=1
export JAEGER_REPORTER_LOG_SPANS=true
export JAEGER_AGENT_HOST="127.0.0.1"
export JAEGER_AGENT_PORT=6831

build/user &
build/video &
build/comment &
build/api