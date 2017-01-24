#!/usr/bin/bash

# run linters
docker run --rm -it -v `pwd`/go:/go dfinity/build-env gometalinter --deadline=10s --disable=gas --disable=gocyclo --disable=gotype --disable=interfacer src/bls/...
if [ $? -ne 0 ];
then
    echo "failed gometalinter"
    exit 1
fi
echo "OK for gometalinter bls"
docker run --rm -it -v `pwd`/go:/go dfinity/build-env gometalinter --deadline=10s --disable=gas --disable=gocyclo --disable=gotype --disable=interfacer src/blscgo/...
if [ $? -ne 0 ];
then
    echo "failed gometalinter"
    exit 1
fi
echo "OK for gometalinter blscgo"

docker run --rm -it -v `pwd`/go:/go dfinity/build-env gometalinter --deadline=10s --disable=gas --disable=gocyclo --disable=gotype --disable=interfacer src/common2/...
if [ $? -ne 0 ];
then
    echo "failed gometalinter"
    exit 1
fi
echo "OK for gometalinter common2"

docker run --rm -it -v `pwd`/go:/go dfinity/build-env gometalinter --deadline=10s --disable=gas --disable=gocyclo --disable=gotype --disable=interfacer src/sim/...
if [ $? -ne 0 ];
then
    echo "failed gometalinter sim"
    exit 1
fi
echo "OK for gometalinter sim"

docker run --rm -it -v `pwd`/go:/go dfinity/build-env gometalinter --deadline=10s --disable=gas --disable=gocyclo --disable=gotype --disable=interfacer src/state/...
if [ $? -ne 0 ];
then
    echo "failed gometalinter"
    exit 1
fi
echo "OK for gometalinter state"
